package q2p

import (
	"encoding/binary"
	"net"
	"log"
)

type peer_T struct {
	IP string
	Port int
	RemoteSeeds map[string]bool
	NetworkID uint16
	Conn *net.UDPConn
}

func NewPeer(ip string, port int, rAddrs map[string]bool, networkID uint16) *peer_T {
	return &peer_T {ip, port, rAddrs, networkID, nil}
}

func (peer *peer_T)Run() error {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(peer.IP), Port: peer.Port})
	if err != nil {
		return err
	}
	peer.Conn = listener

	log.Println(listener.LocalAddr().String())

	go peer.read()

	peer.join()

	return nil
}

func (peer *peer_T)read() {
	data := make([]byte, 64)
	for {
		n, remoteAddr, err := peer.Conn.ReadFromUDP(data)
		if err != nil {
			log.Println("error during read:", err)
		}
		log.Println(remoteAddr, n, data[:n])

		if n < 4 {
			log.Println("Invalid q2p header length", binary.LittleEndian.Uint16(data[:2]), binary.LittleEndian.Uint16(data[2:]))
			continue
		}

		networkID := binary.LittleEndian.Uint16(data[:2])
		event := binary.LittleEndian.Uint16(data[2:])

		if networkID != peer.NetworkID {
			log.Println("Different q2p network id", networkID, peer.NetworkID)
			continue
		}

		log.Println(networkID)

		err = peer.networking(remoteAddr, event, data[:n])
		if err != nil {
			log.Println(err)
		}
	}
}

func (peer *peer_T)join() {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, JOIN)
	header = append(header, bs...)

	for seed, _ := range peer.RemoteSeeds {
		seedAddr, err := net.ResolveUDPAddr("udp", seed)
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = peer.Conn.WriteToUDP(header, seedAddr)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (peer *peer_T)TouchRequest(rAddr3 *net.UDPAddr) {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, TOUCHREQUEST)
	header = append(header, bs...)

	data := append(header, []byte(rAddr3.String())...)

	for seed, _ := range peer.RemoteSeeds {
		log.Println("seed:", seed)
		if seed == rAddr3.String() {
			continue
		}

		seedAddr, err := net.ResolveUDPAddr("udp", seed)
		if err != nil {
			log.Println(err)
			continue
		}

		for i := 0; i < 1; i++ {
			_, err = peer.Conn.WriteToUDP(data, seedAddr)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

func (peer *peer_T)Touch(rAddr, rAddr3 *net.UDPAddr) {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, TOUCH)
	header = append(header, bs...)

	_, err := peer.Conn.WriteToUDP(header, rAddr3)
	if err != nil {
		log.Fatal(err)
	}

	header = header[:2]
	binary.LittleEndian.PutUint16(bs, TOUCHED)
	header = append(header, bs...)

	data := append(header, []byte(rAddr3.String())...)

	_, err = peer.Conn.WriteToUDP(data, rAddr)
	if err != nil {
		log.Fatal(err)
	}
}

func (peer *peer_T)ConnectRequest(rAddr2, rAddr3 *net.UDPAddr) {
	if rAddr2.String() == rAddr3.String() {
		return
	}

	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, CONNECTREQUEST)
	header = append(header, bs...)

	data := append(header, []byte(rAddr2.String())...)

	_, err := peer.Conn.WriteToUDP(data, rAddr3)
	if err != nil {
		log.Fatal(err)
	}
}

func (peer *peer_T)Connect(rAddr2 *net.UDPAddr) {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, CONNECT)
	header = append(header, bs...)

	_, err := peer.Conn.WriteToUDP(header, rAddr2)
	if err != nil {
		log.Fatal(err)
	}
}

func (peer *peer_T)Connected(rAddr3 *net.UDPAddr) {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, CONNECTED)
	header = append(header, bs...)

	_, err := peer.Conn.WriteToUDP(header, rAddr3)
	if err != nil {
		log.Fatal(err)
	}
}
