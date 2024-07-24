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
	case JOIN:
		log.Println("event: JOIN")
		peer.TouchRequest(rAddr)
		peer.RemoteSeeds = append(peer.RemoteSeeds, rAddr)
	case TOUCHREQUEST:
		log.Println("event: TOUCHREQUEST")
		log.Println("rAddr3:", string(data[4:]))
		rAddr3, err := net.ResolveUDPAddr("udp", string(data[4:]))
		if err != nil {
			log.Println(err)
		}

		peer.Touch(rAddr, rAddr3)
	case TOUCH:
		log.Println("event: TOUCH")
	case TOUCHED:
		log.Println("event: TOUCHED")
		log.Println("rAddr3:", string(data[4:]))
		rAddr3, err := net.ResolveUDPAddr("udp", string(data[4:]))
		if err != nil {
			log.Println(err)
		}

		peer.ConnectRequest(rAddr, rAddr3) // here rAddr is rAddr2
	case CONNECTREQUEST:
		log.Println("event: CONNECTREQUEST")
		log.Println("rAddr2:", string(data[4:]))
		rAddr2, err := net.ResolveUDPAddr("udp", string(data[4:]))
		if err != nil {
			log.Println(err)
		}

		peer.Connect(rAddr2)
	case CONNECT:
		log.Println("event: CONNECT")
		log.Println("from:", rAddr.String())
		peer.RemoteSeeds = append(peer.RemoteSeeds, rAddr)

		peer.Connected(rAddr)
	case CONNECTED:
		log.Println("event: CONNECTED")
		log.Println("from:", rAddr.String())
		peer.RemoteSeeds = append(peer.RemoteSeeds, rAddr)
	default:
		log.Println("Undefined event")
	}

	return nil
}
