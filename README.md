# JS8Call to Pushbullet

## Notification Limit!
Well, since writing this app, I found out that Pushbullet allow a limited number of pushes [(500) per month](https://docs.pushbullet.com/#push-limit) for free accounts. Just want to let you know that up front before you get too excited. I'd never run close to that limit before, but this morning, I got a notification to let me know I was approaching it. This means you'll probably want to 
* be selective with which notification types you enable;
* not run the app when you're in the shack / present at the keyboard / within earshot of audible notifications;
* not leave the app running overnight while you're sleeping and not getting the notifications anyway.
* Or you could sign up to a Pushbullet [Pro subscription (about $40USD/yr)](https://www.pushbullet.com/pro).

## JS8Call
(http://js8call.com/) is an awesome amateur radio mode and program that allows many forms of messaging and more, on top of the very popular FT8.

## Pushbullet
(https://www.pushbullet.com/) is a service and mobile application that I've used for some time to get alerts for various things sent to my phone.

## This
is a small application that can notify you on your phone when there is activity on your JS8Call station.
Personally, I find this useful because I often leave my station running JS8Call when the radio would otherwise be turned off.
This allows for propogation reporting, asyncronous messaging, and relaying between other stations.
While my station is running, I may afk, but close enough to answer activity if I have time. I might be doing some yard work, reading, etc. 
Getting notifications on my phone lets me make the most out of my time and answer JS8Call activity without having to stay sat in the shack or constantly watching a VNC screen.

If this sounds useful for you as well, go ahead and give it a try!

Prerequesites
-------------
* You'll need a Pushbullet account, and you'll need to grab an API token from that account.
* You'll also need to install their app on the device on which you wish to receive notifications.
* That device will need to be connected to the Internet.


Getting Started
---------------
First, either download the application binary for your system from the releases here: https://github.com/CarpeNoctem/js8call_to_pushbullet/releases/tag/v0.1
or compile the binary from the source code in this repository.

Next, copy the config.yml file into the same directory as the application.
Then fill in your callsign, any groups your in, your Pushbullet API token, and set up your notification settings.

After that, make sure you've enabled API access in JS8Call under File>Settings>Reporting>API> "Enable TCP Server API" & "Accept TCP Requests"

Run this app!


Further Notes
-------------
* This is pretty fresh off the keyboard. I've encountered and fixed a few bugs along the way, and I know there will still be a couple more to be fixed.
* There is also still some cleanup to be done. For example, I haven't yet added case in-sensitivity to the app, so please use UPPERCASE for any callsigns in your config.yml
* On the note of YAML, I found that the --- and ... delimiters from the YAML spec weren't necessary for the yaml parsing library I used, so I removed them to simplify the file a little.
* Additionally, for those not familiar with YAML: if you want an empty ignore_list or special_calls, use brackets: []
* If you want to list a few callsigns for either, then use and intent, hyphen, and callsign for each call, and start below the setting text. The default config file shows both in use, so I hope this is clear enough.
* This is not a great solution for off-the-grid ops, as it relies on the third-party Pushbullet.
If you're planning on running off-the-grid or just don't want to use Pushbullet, then I suggest using the built-in sound notification options within JS8Call. I use this with a speaker plugged into my Raspberry Pi or with a laptop, but that doesn't help when I'm not within earshot, hence this app.
Perhaps in the future, I can add the capability for this app to play sound files. That way, one could run it from a remote device with network connectivity to the JS8Call device so that you can get sound notifications remotely (but without Pushbullet or the Internet). You could also consider using a bluetooth speaker connected to the JS8Call device.

