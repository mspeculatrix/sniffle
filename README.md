# sniffle

Simple Go program for logging the connections made by a device on the local network.

Intended for (and tested on) macOS and Linux.

It just logs the SRC and DST IPs and packet types. It doesn't attempt to look at the contents of the packets. If you want clever stuff, use Wireshark, tshark or tcpdump.

Command line options are:

- -a \<ip>   IP address of device to sniff. REQUIRED.
- -d \<dev>  Name of a predefined device to sniff. Optional.
- -f \<name> Name of a predefined BPF-format filter to use. Optional.
- -i \<if>   Interface to use for sniffing (default: enp3s0). Optional.
- -t \<secs> Timeout in secs (default: 30). Optional.
- -v         Turn on verbose mode (otherwise there's no STDOUT output)

The first three are treated in that priority - ie, if -a is used, then -d and -f are ignored. If -d is used, -f is ignored. At least one of these options must be used. Sniffle can't guess what you want to do.

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

With all that done, run:

```sh
./build
```

## Config files

You need these config files. The program looks for them in the following locations (and in this order of priority): `/etc/`, an `etc/` subdir in the current working directory and the current working directory itself.

- **localhosts.cfg** - a key/value list of devices on the local network, in the format `<name>: <ip>`. It might be useful to have this in `/etc/` on *nix systems.
- **snarfle.cfg** - a key/value list of basic config settings, in the format `<name>: <value>`.
- **bpf_filters.cfg** - a key/value list of basic config settings, in the format `<name>: <filter_string>`.

Examples are provided in the etc folder. I recommend editing them then copying them to your /etc folder.

The `localhosts.cfg` file is (or can be) shared with snarfle, so if you already have that installed in `/etc/` you don't need this copy.

## Documentation

[GoPacket](https://github.com/google/gopacket)

- [gopacket](https://pkg.go.dev/github.com/google/gopacket)
- [pcap](https://pkg.go.dev/github.com/google/gopacket@v1.1.19/pcap#section-documentation)
- [layers](https://pkg.go.dev/github.com/google/gopacket/layers#section-documentation)

[Info on BPF filters](https://www.ibm.com/docs/en/qsip/7.5?topic=queries-berkeley-packet-filters)

## Guarantees

There are no guarantees. This is the work of a dilettante amateur coder. Use at your own risk. If your house explodes or the internet catches fire, don't come running to me.

[ENDS]
