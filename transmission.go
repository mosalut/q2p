package q2p

import (
	"net"
	"context"
	"fmt"
	"log"
)

// pool of transmission, string is a MD5 hash string, uint32 is a SYN
var transmissionM = make(map[string]map[uint32][]byte) // Senders' data
var transmissionR = make(map[string][]byte) // Receiver' data
var transmissionRSYNS = make(map[string]map[uint32]bool) // Receiver' data

func transmissionSending(ctx context.Context, key, addr string) {
	select {
	case <-ctx.Done():
		delete(transmissionM, key)
		log.Println(key, addr, "transport over")
		log.Println(ctx.Err())
	}
}

func transmissionReceiving(ctx context.Context, peer *Peer_T, hash []byte, addr string) {
	key := fmt.Sprintf("%x", hash)

	select {
	case <-ctx.Done():
		if len(transmissionRSYNS[key]) != 0 {
			log.Println("transport failed")

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
		delete(transmissionR, key)
		log.Println(key, addr, "recieving over")
		log.Println(ctx.Err())
	}
}
