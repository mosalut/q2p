package q2p

import (
	"crypto/md5"
	"encoding/binary"
	"context"
	"net"
	"errors"
	"time"
	"fmt"
	"log"
)

type Peer_T struct {
	IP string `json:"ip"`
	Port int `json:"port"`
	RemoteSeeds map[string]bool `json:"remote_seeds"`
	NetworkID uint16 `json:"network_id"`
	Conn *net.UDPConn `json:"conn"`
	Callback func([]byte) `json:"-"`
}

func NewPeer(ip string, port int, rAddrs map[string]bool, networkID uint16, callback func([]byte)) *Peer_T {
	return &Peer_T {ip, port, rAddrs, networkID, nil, callback}
}

func (peer *Peer_T)Run() error {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(peer.IP), Port: peer.Port})
	if err != nil {
		return err
	}
	peer.Conn = listener

	log.Println(listener.LocalAddr().String())
	log.Println(peer.RemoteSeeds)

	go peer.read()

	peer.join()

	return nil
}

func (peer *Peer_T)read() {
	data := make([]byte, 512)
	for {
		n, remoteAddr, err := peer.Conn.ReadFromUDP(data)
		if err != nil {
			log.Println("error during read:", err)
		}
	//	log.Println(remoteAddr, n, data[:n])

		if n < 4 {
			log.Println("Invalid q2p header length")
			continue
		}
		networkID := binary.LittleEndian.Uint16(data[:2])
		event := binary.LittleEndian.Uint16(data[2:4])

		log.Println("networkID =", networkID)
		log.Println("event =", event)

		if networkID != peer.NetworkID {
			log.Println("Different q2p network id", networkID, peer.NetworkID)
			continue
		}

		log.Println(networkID)

		err = peer.networking(remoteAddr, event, data[:n], peer.Callback)
		if err != nil {
			log.Println(err)
		}
	}
}

func (peer *Peer_T)join() {
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
			log.Println(err)
		}
	}
}

func (peer *Peer_T)TouchRequest(rAddr3 *net.UDPAddr) {
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

		_, err = peer.Conn.WriteToUDP(data, seedAddr)
		if err != nil {
			log.Println(err)
		}
	}

}

func (peer *Peer_T)Touch(rAddr, rAddr3 *net.UDPAddr) {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, TOUCH)
	header = append(header, bs...)

	_, err := peer.Conn.WriteToUDP(header, rAddr3)
	if err != nil {
		log.Println(err)
	}

	header = header[:2]
	binary.LittleEndian.PutUint16(bs, TOUCHED)
	header = append(header, bs...)

	data := append(header, []byte(rAddr3.String())...)

	_, err = peer.Conn.WriteToUDP(data, rAddr)
	if err != nil {
		log.Println(err)
	}
}

func (peer *Peer_T)ConnectRequest(rAddr2, rAddr3 *net.UDPAddr) {
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
		log.Println(err)
	}
}

func (peer *Peer_T)Connect(rAddr2 *net.UDPAddr) {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, CONNECT)
	header = append(header, bs...)

	_, err := peer.Conn.WriteToUDP(header, rAddr2)
	if err != nil {
		log.Println(err)
	}
}

func (peer *Peer_T)Connected(rAddr3 *net.UDPAddr) {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, CONNECTED)
	header = append(header, bs...)

	_, err := peer.Conn.WriteToUDP(header, rAddr3)
	if err != nil {
		log.Println(err)
	}
}

func (peer *Peer_T)Transport(rAddr *net.UDPAddr, data []byte) error {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, TRANSPORT)
	header = append(header, bs...)

	hash := md5.Sum(data)
	fmt.Println("hash:", hash)

	key := fmt.Sprintf("%x", hash)
	fmt.Println(key)

	length := len(data)
	if(length > 2078764170780) { // 2078764170780 = math.MaxUint32 * 484, 484 is each packet's body length
		return errors.New("too long data: should be less than 2078764170780")
	}
	packetNum := length / 484

	if length % 484 != 0 {
		packetNum++
	}

	bs = make([]byte, 4, 4)

	transmissionHead := make([]byte, 0, 24)
	transmissionHead = append(transmissionHead, hash[:]...) // 0th ~ 16th bytes: hash
	binary.LittleEndian.PutUint32(bs, uint32(length))
	transmissionHead = append(transmissionHead, bs...)  // 16th ~ 20nd bytes: length
	transmissionHead = append(transmissionHead, []byte{0, 0, 0, 0}...) // 20nd ~ 24th bytes: leave space empty for each SYN

	transmissionM[key] = make(map[uint32][]byte)

	for i := 0; i < packetNum; i++ {
		/* for packet losing test
		if i == 3 || i == 5 {
			continue
		}
		*/
		start := i * 484
		end := start + 484
		if(length < end) {
			end = length
		}

		log.Println("start:", start)
		body := data[start: end]
		transmissionM[key][uint32(i)] = body

		// SYN
		binary.LittleEndian.PutUint32(bs, uint32(i))
		copy(transmissionHead[20:], bs)

		transm := append(transmissionHead, body...)
		transmission := append(header, transm...)
		_, err := peer.Conn.WriteToUDP(transmission, rAddr)
		if err != nil {
			log.Println(err)
		}
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Second * 10)
	go transmissionSending(ctx, key, rAddr.String())

	return nil

}

func (peer *Peer_T)transportFailed(rAddr *net.UDPAddr, hash []byte, syns []uint32) {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, TRANSPORT_FAILED)
	header = append(header, bs...)

	bs = make([]byte, 4, 4)
	body := hash
	for _, v := range syns {
		binary.LittleEndian.PutUint32(bs, v)
		body = append(body, bs...)
	}

	transmission := append(header, body...)
	_, err := peer.Conn.WriteToUDP(transmission, rAddr)
	if err != nil {
		log.Println(err)
	}
}
