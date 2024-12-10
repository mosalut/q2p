package q2p // github.com/mosalut/q2p

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"context"
	"net"
	"errors"
	"time"
	"fmt"
	"log"
)

const PACKET_LEN = 484

type Peer_T struct {
	IP string `json:"ip"` // Host address where the node starts. The default is 0.0.0.0
	Port int `json:"port"` // Port where the node starts. The default is 10000.
	RemoteSeeds map[string]bool `json:"remote_seeds"` // List of seed nodes. The default is empty.
	NetworkID uint16 `json:"network_id"`
	Conn *net.UDPConn `json:"conn"`
	TimeSendLost int `json:"time_send_again"` // How often the receiver checks for lost packets. If there are lost packets, it will inform the sender. The default value is 5(SECs).
	Timeout int `json:"timeout"` // Timeout for the receiver to wait for complete data. If it times out, it will inform the sender. The default value is 5(SECs).
	Callback func(*Peer_T, *net.UDPAddr, string, []byte) `json:"-"` // Callback function executed when complete data is received. The default function member is to type success transmission HASH.
	CallbackFailed func(*Peer_T, *net.UDPAddr, string, []uint32) `json:"-"` // Callback function executed upon failure. If the length of the last parameter is 0, it indicates a timeout; otherwise, it represents the SYN position where the packet was lost. The default function member is to type timeout transmission HASH or lost transmission HASH that lost packet in with SYNs' positions.
}

// Create a new peer
// Parameter list:
//	* see Peer_T's members
func NewPeer(ip string, port int, rAddrs map[string]bool, networkID uint16) *Peer_T {
	return &Peer_T {
		ip, port, rAddrs, networkID, nil, 6, 5,
		func(peer *Peer_T, rAddr *net.UDPAddr, key string, body []byte) {
			fmt.Println("Succeeded transmission hash:", key)
		},
		func(peer *Peer_T, rAddr *net.UDPAddr, key string, syns []uint32) {
			if len(syns) == 0 {
				fmt.Println("transmission Failed, timeout, hash:", key)
				return
			}

			fmt.Println("Lost packet: The hash in transmission:", key)
			fmt.Println("Lost packet: The SYNs:", syns)
		},
	}
}

// Start peer
func (peer *Peer_T)Run() error {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(peer.IP), Port: peer.Port})
	if err != nil {
		return err
	}
	peer.Conn = listener

	go peer.read()

	peer.join()

	return nil
}

func (peer *Peer_T)read() {
	for {
		data := make([]byte, 512, 512)
		n, remoteAddr, err := peer.Conn.ReadFromUDP(data)
		if err != nil {
			log.Println("error during read:", err)
		}

		if n < 4 {
			log.Println("Invalid q2p header length")
			continue
		}

		peer.networking(remoteAddr, data[:n])
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

func (peer *Peer_T)touchRequest(rAddr3 *net.UDPAddr) {
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

func (peer *Peer_T)touch(rAddr, rAddr3 *net.UDPAddr) {
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

func (peer *Peer_T)connectRequest(rAddr2, rAddr3 *net.UDPAddr) {
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

func (peer *Peer_T)connect(rAddr2 *net.UDPAddr) {
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

func (peer *Peer_T)connected(rAddr3 *net.UDPAddr) {
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

func (peer *Peer_T)connectFailed(rAddr *net.UDPAddr, message []byte) {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, CONNECT_FAILED)
	header = append(header, bs...)

	data := append(header, message...)

	_, err := peer.Conn.WriteToUDP(data, rAddr)
	if err != nil {
		log.Println(err)
	}
}

// Transfer a single packet to the UDP address: `rAddr`.
// It usually used to resend a lost packet.
// `key` is the transmission's HASH that the lost packet is in.
// `syn` is the packet's SYN.
// `body` is the body for resend.
func (peer *Peer_T)TransportAPacket(rAddr *net.UDPAddr, key string, syn uint32, body []byte) error {
	log.Println("send again:", key, syn)
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, TRANSPORT)
	header = append(header, bs...)

	hash, err := hex.DecodeString(key)
	if err != nil {
		return err
	}

	bs = make([]byte, 4, 4)
	transmissionHead := make([]byte, 0, 24)
	transmissionHead = append(transmissionHead, hash[:]...) // 0th ~ 16th bytes: hash
	transmissionHead = append(transmissionHead, []byte{0, 0, 0, 0}...) // 16th ~ 20th bytes: length always = 0
	binary.LittleEndian.PutUint32(bs, syn)
	transmissionHead = append(transmissionHead, bs...) // 20nd ~ 24th bytes: SYN

	transm := append(transmissionHead, body...)
	transmission := append(header, transm...)
	_, err = peer.Conn.WriteToUDP(transmission, rAddr)
	if err != nil {
		return err
	}

	return nil
}

// Transfer `data` to the UDP address: `rAddr`.
// If the length of data > 2078764170780 (2078764170780 = math.MaxUint32 * 484, 484 is each packet's body length),
// it'll return an empty string and an error. Or the transmission HASH and nil. 
func (peer *Peer_T)Transport(rAddr *net.UDPAddr, data []byte) (string, error) {
	header := make([]byte, 0, 4)
	bs := make([]byte, 2, 2)
	log.Println("NetworkID:", peer.NetworkID)
	binary.LittleEndian.PutUint16(bs, peer.NetworkID)
	header = append(header, bs...)

	binary.LittleEndian.PutUint16(bs, TRANSPORT)
	header = append(header, bs...)

	hash := md5.Sum(data)

	key := fmt.Sprintf("%x", hash)

	length := len(data)
	if(length > 2078764170780) { // 2078764170780 = math.MaxUint32 * 484, 484 is each packet's body length
		return "", errors.New("too long data: should be less than 2078764170780")
	}
	packetNum := length / PACKET_LEN

	if length % PACKET_LEN != 0 {
		packetNum++
	}

	bs = make([]byte, 4, 4)

	transmissionHead := make([]byte, 0, 24)
	transmissionHead = append(transmissionHead, hash[:]...) // 0th ~ 16th bytes: hash
	binary.LittleEndian.PutUint32(bs, uint32(length))
	transmissionHead = append(transmissionHead, bs...)  // 16th ~ 20nd bytes: length
	transmissionHead = append(transmissionHead, []byte{0, 0, 0, 0}...) // 20nd ~ 24th bytes: leave space empty for each SYN

	for i := 0; i < packetNum; i++ {
		// for packet losing test
		/*
		if i == 3 || i == 5 {
			continue
		}
		*/

		start := i * PACKET_LEN
		end := start + PACKET_LEN
		if(length < end) {
			end = length
		}

		body := data[start: end]

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

	ctx, _ := context.WithTimeout(context.TODO(), time.Second * time.Duration(peer.Timeout))
	go transmissionSending(ctx, key, rAddr.String())

	return key, nil

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
