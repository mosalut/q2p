package q2p

import (
	"sync"
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

var transmissionS = make(map[string][]byte) // Senders' data
var transmissionR = make(map[string][]byte) // Receivers' data
var transmissionRSYNS = make(map[string]map[uint32]bool) // Receivers' data SYNs
var transmissionCTXM = make(map[string]*ctx_T)

var mutex = &sync.RWMutex{}

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
			mutex.Lock()
			if len(transmissionRSYNS[key]) != 0 {
				log.Println("transport failed")

				rAddr, err := net.ResolveUDPAddr("udp", addr)
				if err != nil {
					log.Println(err)
					log.Println("transmissionReceiving Done timeout")
					mutex.Unlock()
					return
				}

				peer.transportFailed(rAddr, hash, []uint32{})
			}
			delete(transmissionR, key)
			delete(transmissionRSYNS, key)
			delete(transmissionCTXM, key)
			mutex.Unlock()
			log.Println(key, addr, "recieving over")
			log.Println(ctx.Err())

			log.Println("transmissionReceiving Done")
			return
		default:
			time.Sleep(time.Second * time.Duration(peer.TimeSendLost))
			log.Println("RSYNS length:", len(transmissionRSYNS[key]))
			mutex.Lock()
			if len(transmissionRSYNS[key]) != 0 {
				log.Println("Packet lost")

				syns := make([]uint32, 0, len(transmissionRSYNS[key]))
				for k, _ := range transmissionRSYNS[key] {
					syns = append(syns, k)
				}

				rAddr, err := net.ResolveUDPAddr("udp", addr)
				if err != nil {
					log.Println(err)
					mutex.Unlock()
					return
				}

				peer.transportFailed(rAddr, hash, syns)
			}
			mutex.Unlock()
		}
	}
}
