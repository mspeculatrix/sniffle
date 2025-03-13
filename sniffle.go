/*
SNIFFLE

A simple program to sniff packets going to or from a specific device on the
local network.
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
	snaplen = int32(1024)
	promisc = true
)

var (
	devIP      = ""        // IPv4 address for a specific device
	devName    = ""        // for searching in devices map
	filterName = ""        // for searching in filters map
	filter     = ""        // BPF filter string
	filePrefix = "sniffle" // default prefix
	iface      = "enp3s0"  // should probably make this a var & cmd line option
	timeoutT   = 30        // timeout in seconds
	verbose    = false     // naturally silent
	defaults   = map[string]string{
		// defaults to be used if no settings.cfg file found
		"captureDir": ".",
		"captureExt": "log",
	}
)

func main() {
	/*  READ CONFIG FILES  */
	// Read settings.cfg file
	settings, settingsErr := filelib.ReadConfig("sniffle.cfg")
	if settingsErr != nil {
		// Had a problem reading settings file, so use defaults
		settings = defaults
	}
	// Read the devices.cfg file to create the devices map
	devices, _ := filelib.ReadConfig("localhosts.cfg")
	// Read the filters.cfg file to create the filters map
	filters, _ := filelib.ReadConfig("bpf_filters.cfg")

	/*  GET COMMAND LINE PARAMETERS  */
	flag.StringVar(&devIP, "a", devIP, "IP of device to sniff")
	flag.StringVar(&devName, "d", devName, "Name of predefined device to sniff")
	flag.StringVar(&filterName, "f", filter, "Predefined filter")
	flag.StringVar(&iface, "i", iface, "Interface to use fo sniffing")
	flag.IntVar(&timeoutT, "t", timeoutT, "Timeout in secs")
	flag.BoolVar(&verbose, "v", false, "Use verbose mode")
	flag.Parse()

	/*  SET FILTER  */
	var filterOK = false
	// First we'll look to see if the user passed a device IP using the -a flag
	if devIP != "" {
		filter = "host " + devIP
		filterOK = true
		filePrefix = devIP
		// If not, was a device name specified with the -d flag
	} else if devName != "" {
		// If so, look up the name in the devices map
		if devices != nil {
			val, ok := devices[devName]
			if ok {
				filter = "host " + val
				filterOK = true
				filePrefix = devName
			} else {
				log.Fatal("Invalied device name.")
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
				filePrefix = filterName
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
		captureFile := filePrefix + "_" + ts + "." + settings["captureExt"]
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
		// Create a NewTicker. This sends ticks out on the channel timer.C
		timer := time.NewTicker(timeout)
		defer timer.Stop() // not strictly necessary since Go 1.23

		for {
			select {
			case packet := <-packetSource.Packets():
				// If we've received a packet, handle it
				packethandler.HandlePacket(packet, fh, verbose)
			case <-timer.C:
				// If there is nothing on this channel, it means the ticker
				// has finished, and so we're done.
				if verbose {
					fmt.Println("-- done")
				}
				return
			}
		}

	} else {
		log.Fatal("No valid filter")
	}
}
