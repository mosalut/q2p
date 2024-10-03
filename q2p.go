package q2p

import (
	"net"
	"log"
)

type Header_T struct {
	RAddr *net.UDPAddr `json:"raddr"`
	PackageNum int `json:"package_num"`
	Hash [32]byte `json:"hash"`
	Timestamp int64 `json:"timestamp"`
}

type Request_T struct {
	Header *Header_T `json:"header"`
	Body []byte `json:"body"`
	Status int `json:"status"` // 0: header, 1: receiving, 2: received
}

func (peer *peer_T)networking(rAddr *net.UDPAddr, event uint16, data []byte) error {
	log.Println("event:", event)
	switch event {
	case JOIN:
		log.Println("event: JOIN")

		peer.TouchRequest(rAddr)
		if(len(peer.RemoteSeeds) < connection_num) {
			peer.RemoteSeeds[rAddr.String()] = false
		}
	case TOUCHREQUEST:
		log.Println("event: TOUCHREQUEST")
		if(len(peer.RemoteSeeds) >= connection_num) {
			log.Println("The connection number upper limit has been reached")
			break
		}
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
		if(len(peer.RemoteSeeds) >= connection_num) {
			log.Println("The connection number upper limit has been reached")
			break
		}
		log.Println("rAddr2:", string(data[4:]))
		rAddr2, err := net.ResolveUDPAddr("udp", string(data[4:]))
		if err != nil {
			log.Println(err)
		}

		peer.Connect(rAddr2)
	case CONNECT:
		log.Println("event: CONNECT")
		log.Println("from:", rAddr.String())
		peer.RemoteSeeds[rAddr.String()] = false

		peer.Connected(rAddr)
	case CONNECTED:
		log.Println("event: CONNECTED")
		log.Println("from:", rAddr.String())
		peer.RemoteSeeds[rAddr.String()] = false
	case HEADER:

	case BODY:

	case TEST:
		log.Println("event: TEST")
	default:
		log.Println(event)
		log.Println("Undefined event")
	}

	return nil
}
