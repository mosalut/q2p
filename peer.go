package q2p

import (
	"net"
	"log"
)

func Q2P(ip string, port int) error {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(ip), Port: port})
	if err != nil {
		return err
	}

	log.Println(listener.LocalAddr().String())

	go read(listener)

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

func join(hostAddrs []string, networkIdentify string) {
	for k, addr := range hostAddrs {
		remoteAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			log.Fatal(err)
		}

		_, err = listener.WriteToUDP([]byte(networkIdentify), remoteAddr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
