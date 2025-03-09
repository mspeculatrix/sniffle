/*
A simple program to sniff packets going to or from a specific device.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sniffle/filelib"
	"sniffle/packethandler"

	//"net"

	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const (
	iface      = "enp3s0"
	snaplen    = int32(1024)
	promisc    = true
	timeoutT   = 30
	captureDir = "/mnt/nas/sync/captures"
)

var (
	devIP      = ""
	devName    = ""
	filterName = ""
	filter     = ""
	devices    = map[string]string{}
	filters    = map[string]string{}
	verbose    = false
)

func main() {
	/*  READ CONFIG FILES  */
	devices, _ = filelib.ReadKVFile("config/devices.cfg", ":")
	filters, _ = filelib.ReadKVFile("config/filters.cfg", ":")

	/*  GET COMMAND LINE FLAGS  */
	flag.StringVar(&devIP, "a", "", "IP of device to sniff")
	flag.StringVar(&devName, "d", "", "Name of device to sniff")
	flag.StringVar(&filterName, "f", "", "Predefined filter")
	flag.BoolVar(&verbose, "v", false, "Use verbose mode")
	flag.Parse()

	/*  SET FILTER  */
	var filterOK = false
	if devIP != "" {
		filter = "host " + devIP
		filterOK = true
	} else if devName != "" {
		if devices != nil {
			val, ok := devices[devName]
			if ok {
				filter = "host " + val
				filterOK = true
			}
		}
	} else if filterName != "" {
		if filters != nil {
			val, ok := filters[filterName]
			if ok {
				filter = val
				filterOK = true
			}
		}
	} else {
		log.Fatal("No option provided")
	}

	if filterOK {
		filter += " && not broadcast && not arp && not multicast"

		var timeout time.Duration = time.Duration(timeoutT) * time.Second

		// Open file for capture
		ts := time.Now().Format("20060102_150405")
		captureFile := captureDir + "/iot_" + ts + ".log"
		fh, err := os.Create(captureFile)
		if err != nil {
			log.Fatal(err)
		}
		// ensure file gets closed
		defer fh.Close()

		// Open interface
		handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()

		// Apply BPF filter
		filterErr := handle.SetBPFFilter(filter)
		if filterErr != nil {
			log.Fatalf("error applyign BPF Filter %s - %v", filter, err)
		}

		if verbose {
			fmt.Println("SNIFFLE")
			fmt.Println("-------")
			fmt.Println("Sniffing on interface:", iface)
			fmt.Println("Saving to:", captureDir)
			fmt.Println("Using filter:", filter)
		}
		// Start stream
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		// Pull packets from stream
		for packet := range packetSource.Packets() {
			packethandler.HandlePacket(packet, fh, verbose)
		}
	} else {
		log.Fatal("No valid filter")
	}
}
