package q2p

import (
	"encoding/binary"
	"context"
	"net"
	"time"
	"fmt"
	"log"
)

func (peer *Peer_T)networking(rAddr *net.UDPAddr, data []byte) {
	networkID := binary.LittleEndian.Uint16(data[:2])
	event := binary.LittleEndian.Uint16(data[2:4])

	log.Println("Remote networkID =", networkID)
	log.Println("event =", event)

	if networkID != peer.NetworkID && event != CONNECT_FAILED {
		peer.connectFailed(rAddr, []byte("Dismatch networking version"))
		return
	}

	switch event {
	case JOIN:
		log.Println("event: JOIN")

		peer.touchRequest(rAddr)
		if(len(peer.RemoteSeeds) < connection_num) {
			peer.RemoteSeeds[rAddr.String()] = false
		}

		peer.LifeCycle(peer, rAddr, int(event))
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

		peer.touch(rAddr, rAddr3)
	case TOUCH:
		log.Println("event: TOUCH")
	case TOUCHED:
		log.Println("event: TOUCHED")
		log.Println("rAddr3:", string(data[4:]))
		rAddr3, err := net.ResolveUDPAddr("udp", string(data[4:]))
		if err != nil {
			log.Println(err)
		}

		peer.connectRequest(rAddr, rAddr3) // here rAddr is rAddr2
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

		peer.connect(rAddr2)
	case CONNECT:
		log.Println("event: CONNECT")
		log.Println("from:", rAddr.String())
		peer.RemoteSeeds[rAddr.String()] = false

		peer.connected(rAddr)
		peer.LifeCycle(peer, rAddr, int(event))
	case CONNECTED:
		log.Println("event: CONNECTED")
		log.Println("from:", rAddr.String())
		peer.RemoteSeeds[rAddr.String()] = false
		peer.LifeCycle(peer, rAddr, int(event))
	case CONNECT_FAILED:
		log.Println("event: CONNECT_FAILED")
		log.Println("from:", rAddr.String())
		log.Println("Network error:", rAddr.String(), string(data[4:]))
	case TRANSPORT:
		event := binary.LittleEndian.Uint16(data[2:4])

		log.Println("event: TRANSPORT", event)
		log.Println("from:", rAddr.String())

		hash := data[4:20]
		length := binary.LittleEndian.Uint32(data[20:24])
		syn := binary.LittleEndian.Uint32(data[24:28])

		key := fmt.Sprintf("%x", hash)

		_, ok := transmissionR[key]
		if !ok {
			transmissionR[key] = make([]byte, length, length)

			transmissionCTXM[key] = &ctx_T{}
			transmissionCTXM[key].ctx, transmissionCTXM[key].cancel = context.WithTimeout(context.TODO(), time.Second * time.Duration(peer.Timeout))

			go transmissionReceiving(transmissionCTXM[key].ctx, peer, hash, rAddr.String())

			packetNum := length / 484
			if length % 484 != 0 {
				packetNum += 1
			}
			transmissionRSYNS[key] = make(map[uint32]bool)

			for i := 0; i < int(packetNum); i++ {
				transmissionRSYNS[key][uint32(i)] = false
			}
		}

		/*
		fmt.Println("length:", length)
		fmt.Println("syn:", syn)
		fmt.Println("length:", len(data[28:]))
		*/

		start := int(syn) * 484
		end := int(start) + len(data[28:])

		/*
		if len(transmissionR[key]) == 0 {
			break
		}
		*/

		copy(transmissionR[key][start:end], data[28:])

		delete(transmissionRSYNS[key], syn)

		if len(transmissionRSYNS[key]) == 0 {
			go peer.Successed(peer, rAddr, key, transmissionR[key])
			transmissionCTXM[key].cancel()
		}

	case TRANSPORT_FAILED:
		event := binary.LittleEndian.Uint16(data[2:4])
		log.Println("event: TRANSPORT FAILED", event)
		log.Println("from:", rAddr.String())

		hash := data[4:20]
		key := fmt.Sprintf("%x", hash)

		lengthOfSyns := (len(data) - 20) / 4
		syns := make([]uint32, 0, lengthOfSyns)
		for i := 0; i < lengthOfSyns; i++ {
			start := 20 + i * 4
			end := start + 4
			syn := binary.LittleEndian.Uint32(data[start:end])
			syns = append(syns, syn)
		}

		peer.Failed(peer, rAddr, key, syns)
	case TEST:
		log.Println("event: TEST")
	default:
		log.Println(event)
		log.Println("Undefined event")
	}
}
