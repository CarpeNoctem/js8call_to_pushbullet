#!/bin/bash

# Sends a notification whenever there are new CQ calls.

# Your Pushbullet API key here
APIKEY=""

while true; do

  grep ': @ALLCALL CQ' ${HOME}/.local/share/JS8Call/DIRECTED.TXT > /tmp/js8c-cq.new;

  NEW=$(diff /tmp/js8c-cq.old /tmp/js8c-cq.new |grep -o '[A-Z0-9]*: .*')

  cat /tmp/js8c-cq.new > /tmp/js8c-cq.old

  NEWCHARS=$(echo ${NEW} | wc -c)
  if [ "${NEWCHARS}" -gt "3" ]; then
    echo "CQ: ${NEW}" |head -c100;
    curl https://api.pushbullet.com/api/pushes -u ${APIKEY}: \
        --output /dev/null --silent --max-time 5 \
        -X POST \
        -d type=note -d title="JS8Call CQ CQ CQ" \
        -d body="$(echo ${NEW} |head -c100)"
  fi

  sleep 30;

done
