package q2p

import (
	"net"
	"log"
)

func callback(peer *peer_T, rAddr *net.UDPAddr, event uint16, data []byte) error {
	log.Println("event:", event)
	switch event {
	case JOIN:
		log.Println("event: JOIN")
		log.Println("remote seeds")
		for seed, _ := range peer.RemoteSeeds {
			log.Println(seed.String())
		}
		peer.TouchRequest(rAddr)
		peer.RemoteSeeds[rAddr] = false
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
		peer.RemoteSeeds[rAddr] = false

		peer.Connected(rAddr)
	case CONNECTED:
		log.Println("event: CONNECTED")
		log.Println("from:", rAddr.String())
		peer.RemoteSeeds[rAddr] = false
	case TEST:
		log.Println("event: TEST")
	default:
		log.Println(event)
		log.Println("Undefined event")
	}

	return nil
}
