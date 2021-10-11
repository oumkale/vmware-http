echo "Stopping stress-ng chaos" 
kill -9 $(ps aux | grep [s]tress-ng | awk '{print $2}')