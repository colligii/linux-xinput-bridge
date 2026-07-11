#!/bin/bash

evtest=$(jq -r '.evtest' defaultConfig.json)
sudo chmod +r "$evtest"
echo "Everything is setted, try run npm start"