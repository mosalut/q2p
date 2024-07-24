package q2p

import (
	"encoding/binary"
	"net"
	"log"
)

type peer_T struct {
	IP string
	Port int
	RemoteSeeds []*net.UDPAddr
	NetworkID uint16
	Conn *net.UDPConn
	Callback func(*peer_T, *net.UDPAddr, uint16, []byte) error
}

func NewPeer(ip string, port int, rAddrs []*net.UDPAddr, networkID uint16) *peer_T {
	return &peer_T {ip, port, rAddrs, networkID, nil, nil}
}

func (peer *peer_T)Run() error {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(peer.IP), Port: peer.Port})
	if err != nil {
		return err
	}
	peer.Conn = listener

	log.Println(listener.LocalAddr().String())

	go peer.read()

	/*
	if peer.RHAddr != "" {
		joinOne(listener, peer.RHAddr, peer.NetworkID)
	}
	*/

	peer.join()

	return nil
}

func (peer *peer_T)read() {
	data := make([]byte, 1024)
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

		peer.Callback(peer, remoteAddr, event, data[:n])
		/*
		_, err = peer.Conn.WriteToUDP([]byte("world"), remoteAddr)
		if err != nil {
			log.Println(err)
		}
		*/
	}
}

func (peer *peer_T)join() {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, JOIN)
	header = append(header, bs...)

	for _, seed := range peer.RemoteSeeds {
		_, err := peer.Conn.WriteToUDP(header, seed)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (peer *peer_T)TouchRequest() {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, TOUCH)
	header = append(header, bs...)

	for _, seed := range peer.RemoteSeeds {
		_, err := peer.Conn.WriteToUDP(header, seed)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (peer *peer_T)ConnectRequest(rAddr *net.UDPAddr) {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, CONNECT)
	header = append(header, bs...)

	_, err := peer.Conn.WriteToUDP(header, rAddr)
	if err != nil {
		log.Fatal(err)
	}
}
