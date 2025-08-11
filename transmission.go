package q2p

import (
	"encoding/hex"
	"sync"
	"net"
	"context"
	"time"
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
	print(log_info, "transmissionSending called")
	select {
	case <-ctx.Done():
		print(log_info, "transmissionSending Done")
		print(log_error, ctx.Err())
		print(log_info, key, addr, "transport over")
	}
}

func transmissionReceiving(ctx context.Context, peer *Peer_T, hash []byte, addr string) {
	print(log_info, "transmissionReceiving called")

	key := hex.EncodeToString(hash)

	for {
		select {
		case <-ctx.Done():
			mutex.Lock()
			if len(transmissionRSYNS[key]) != 0 {
				print(log_warning, "transport failed")

				rAddr, _ := net.ResolveUDPAddr("udp", addr)

				peer.transportFailed(rAddr, hash, []uint32{})
			}
			delete(transmissionR, key)
			delete(transmissionRSYNS, key)
			delete(transmissionCTXM, key)
			mutex.Unlock()
			print(log_info, key, addr, "recieving over")
			print(log_error, ctx.Err())

			print(log_info, "transmissionReceiving Done")
			return
		default:
			time.Sleep(time.Second * time.Duration(peer.TimeSendLost))
			mutex.Lock()
			print(log_debug, "RSYNS length:", len(transmissionRSYNS[key]))
			if len(transmissionRSYNS[key]) != 0 {
				print(log_warning, "Packet lost")

				syns := make([]uint32, 0, len(transmissionRSYNS[key]))
				for k, _ := range transmissionRSYNS[key] {
					syns = append(syns, k)
				}

				rAddr, _ := net.ResolveUDPAddr("udp", addr)

				peer.transportFailed(rAddr, hash, syns)
			}
			mutex.Unlock()
		}
	}
}
