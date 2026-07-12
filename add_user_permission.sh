#!/bin/bash

sudo chown "$USER":"$USER" /dev/input/event*
sudo chmod 660 /dev/input/event*

echo "All setted for user: $USER."