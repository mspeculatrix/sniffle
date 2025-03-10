package packethandler

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// The callback function for handling packets as they come in.
func HandlePacket(packet gopacket.Packet, fileh *os.File, verbose bool) {
	var ip4SrcIP = ""
	var ip4DstIP = ""
	var ip4Proto = ""
	var ptype = ""
	var srcPort uint16 = 0
	var dstPort uint16 = 0
	var note = ""
	logpkt := false

	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		logpkt = true
		tcp, _ := tcpLayer.(*layers.TCP) // Get TCP data from this layer
		srcPort = uint16(tcp.SrcPort)
		dstPort = uint16(tcp.DstPort)
		if tcp.ACK {
			ptype = "ACK"
		}
		if tcp.FIN {
			ptype = "FIN"
		}
		if tcp.RST {
			ptype = "RST"
		}
		if tcp.SYN {
			ptype = "SYN"
		}
		// Other TCP properties
		// PSH, URG, ECE, CWR, NS - bool
		// Ack - uint32
		//
	}

	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		logpkt = true
		udp, _ := udpLayer.(*layers.UDP) // Get UDP data from this layer
		srcPort = uint16(udp.SrcPort)
		dstPort = uint16(udp.DstPort)
	}

	if ip4Layer := packet.Layer(layers.LayerTypeIPv4); ip4Layer != nil {
		logpkt = true
		ip4, _ := ip4Layer.(*layers.IPv4) // Get TCP data from this layer
		ip4SrcIP = ip4.SrcIP.String()
		ip4DstIP = ip4.DstIP.String()
		ip4Proto = ip4.Protocol.String()
	}

	if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
		ptype = "DNS"
		dns, _ := dnsLayer.(*layers.DNS) // Get TCP data from this layer
		if dns.Answers != nil {
			answers := dns.Answers // of type DNSResourceRecord
			IPs := []string{}
			for answ := range answers {
				IPs = append(IPs, answers[answ].IP.String())
			}
			note = "ans " + strings.Join(IPs, ", ")
		} else {
			question := dns.Questions
			note = "qry " + string(question[0].Name)
		}
	}

	if ntpLayer := packet.Layer(layers.LayerTypeNTP); ntpLayer != nil {
		ptype = "NTP"
	}

	if logpkt {
		ts := time.Now().Format("20060102_150405") // timestamp string
		logline := fmt.Sprintf("%-15s %-7s %3s %15s : %-5d %15s : %-5d  %-s\n", ts, ip4Proto, ptype, ip4SrcIP, srcPort, ip4DstIP, dstPort, note)
		if verbose {
			fmt.Print(logline)
		}
		_, ferr := fileh.WriteString(logline)
		if ferr != nil {
			log.Fatal(ferr)
		}
	}
}
