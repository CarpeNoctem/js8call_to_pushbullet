package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// TODO: Add CLI option to specify config file
// TODO: Errs to stderr

type Config struct {
	MyCall     string   `yaml:"my_call"`
	MyGroups   []string `yaml:"my_groups"`
	Pushbullet struct {
		Token string `yaml:"token"`
		// Ident string `yaml:"identity"`
	} `yaml:"pushbullet"`
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Notifications struct {
		CQ                  bool     `yaml:"cq"`
		Heartbeat           bool     `yaml:"heartbeat"`
		HeartbeatAck        bool     `yaml:"heartbeat_ack"`
		DirectMsg           bool     `yaml:"direct_msg"`
		RxSpot              bool     `yaml:"callsign_spot"`
		RigTx               bool     `yaml:"rig_tx"`
		SpecialCallMentions bool     `yaml:"special_call_mentions"`
		SpecialCalls        []string `yaml:"special_calls"`
		IgnoreCalls         []string `yaml:"ignore_calls"`
	} `yaml:"notifications"`
}

type Js8Event struct {
	Params struct {
		CMD    string  `json:"CMD"`
		CALL   string  `json:"CALL"`
		DIAL   int     `json:"DIAL"`
		EXTRA  string  `json:"EXTRA"`
		FREQ   int     `json:"FREQ"`
		FROM   string  `json:"FROM"`
		GRID   string  `json:"GRID"`
		OFFSET int     `json:"OFFSET"`
		SNR    int     `json:"SNR"`
		SPEED  int     `json:"SPEED"`
		TDRIFT float64 `json:"TDRIFT"`
		TEXT   string  `json:"TEXT"`
		TO     string  `json:"TO"`
		UTC    int64   `json:"UTC"`
	} `json:"params"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func main() {
	var config Config
	loadConfig("config.yml", &config)

	events := make(chan string, 20)

	go parseEvents(events, config)

	conn, err := net.Dial("tcp", fmt.Sprint(config.Server.Host, ":", config.Server.Port))
	if err != nil {
		println("Couldn't connect to JS8Call:", err.Error())
		os.Exit(1)
	}

	println("Connected!")

	for {
		reply, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				// TODO: Add reconnect attempt. Can run in a bash loop in the meantime...
				println("Connection lost with JS8Call. Quitting.")
				pushNotification("JS8Call Connection Problem", "Connection lost with JS8Call. Quitting.", config) // purposely do not run this in an async goroutine. We wait for the notification to fire before exiting.
				os.Exit(1)
			} else {
				println("Woah, that was unexpected:", err.Error())
			}
		}
		events <- reply
	}

	conn.Close()
}

func parseEvents(events <-chan string, config Config) {
	for rawEvent := range events {
		//println("DEBUG: Got raw event:", rawEvent)
		var eventObject Js8Event
		jsonErr := json.Unmarshal([]byte(rawEvent), &eventObject)
		if jsonErr != nil {
			println("Hmm, couldn't parse event:", jsonErr.Error())
		}
		//println("DEBUG: Parsed event of type", eventObject.Type)
		//println("DEBUG: Command is:", eventObject.Params.CMD)
		//fmt.Printf("DEBUG: %+v\n", eventObject)
		handleEvent(eventObject, config)
	}
}

func handleEvent(event Js8Event, config Config) {
	switch event.Type {
	case "RX.SPOT":
		if inSlice(event.Params.CALL, config.Notifications.IgnoreCalls) { // TODO: test if this can be moved out a level (does CALL == "" or == nil when not present in json?)
			return
		}
		if inSlice(event.Params.CALL, config.Notifications.SpecialCalls) {
			go pushNotification("JS8Call Special Call On-Air Direct", fmt.Sprintf("Call: %s SNR:%d QTH:%s", event.Params.CALL, event.Params.SNR, event.Params.GRID), config)
		} else if config.Notifications.RxSpot { // don't double-up on spots
			go pushNotification("JS8Call Callsign Spot", fmt.Sprintf("Call: %s SNR:%d QTH:%s", event.Params.CALL, event.Params.SNR, event.Params.GRID), config)
		}

	case "RIG.PTT":
		if event.Value == "on" && config.Notifications.RigTx {
			go pushNotification("JS8Call TX", "Rig transmitting", config)
		}
	case "RX.DIRECTED":
		if inSlice(event.Params.FROM, config.Notifications.IgnoreCalls) {
			return
		}
		if event.Params.CMD == " CQ" { // Separating this and the below condition so that CQs don't trigger the myGroups condition below when CQ notification setting disabled.
			if config.Notifications.CQ {
				go pushNotification("JS8Call CQ CQ CQ", fmt.Sprintf("Call: %s SNR:%d QTH:%s", event.Params.FROM, event.Params.SNR, event.Params.GRID), config)
			}
		} else if event.Params.CMD == " HEARTBEAT" { // Splitting this and below cond for same reason as above, except to weed out from beta debugging log. can probably put back onto one line after release.
			if config.Notifications.Heartbeat {
				go pushNotification("JS8Call Heartbeat", fmt.Sprintf("Call: %s SNR:%d QTH:%s", event.Params.FROM, event.Params.SNR, event.Params.GRID), config)
			}
		} else if event.Params.CMD == " SNR" || event.Params.CMD == " HEARTBEAT SNR" {
			if event.Params.TO == config.MyCall && config.Notifications.HeartbeatAck { // Separating this and the above condition so that heartbeat acks don't trigger below new msg condition when heartbeat ack notifications disabled
				go pushNotification("JS8Call Heartbeat ACK", fmt.Sprintf("My SNR:%s From %s their SNR:%d", event.Params.EXTRA, event.Params.FROM, event.Params.SNR), config)
			}
		} else if event.Params.TO == config.MyCall {
			replacement := fmt.Sprintf("%s: %s ", event.Params.FROM, event.Params.TO)
			message := strings.Replace(event.Value, replacement, "", 1)
			if len(message) > 100 {
				message = fmt.Sprint(message[0:100], "...")
			}
			go pushNotification(fmt.Sprintf("JS8Call new msg from %s", event.Params.FROM), message, config)
		} else if inSlice(event.Params.TO, config.MyGroups) {
			replacement := fmt.Sprintf("%s: %s ", event.Params.FROM, event.Params.TO)
			message := strings.Replace(event.Value, replacement, "", 1)
			if len(message) > 100 {
				message = fmt.Sprint(message[0:100], "...")
			}
			go pushNotification(fmt.Sprintf("JS8Call msg to @ALLCALL frm %s", event.Params.FROM), message, config)
		} else if config.Notifications.SpecialCallMentions && inSlice(event.Params.TO, config.Notifications.SpecialCalls) {
			go pushNotification("JS8Call Special Call On-Air Indirect", fmt.Sprintf("Call: %s Sent from:%s SNR:%d", event.Params.TO, event.Params.FROM, event.Params.SNR), config)
		} //else if event.Params.CMD != " HEARTBEAT SNR" {
		//	fmt.Printf("DEBUG: Other Directed: %+v\n", event)
		//}
		//		default:
		//			fmt.Printf("DEBUG: Ignoring: %+v\n", event)
	}
}

func pushNotification(title string, message string, config Config) {
	fmt.Printf("* %s *\n  - %s\n\n", title, message)

	reqBody := url.Values{
		"type":  {"note"},
		"title": {title},
		"body":  {message},
	}

	req, err := http.NewRequest("POST", "https://api.pushbullet.com/api/pushes", strings.NewReader(reqBody.Encode()))

	if err != nil {
		fmt.Println("Couldn't build req:", err.Error())
	}

	req.SetBasicAuth(config.Pushbullet.Token, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Problem with call to Pushbullet:", err.Error())
		return
	}
	if resp.Status != "200 OK" {
		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Problem sending notification via Pushbullet: %s\n", respBody)
		//fmt.Println("Status code:", resp.Status)
	}
	resp.Body.Close()
}

// TODO: check for missing config items (namely mycall, api token, ...)
// TODO: convirt config values to uppercase upon loading
func loadConfig(filename string, config *Config) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Could not open", filename, ":", err.Error())
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Could not parse", filename, ":", err.Error())
	}
}

func inSlice(needle string, haystack []string) bool {
	for _, item := range haystack {
		if needle == item {
			return true
		}
	}
	return false
}
