package q2p

import (
	"net"
	"context"
	"time"
	"fmt"
	"log"
)

type ctx_T struct {
	ctx context.Context
	cancel context.CancelFunc
}

var transmissionR = make(map[string][]byte) // Receivers' data
var transmissionRSYNS = make(map[string]map[uint32]bool) // Receivers' data SYNs
var transmissionCTXM = make(map[string]*ctx_T)

func transmissionSending(ctx context.Context, key, addr string) {
	log.Println("transmissionSending called")
	select {
	case <-ctx.Done():
		log.Println("transmissionSending Done")
		log.Println(ctx.Err())
		log.Println(key, addr, "transport over")
	}
}

func transmissionReceiving(ctx context.Context, peer *Peer_T, hash []byte, addr string) {
	log.Println("transmissionReceiving called")

	key := fmt.Sprintf("%x", hash)

	for {
		select {
		case <-ctx.Done():
			if len(transmissionRSYNS[key]) != 0 {
				log.Println("transport failed")

				rAddr, err := net.ResolveUDPAddr("udp", addr)
				if err != nil {
					log.Println(err)
					log.Println("transmissionReceiving Done timeout")
					return
				}

				peer.transportFailed(rAddr, hash, []uint32{})
			}
			delete(transmissionR, key)
			delete(transmissionRSYNS, key)
			delete(transmissionCTXM, key)
			log.Println(key, addr, "recieving over")
			log.Println(ctx.Err())

			log.Println("transmissionReceiving Done")
			return
		default:
			time.Sleep(time.Second * time.Duration(peer.TimeSendLost))
			log.Println("RSYNS length:", len(transmissionRSYNS[key]))
			if len(transmissionRSYNS[key]) != 0 {
				log.Println("Packet lost")

				syns := make([]uint32, 0, len(transmissionRSYNS[key]))
				for k, _ := range transmissionRSYNS[key] {
					syns = append(syns, k)
				}

				rAddr, err := net.ResolveUDPAddr("udp", addr)
				if err != nil {
					log.Println(err)
					return
				}

				peer.transportFailed(rAddr, hash, syns)
			}
		}
	}
}
