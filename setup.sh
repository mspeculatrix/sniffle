#!/usr/bin/env bash

echo "Setting up sniffle."
echo "You shouldn't have to run this script more than once."
echo
echo "Getting packages..."
go get github.com/google/gopacket
go get github.com/google/gopacket/pcap
go get github.com/google/gopacket/layers

echo "Initialising project..."
go mod init sniffle

echo "Run ./build to build the executable."

exit 0
