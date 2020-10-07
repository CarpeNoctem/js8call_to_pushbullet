#!/bin/bash

# Sends a notification whenever someone is sending you a new, non HEARTBEAT-related message.

# Your Pushbullet API key here
APIKEY=""

# Your Callsign here
MYCALL=""

while true; do

  grep -v HEARTBEAT ${HOME}/.local/share/JS8Call/DIRECTED.TXT |grep ": ${MYCALL}" > /tmp/js8c.new;

  NEW=$(diff /tmp/js8c.old /tmp/js8c.new |grep -o "[A-Z0-9]*: ${MYCALL}.*")

  cat /tmp/js8c.new > /tmp/js8c.old

  NEWCHARS=$(echo ${NEW} | wc -c)
  if [ "${NEWCHARS}" -gt "3" ]; then
    echo 'New message';
    curl https://api.pushbullet.com/api/pushes -u ${APIKEY}: \
        --output /dev/null --silent --max-time 5 \
        -X POST \
        -d type=note -d title="New JS8Call Message" \
        -d body="$(echo ${NEW} |head -c100)"
  fi

  sleep 10;

done
