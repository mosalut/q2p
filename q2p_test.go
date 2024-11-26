package q2p

import (
	"testing"
	"syscall"
	"os"
	"os/signal"
	"flag"
	"net"
	"time"
	"log"
	"fmt"
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

	seedAddrs := make(map[string]bool)
	if cmdFlag.remoteHost != "" {
		seedAddrs[cmdFlag.remoteHost] = false
	}

	peer := NewPeer(cmdFlag.ip, cmdFlag.port, seedAddrs, cmdFlag.networkID, callback)
	err := peer.Run()
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	t.Log("Received signal, shutting down...")
}

func TestTransport(t *testing.T) {
	t.Log(*cmdFlag)

	seedAddrs := make(map[string]bool)
	if cmdFlag.remoteHost != "" {
		seedAddrs[cmdFlag.remoteHost] = false
	}

	peer := NewPeer(cmdFlag.ip, cmdFlag.port, seedAddrs, cmdFlag.networkID, callback)
	err := peer.Run()
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission")
OUTER:
	for {
		for key, _ := range peer.RemoteSeeds {
			rAddr, err := net.ResolveUDPAddr("udp", key)
			if err != nil {
				t.Error(err, key)
			}

			err = peer.Transport(rAddr, data)
			if err != nil {
				t.Error(err, key)
			}
			t.Log("returned returned returned returned returned returned returned returned returned")
			break OUTER
		}
		t.Log("waiting waiting waiting")
		time.Sleep(time.Second)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	t.Log("Received signal, shutting down...")
}

func callback(data []byte) {
	fmt.Println(string(data))
}

func TestTransport10001(t *testing.T) {
	t.Log(*cmdFlag)

	seedAddrs := make(map[string]bool)
	if cmdFlag.remoteHost != "" {
		seedAddrs[cmdFlag.remoteHost] = false
	}

	peer := NewPeer(cmdFlag.ip, cmdFlag.port, seedAddrs, cmdFlag.networkID, callback)
	err := peer.Run()
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission")

	key := "127.0.0.1:10001"
	time.Sleep(time.Second * 3)
	rAddr, err := net.ResolveUDPAddr("udp", key)
	if err != nil {
		t.Error(err, key)
	}

	err = peer.Transport(rAddr, data)
	if err != nil {
		t.Error(err, key)
	}
	t.Log("returned returned returned returned returned returned returned returned returned")

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
