#!/bin/bash

# Sends a notification whenever a callsign of interest appears on-air.

# Your Pushbullet API key here
APIKEY=""

while true; do

  # Replace <SomePattern> with regex that matches the station(s) for which you want on-air alerts.
  # Using ': ' will notify for tx from them, and omitting that will notify for tx to or from them.
  #grep '<SomePattern>: ' ${HOME}/.local/share/JS8Call/DIRECTED.TXT > /tmp/js8c-special.new;
  grep '<SomePattern' ${HOME}/.local/share/JS8Call/DIRECTED.TXT > /tmp/js8c-special.new;

  NEW=$(diff /tmp/js8c-special.old /tmp/js8c-special.new |grep -o '[A-Z0-9]*: .*')

  cat /tmp/js8c-special.new > /tmp/js8c-special.old

  NEWCHARS=$(echo ${NEW} | wc -c)
  if [ "${NEWCHARS}" -gt "3" ]; then
    echo "Special Call:: ${NEW}" |head -c100;
    curl https://api.pushbullet.com/api/pushes -u ${APIKEY}: \
        --output /dev/null --silent --max-time 5 \
        -X POST \
        -d type=note -d title="JS8Call Callsign On-Air" \
        -d body="$(echo ${NEW} |head -c100)"
  fi

  sleep 60;

done
