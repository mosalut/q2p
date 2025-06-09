package q2p

import (
	"sync"
	"testing"
	"syscall"
	"os"
	"os/signal"
	"flag"
	"net"
	"time"
)

type cmdFlag_T struct {
	ip string
	port int
	remoteHost string
	networkID uint16
}

var cmdFlag *cmdFlag_T
var transmissionCache = make(map[string][]byte)
var mutex0 = &sync.Mutex{}

func init() {
	cmdFlag = &cmdFlag_T{}
	readFlags(cmdFlag)
	cmdFlag.networkID = 0
}

func TestQ2P(t *testing.T) {
	seedAddrs := make(map[string]bool)
	if cmdFlag.remoteHost != "" {
		seedAddrs[cmdFlag.remoteHost] = false
	}

	t.Log(*cmdFlag)

	peer := NewPeer(cmdFlag.ip, cmdFlag.port, seedAddrs, cmdFlag.networkID)
	peer.TimeSendLost = 3
	peer.Timeout = 10
//	peer.LifeCycle = lifeCycle
	peer.Successed = successed
	peer.Failed = failed
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
	seedAddrs := make(map[string]bool)
	if cmdFlag.remoteHost != "" {
		seedAddrs[cmdFlag.remoteHost] = false
	}

	t.Log(*cmdFlag)

	peer := NewPeer(cmdFlag.ip, cmdFlag.port, seedAddrs, cmdFlag.networkID)
	peer.TimeSendLost = 3
	peer.Timeout = 10
	peer.LifeCycle = lifeCycle
	peer.Successed = successed
	peer.Failed = failed
	err := peer.Run()
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	t.Log("Received signal, shutting down...")
}

func lifeCycle(peer *Peer_T, rAddr *net.UDPAddr, event int) {
	print(log_info, "on life cycle", EventName[event], ":", rAddr.String())
	switch event {
	case STARTRUN:
		line := []byte("\nHello, transmission")
		data := []byte("")

/*
		key, err := peer.Transport(rAddr, data)
		if err != nil {
			print(log_debug, err, key)
		}
		print(log_debug, "returned returned returned returned returned returned returned returned returned")

		transmissionCache[key] = data
		*/

		for {
			for i := 0; i < 100; i++ {
				data = append(data, line...)
			}

			key, err := peer.Transport(rAddr, data)
			if err != nil {
				print(log_error, err, key)
			}
			print(log_debug, "returned returned returned returned returned returned returned returned returned")
			mutex0.Lock()
			transmissionCache[key] = data
			mutex0.Unlock()
		}
	}
}

func successed(peer *Peer_T, rAddr *net.UDPAddr, key string, data []byte) {
	print(log_info, "Successed transmission hash:", key)
	print(log_info, "Received data:", string(data))
	mutex0.Lock()
	delete(transmissionCache, key)
	mutex0.Unlock()
}

func failed(peer *Peer_T, rAddr *net.UDPAddr, key string, syns []uint32) {
	if len(syns) == 0 {
		print(log_warning, "transmission Failed, hash:", key)
		return
	}

	print(log_info, "Lost packet: The hash in transmission:", key)
	mutex0.Lock()
	data := transmissionCache[key]
	mutex0.Unlock()

	var start int
	var end int
	for _, syn := range syns {
		start = int(syn) * PACKET_LEN
		if start + PACKET_LEN > len(data) {
			end = len(data)
		} else {
			end = start + PACKET_LEN
		}

	//	print(log_debug, "lost SYN:", syn, start, end, "data:", string(data[start:end]))
		print(log_debug, "lost SYN:", syn, start, end)

		err := peer.TransportAPacket(rAddr, key, syn, data[start:end])
		if err != nil {
			print(log_error, err)
		}
	}
}




func TestTransport10001(t *testing.T) {
	seedAddrs := make(map[string]bool)
	if cmdFlag.remoteHost != "" {
		seedAddrs[cmdFlag.remoteHost] = false
	}

	peer := NewPeer(cmdFlag.ip, cmdFlag.port, seedAddrs, cmdFlag.networkID)
	err := peer.Run()
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission\nHello, transmission")

	addr := "127.0.0.1:10001"
	time.Sleep(time.Second * 3)
	rAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		t.Error(err, addr)
	}

	key, err := peer.Transport(rAddr, data)
	if err != nil {
		t.Error(err, addr)
	}
	t.Log("returned returned returned returned returned returned returned returned returned")
	mutex0.Lock()
	transmissionCache[key] = data
	mutex0.Unlock()

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
