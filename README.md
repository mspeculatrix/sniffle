# sniffle

Simple Go program for logging the connections made by a device on the local network.

It just logs the SRC and DST IPs and packet types. It doesn't attempt to look at the contents of the packets.

Command line options are:

- -a <ip>   IP address of device to sniff
- -d <dev>  Name of a predefined device to sniff
- -f <name> Name of a predefined filter to use

They are treated in that priority - ie, if -a is used, then -d and -f are ignored. If -d is used, -f is ignored.

After cloning or downloading this code, `cd` into the sniffle directory.

You'll need some packages.

```sh
go get github.com/google/gopacket
go get github.com/google/gopacket/pcap
go get github.com/google/gopacket/layers
```

And then:

```sh
go mod init sniffle
```



## Documentation

[GoPacket](https://github.com/google/gopacket)

- [gopacket](https://pkg.go.dev/github.com/google/gopacket)
- [pcap](https://pkg.go.dev/github.com/google/gopacket@v1.1.19/pcap#section-documentation)
- [layers](https://pkg.go.dev/github.com/google/gopacket/layers#section-documentation)

[Info on BPF filters](https://www.ibm.com/docs/en/qsip/7.5?topic=queries-berkeley-packet-filters)

## Guarantees

There are no guarantees. This is the work of a dilettante amateur coder. Use at your own risk. If your house explodes or the internet catches fire, don't come running to me.
