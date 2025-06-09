package q2p

import (
	"encoding/binary"
	"encoding/hex"
	"context"
	"net"
	"time"
)

func (peer *Peer_T)networking(rAddr *net.UDPAddr, data []byte) {
	networkID := binary.LittleEndian.Uint16(data[:2])
	event := binary.LittleEndian.Uint16(data[2:4])

	print(log_info, "Remote networkID =", networkID)
	print(log_info, "event =", event)

	if networkID != peer.NetworkID && event != CONNECT_FAILED {
		peer.connectFailed(rAddr, []byte("Dismatch networking version"))
		return
	}

	switch event {
	case JOIN:
		print(log_info, "event: JOIN")

		peer.touchRequest(rAddr)
		if(len(peer.RemoteSeeds) < connection_num) {
			peer.RemoteSeeds[rAddr.String()] = false
		}

		peer.LifeCycle(peer, rAddr, int(event))
	case TOUCHREQUEST:
		print(log_info, "event: TOUCHREQUEST")
		if(len(peer.RemoteSeeds) >= connection_num) {
			print(log_warning, "The connection number upper limit has been reached")
			break
		}
		print(log_info, "rAddr3:", string(data[4:]))
		rAddr3, err := net.ResolveUDPAddr("udp", string(data[4:]))
		if err != nil {
			print(log_error, err)
		}

		peer.touch(rAddr, rAddr3)
	case TOUCH:
		print(log_info, "event: TOUCH")
	case TOUCHED:
		print(log_info, "event: TOUCHED")
		print(log_info, "rAddr3:", string(data[4:]))
		rAddr3, err := net.ResolveUDPAddr("udp", string(data[4:]))
		if err != nil {
			print(log_error, err)
		}

		peer.connectRequest(rAddr, rAddr3) // here rAddr is rAddr2
	case CONNECTREQUEST:
		print(log_info, "event: CONNECTREQUEST")
		if(len(peer.RemoteSeeds) >= connection_num) {
			print(log_warning, "The connection number upper limit has been reached")
			break
		}
		print(log_info, "rAddr2:", string(data[4:]))
		rAddr2, err := net.ResolveUDPAddr("udp", string(data[4:]))
		if err != nil {
			print(log_error, err)
		}

		peer.connect(rAddr2)
	case CONNECT:
		print(log_info, "event: CONNECT")
		print(log_info, "from:", rAddr.String())
		peer.RemoteSeeds[rAddr.String()] = false

		peer.connected(rAddr)
		peer.LifeCycle(peer, rAddr, int(event))
	case CONNECTED:
		print(log_info, "event: CONNECTED")
		print(log_info, "from:", rAddr.String())
		peer.RemoteSeeds[rAddr.String()] = false
		peer.LifeCycle(peer, rAddr, int(event))
	case CONNECT_FAILED:
		print(log_error, "event: CONNECT_FAILED")
		print(log_error, "from:", rAddr.String())
		print(log_error, "Network error:", rAddr.String(), string(data[4:]))
	case OPTIONS:
		if len(data) != 24 {
			break
		}

		event := binary.LittleEndian.Uint16(data[2:4])

		print(log_info, "event: OPTIONS", event)
		print(log_info, "from:", rAddr.String())

		hash := data[4:20]
		length := binary.LittleEndian.Uint32(data[20:24])

		key := hex.EncodeToString(hash)

		// TODO issue?
		mutex.Lock()
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
		mutex.Unlock()

		peer.optionsFeedback(rAddr, hash)

	case OPTIONSFEEDBACK:
		if len(data) != 20 {
			break
		}

		event := binary.LittleEndian.Uint16(data[2:4])

		print(log_info, "event: OPTIONSFEEDBACK", event)
		print(log_info, "from:", rAddr.String())

		hash := data[4:20]

		peer.transfer(rAddr, hash)

	case TRANSPORT:
		if len(data) < 24 {
			break
		}

		event := binary.LittleEndian.Uint16(data[2:4])

		print(log_info, "event: TRANSPORT", event)
		print(log_info, "from:", rAddr.String())

		hash := data[4:20]
		syn := binary.LittleEndian.Uint32(data[20:24])

		key := hex.EncodeToString(hash)

		/*
		print(log_debug, "data length:", len(data))
		print(log_debug, "syn:", syn)
		print(log_debug, "length:", len(data[24:]))
		*/

		start := int(syn) * 484
		end := int(start) + len(data[24:])

		mutex.Lock()
		if len(transmissionR[key]) < end {
			mutex.Unlock()
			break
		}
		copy(transmissionR[key][start:end], data[24:])

		delete(transmissionRSYNS[key], syn)

		if len(transmissionRSYNS[key]) == 0 {
			go peer.Successed(peer, rAddr, key, transmissionR[key])
			transmissionCTXM[key].cancel()
		}
		mutex.Unlock()

	case TRANSPORT_FAILED:
		event := binary.LittleEndian.Uint16(data[2:4])
		print(log_error, "event: TRANSPORT FAILED", event)
		print(log_error, "from:", rAddr.String())

		hash := data[4:20]
		key := hex.EncodeToString(hash)

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
		print(log_info, "event: TEST")
	default:
		print(log_warning, event)
		print(log_warning, "Undefined event")
	}
}
