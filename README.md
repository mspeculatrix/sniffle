# sniffle

Simple Go program for logging the connections made by a device on the local network.

It just logs the SRC and DST IPs and packet types. It doesn't attempt to look at the contents of the packets.

Command line options are:

- -a \<ip>   IP address of device to sniff
- -d \<dev>  Name of a predefined device to sniff
- -f \<name> Name of a predefined BPF-format filter to use
- -t \<secs> Timeout in secs - default 30
- -v         Turn on verbose mode (otherwise there's no STDOUT output)

The first three are treated in that priority - ie, if -a is used, then -d and -f are ignored. If -d is used, -f is ignored.

After cloning or downloading this code, `cd` into the sniffle directory and run the `setup.sh` script once.

```sh
./setup.sh
```

This just gets the necessary packages and runs init - something along the lines of:

```sh
go get github.com/google/gopacket
go get github.com/google/gopacket/pcap
go get github.com/google/gopacket/layers
go mod init sniffle
```

You then run the `build` script to create the executable. I've made this a separate script because you'll no doubt want to tinker with the source code and can use this handy shortcut after each change.

You need to modify the first four variables in the build script to suit the platform you're building for.

```sh
TARGET_OS=linux         # linux, darwin, windows
TARGET_ARCH=amd64       # arm (RPi), amd64 (Intel), arm64 (Apple M-series)
ARM_VERSION=""          # mostly used for RPis. Typically 7
TARGET_BIN=             # name to give binary, default is directory
```

By default, it comes configured for Linux on amd64. If you're not building for ARM you can ignore ARM_VERSION. You can also ignore TARGET_BIN, in which case the executable will be called `sniffle`.

```sh
./build
```

## Documentation

[GoPacket](https://github.com/google/gopacket)

- [gopacket](https://pkg.go.dev/github.com/google/gopacket)
- [pcap](https://pkg.go.dev/github.com/google/gopacket@v1.1.19/pcap#section-documentation)
- [layers](https://pkg.go.dev/github.com/google/gopacket/layers#section-documentation)

[Info on BPF filters](https://www.ibm.com/docs/en/qsip/7.5?topic=queries-berkeley-packet-filters)

## Guarantees

There are no guarantees. This is the work of a dilettante amateur coder. Use at your own risk. If your house explodes or the internet catches fire, don't come running to me.
