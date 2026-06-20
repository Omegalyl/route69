# route69

Me learning TCP/IP and networking by building a raw packet sniffer in Go.
No libpcap, no gopacket — just using raw sockets and parsing bytes by hand.

## Run

```bash
make route69
sudo ./dist/route69   # needs root for raw sockets

OR

make run
```

Interface name is hardcoded in `cmd/route69/main.go` — swap it for yours.

