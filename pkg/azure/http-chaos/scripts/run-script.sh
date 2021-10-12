#!/bin/bash
if  [[ "${InstallDependency}" == True ]] ; then
    if [[ "$( which toxiproxy-server 2>/dev/null )" ]] ; then echo Dependency is already installed. 
    else 
        echo "Installing required dependencies" 
        if cat /etc/issue | grep -i Ubuntu ; then
            sudo apt-get update -y -qq
            wget -O toxiproxy-2.1.4.deb https://github.com/Shopify/toxiproxy/releases/download/v2.1.4/toxiproxy_2.1.4_amd64.deb -q
            sudo dpkg -i toxiproxy-2.1.4.deb
        else
            echo "There was a problem installing dependencies."
            exit 1
        fi
    fi
fi

echo "Starting toxiproxy server"
toxiproxy-server > /dev/null 2>&1 &
sleep 2s
declare -a toxics=($(echo ${ToxicType} | tr "," " "))
declare -a values=($(echo ${ToxicValue} | tr "," " "))
port=${ListenPort}

echo "Creating toxics"
toxiproxy-cli create ${ToxicName} --listen localhost:${ListenPort} --${StreamType} localhost:${StreamPort}

for i in "${!toxics[@]}"; do
    toxiproxy-cli toxic add ${ToxicName} --type ${toxics[i]} --attribute ${toxics[i]}=${values[i]}
done
echo "Process completed"