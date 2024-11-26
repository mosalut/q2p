package q2p

import (
	"encoding/binary"
	"context"
	"net"
	"time"
	"errors"
	"fmt"
	"log"
)

func (peer *Peer_T)networking(rAddr *net.UDPAddr, event uint16, data []byte, callback func(data []byte)) error {
	version := binary.LittleEndian.Uint16(data[:2])
	log.Println("version:", version)
	if version != 0 {
		return errors.New("Dismatch networking version")
	}

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
	case TRANSPORT:
		event := binary.LittleEndian.Uint16(data[2:4])

		log.Println("event: TRANSPORT", event)
		log.Println("from:", rAddr.String())

		hash := data[4:20]
		length := binary.LittleEndian.Uint32(data[20:24])
		syn := binary.LittleEndian.Uint32(data[24:28])

		key := fmt.Sprintf("%x", hash)
		fmt.Println("hash:", key)

		_, ok := transmissionR[key]
		if !ok {
			transmissionR[key] = make([]byte, length, length)

			transmissionCTXM[key] = &ctx_T{}
			transmissionCTXM[key].ctx, transmissionCTXM[key].cancel = context.WithTimeout(context.TODO(), time.Second * 10)
			go transmissionReceiving(transmissionCTXM[key].ctx, peer, hash, rAddr.String())

			packetNum := length / 484
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

		copy(transmissionR[key][start:end], data[28:])

		delete(transmissionRSYNS[key], syn)

		if(len(transmissionRSYNS[key]) == 0) {
			go callback(transmissionR[key])
			transmissionCTXM[key].cancel()
		}

	case TRANSPORT_FAILED:
		event := binary.LittleEndian.Uint16(data[2:4])
		log.Println("event: TRANSPORT FAILED", event)
		log.Println("from:", rAddr.String())

		hash := data[4:20]
		fmt.Printf("%x failed\n", hash)

		lengthOfSyns := (len(data) - 20) / 4
		syns := make([]uint32, 0, lengthOfSyns)
		for i := 0; i < lengthOfSyns; i++ {
			start := 20 + i * 4
			end := start + 4
			syn := binary.LittleEndian.Uint32(data[start:end])
			syns = append(syns, syn)
		}
		fmt.Println(syns)
	case TEST:
		log.Println("event: TEST")
	default:
		log.Println(event)
		log.Println("Undefined event")
	}

	return nil
}
