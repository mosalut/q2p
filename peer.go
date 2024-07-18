package q2p

import (
	"net"
	"log"
)

type q2pParams_T struct {
	IP string
	Port int
	RHAddr string
	NetworkIdentify string
}

var remoteHostAddresses = make([]string, 0, 16)

func Q2P(params *q2pParams_T) error {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(params.IP), Port: params.Port})
	if err != nil {
		return err
	}

	log.Println(listener.LocalAddr().String())

	go read(listener)

	if params.RHAddr != "" {
		joinOne(listener, params.RHAddr, params.NetworkIdentify)
	}

	join(listener, params.NetworkIdentify)

	return nil
}

func read(conn *net.UDPConn) {
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println("error during read:", err)
		}

		log.Println(remoteAddr, string(data[:n]))

		/*
		_, err = conn.WriteToUDP([]byte("world"), remoteAddr)
		if err != nil {
			log.Println(err)
		}
		*/
	}
}

func joinOne(conn *net.UDPConn, rHAddr, networkIdentify string) {
	remoteAddr, err := net.ResolveUDPAddr("udp", rHAddr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.WriteToUDP([]byte(networkIdentify), remoteAddr)
	if err != nil {
		log.Println(err)
	}
}

func join(conn *net.UDPConn, networkIdentify string) {
	for _, addr := range remoteHostAddresses {
		remoteAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			log.Fatal(err)
		}

		_, err = conn.WriteToUDP([]byte(networkIdentify), remoteAddr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
