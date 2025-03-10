/*
A simple program to sniff packets going to or from a specific device.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sniffle/filelib"
	"sniffle/packethandler"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const (
	iface   = "enp3s0"
	snaplen = int32(1024)
	promisc = true
)

var (
	devIP      = ""
	devName    = ""
	filterName = ""
	filter     = ""
	fileID     = "sniffle"
	timeoutT   = 30
	verbose    = false
	defaults   = map[string]string{
		"captureDir": ".",
		"captureExt": "log",
	}
)

func main() {
	/*  READ CONFIG FILES  */
	// Read settings.cfg file
	settings, settingsErr := filelib.ReadKVFile("config/settings.cfg", ":")
	if settingsErr != nil {
		// Had a problem reading settings file, so use defaults
		settings = defaults
	}
	// Read the devices.cfg file to create the devices map
	devices, _ := filelib.ReadKVFile("config/devices.cfg", ":")
	// Read the filters.cfg file to create the filters map
	filters, _ := filelib.ReadKVFile("config/filters.cfg", ":")

	/*  GET COMMAND LINE FLAGS  */
	flag.StringVar(&devIP, "a", devIP, "IP of device to sniff")
	flag.StringVar(&devName, "d", devName, "Name of predefined device to sniff")
	flag.StringVar(&filterName, "f", filter, "Predefined filter")
	flag.IntVar(&timeoutT, "t", timeoutT, "Timeout in secs")
	flag.BoolVar(&verbose, "v", false, "Use verbose mode")
	flag.Parse()

	/*  SET FILTER  */
	var filterOK = false
	// First we'll look to see if the user passed a device IP using the -a flag
	if devIP != "" {
		filter = "host " + devIP
		filterOK = true
		fileID = devIP
		// If not, was a device name specified with the -d flag
	} else if devName != "" {
		// If so, look up the name in the devices map
		if devices != nil {
			val, ok := devices[devName]
			if ok {
				filter = "host " + val
				filterOK = true
				fileID = devName
			}
		}
		// Did the user specify a named filter with the -f flag?
	} else if filterName != "" {
		// If so, look up the name in the filters map.
		if filters != nil {
			val, ok := filters[filterName]
			if ok {
				filter = val
				filterOK = true
				fileID = filterName
			}
		}
	} else {
		log.Fatal("No filter set")
	}

	if filterOK {
		// We'll add the following by default as, generally, these are not
		// things I want. Your mileage may vary, so feel free to amend.
		filter += " && not broadcast && not arp && not multicast"

		// Open file for capture
		ts := time.Now().Format("20060102_150405")
		captureFile := fileID + "_" + ts + "." + settings["captureExt"]
		capturePath := filepath.Join(settings["captureDir"], captureFile)
		fh, err := os.Create(capturePath)
		if err != nil {
			log.Fatal(err)
		}
		defer fh.Close() // ensure file gets closed

		// Open interface
		var timeout time.Duration = time.Duration(timeoutT) * time.Second
		handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()

		// Apply BPF filter
		filterErr := handle.SetBPFFilter(filter)
		if filterErr != nil {
			log.Fatalf("Error applying BPF Filter %s - %v", filter, err)
		}

		if verbose {
			fmt.Println("SNIFFLE")
			fmt.Println("-------")
			fmt.Println("Sniffing on interface:", iface)
			fmt.Println("Saving to:", capturePath)
			fmt.Println("Using filter:", filter)
		}
		// Start stream
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		// Get timer going for timeout
		timer := time.NewTicker(timeout)
		defer timer.Stop()

		for {
			select {
			case packet := <-packetSource.Packets():
				// If we've received a packet, handle it
				packethandler.HandlePacket(packet, fh, verbose)
			case <-timer.C:
				// If the time has finished, our work is done
				return
			}
		}

	} else {
		log.Fatal("No valid filter")
	}
}
