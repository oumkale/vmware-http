#!/bin/bash
echo "Deleting toxics"
toxiproxy-cli delete ${ToxicName}
echo "Stopping toxiproxy server"
kill -9 $(ps aux | grep [t]oxiproxy | awk '{print $2}')