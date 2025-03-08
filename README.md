# sniffle

Simple Go program for logging the connections made by a device on the local network.

It just logs the SRC and DST IPs and packet types. It doesn't attempt to look at the contents of the packets.

Command line options are:

- -a <ip>   IP address of device to sniff
- -d <dev>  Name of a predefined device to sniff
- -f <name> Name of a predefined filter to use

They are treated in that priority - ie, if -a is used, then -d and -f are ignored. If -d is used, -f is ignored.

```sh
go get github.com/google/gopacket
go get github.com/google/gopacket/pcap
go get github.com/google/gopacket/layers
```

## Documentation

[GoPacket](https://github.com/google/gopacket)

- [gopacket](https://pkg.go.dev/github.com/google/gopacket)
- [pcap](https://pkg.go.dev/github.com/google/gopacket@v1.1.19/pcap#section-documentation)
- [layers](https://pkg.go.dev/github.com/google/gopacket/layers#section-documentation)

[Info on BPF filters](https://www.ibm.com/docs/en/qsip/7.5?topic=queries-berkeley-packet-filters)