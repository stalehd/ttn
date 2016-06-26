#!/bin/bash
export TTNCTL_DEBUG=true
export TTNCTL_TTN_ROUTER=localhost:1700
export TTNCTL_MQTT_BROKER=localhost:1883
export TTNCTL_TTN_HANDLER=localhost:1782
export TTNCTL_TTN_ACCOUNT_SERVER=http://localhost:8080
release/ttnctl-darwin-amd64 $1 $2 $3 $4 $5
