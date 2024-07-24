package q2p

import (
	"testing"
	"syscall"
	"os"
	"os/signal"
	"flag"
	"net"
	"log"
)

type cmdFlag_T struct {
	ip string
	port int
	remoteHost string
	networkID uint16
}

var cmdFlag *cmdFlag_T

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmdFlag = &cmdFlag_T{}
	readFlags(cmdFlag)
	cmdFlag.networkID = 0
}

func TestQ2P(t *testing.T) {
	t.Log(*cmdFlag)

	seedAddrs := make(map[*net.UDPAddr]bool)
	if cmdFlag.remoteHost != "" {
		remoteAddr, err := net.ResolveUDPAddr("udp", cmdFlag.remoteHost)
		if err != nil {
			t.Fatal(err)
		}

		seedAddrs[remoteAddr] = false
	}

	peer := NewPeer(cmdFlag.ip, cmdFlag.port, seedAddrs, cmdFlag.networkID)
	peer.Callback = callback
	err := peer.Run()
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	t.Log("Received signal, shutting down...")
}

func readFlags(cmdFlag *cmdFlag_T) {
	flag.StringVar(&cmdFlag.ip, "ip", "0.0.0.0", "UDP host IP")
	flag.IntVar(&cmdFlag.port, "port", 10000, "UDP host Port")
	flag.StringVar(&cmdFlag.remoteHost, "remote_host", "", "remote host address")
}
