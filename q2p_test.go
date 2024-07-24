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
	cmdFlag = &cmdFlag_T{}
	readFlags(cmdFlag)
	cmdFlag.networkID = 0
}

func TestQ2P(t *testing.T) {
	t.Log(*cmdFlag)

	seedAddrs := make([]*net.UDPAddr, 0, 8)
	if cmdFlag.remoteHost != "" {
		remoteAddr, err := net.ResolveUDPAddr("udp", cmdFlag.remoteHost)
		if err != nil {
			t.Fatal(err)
		}

		seedAddrs = append(seedAddrs, remoteAddr)
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

func callback(peer *peer_T, rAddr *net.UDPAddr, event uint16, data []byte) error {
	log.Println("event:", event)
	switch event {
	case 0:
		log.Println("event: JOIN")
		peer.TouchRequest()
		peer.RemoteSeeds = append(peer.RemoteSeeds, rAddr)
	case 1:
		log.Println("event: TOUCH")
		peer.ConnectRequest(rAddr)
	case 2:
		log.Println("event: CONNECT")
	default:
		log.Println("Undefined event")
	}

	log.Println(data)
	return nil
}
