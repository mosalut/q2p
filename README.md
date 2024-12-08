# Q2P

`Q2P` is a peer-to-peer (P2P) network communication library implemented based on the UDP protocol. It aims to support connection and data exchange between multiple nodes through a simple network event-driven and data transmission protocol. This library is written in Golang and provides a simple way to implement device connection, communication and data transmission in P2P networks.

`Q2P` is divided into two layers: Protocol Layer and Transport Layer.

`Q` stands for quick, because its implementation method and characteristics of the transport layer are very similar to the QUIC protocol.

## Features

- **Network Event-Driven**: Supports handling various network events (such as connection requests, data transmission, etc.).
- **Peer-to-Peer Communication**: Conducts efficient data transmission and node connections based on the UDP protocol.
- **Reliable Connection Management**: Implements processes such as connection requests and confirmations between nodes through an event mechanism.
- **Data Fragmentation Transmission**: Supports fragmentation transmission of large data volumes and retries in case of transmission failures.
- **Transmission Timeout Control**: Supports timeout detection and failure handling during transmission.

### Environmental Requirements

- Go 1.22

### Install Dependencies
```bash
go get github.com/mosalut/q2p
```

## Quick Start

### Running
Assuming a seed node has been started at 127.0.0.1:10000

```go
package main

import (
	"github.com/mosalut/q2p"
	"log"
)

/* Custom Callback Function
	peer *Peer_T	// Current peer
	rAddr *net.UDPAddr	// Sender's address
	key string	// Transmission HASH
	data []byte	// The complete data received
*/
func callback(peer *Peer_T, rAddr *net.UDPAddr, key string, data []byte) {
	fmt.Println(key) // key is the hash of this transmission
	fmt.Println(string(data))
}

/* Custom callback function for transport failure
	peer *Peer_T	// Current peer
	rAddr *net.UDPAddr	// Sender's address
	key string	// Transmission HASH
	syns []uint32	// List of lost SYN positions
*/
func callbackFailed(peer *q2p.Peer_T, rAddr *net.UDPAddr, key string, syns []uint32) { // rAddr is the UDP address of the node that sends the TRANSPORT_FAILED event, and it is also the address that receives the TRANSPORT event.
	fmt.Println(key) // key is the hash of the failed transmission
	fmt.Println(syns []uint32) // The SYN of the lost packets
}

func main() {
	seedAddrs = make(map[string]bool)
	seedAddrs[127.0.0.1:10000] = false // Put the known seed nodes into seedAddrs

        //    Create a new node
	//    Parameter list:
	//	ip string: Host address where the node starts. The default is 0.0.0.0
	//	port int: Port where the node starts. The default is 10000.
	//	rAddrs map[string]bool: List of seed nodes. The default is empty.
	//	networkID uint16: Version number
	//	timeSendAgain int: How often the receiver checks for lost packets. If there are lost packets, it will inform the sender. The default value is 5(SECs).
	//	timeout int: Timeout for the receiver to wait for complete data. If it times out, it will inform the sender. The default value is 5(SECs).
	// 	Callback function executed when complete data is received. The default function member is to type success transmission HASH.
	// 	Callback function executed upon failure. If the length of the last parameter is 0, it indicates a timeout; otherwise, it represents the SYN position where the packet was lost. The default function member is to type timeout transmission HASH or lost transmission HASH that lost packet in with SYNs' positions.
	peer := q2p.NewPeer("127.0.0.1", 10001, seedAddrs, 0})
	err := peer.Run()
	if err != nil {
		log.Fatal(err)
	}
}
```
- Create a peer by executing function q2p.NewPeer(ip, port, seedAddrs, networkID) .
- Run the peer by its 'Run' function.

The parameter ip, port are the UDP addresses of the IP and port where the local node starts the q2p network host.
The parameter seedAddrs is a map, each key of which is a UDP address of a seed node.
The parameter networkID is a network ID number, which is 2 bytes and uses 0 here.
The parameter timeSendLost is a time in seconds, which is used to set how long the data is not received completely and counted as packet loss. According to this time, the sending node can be notified to resend the lost data packet.
The parameter timeout is a time in seconds, which is used to set how long the data has not been received completely and counted as timeout, and notify the sending node that this transmission is timeout.
The parameter callback is a user-defined function used when receiving TRANSPORT events. The parameter data of the callback function is processed by the TRANSPORT event.
The parameter callbackFailed is a callback function when receiving failure.

### Data Transmission

Data transmission is achieved through the `Transport` method, which fragments the data and sends it via UDP to ensure reliable transmission of large amounts of data. If data transmission fails, the user can decide whether to attempt retransmission.

### Transmission Timeout and Failure Handling

In `transmission.go`, the `transmissionSending` and `transmissionReceiving` functions handle the sending and receiving timeouts of data transmission. If the complete data is not received within the specified time or an error occurs, the system records the transmission failure and triggers the corresponding retry or failure handling mechanism.

### Example Code

Here is a simple example :

```go
package main

import (
	"fmt"
	"log"
	"net"
	"github.com/mosalut/q2p"
)

var transmissionCache = make(map[string][]byte)

func callback(peer *Peer_T, rAddr *net.UDPAddr, key string, data []byte) {
	fmt.Println("Transmission hash:", key)
	fmt.Println("Received data:", string(data))
}

func callbackFailed(peer *q2p.Peer_T, rAddr *net.UDPAddr, key string, syns []uint32) {
	if len(syns) == 0 {
		fmt.Println("transmission Failed, hash:", key)
		return
	}

	data := transmissionCache[key]

	var start int
	var end int
	for _, v := range syns {
		start = v * q2p.PACKET_LEN
		if start + q2p.PACKET_LEN > len(data) {
			end = len(data)
		} else {
			end = start + q2p.PACKET_LEN
		}

		fmt.Println("lost SYN:", v, "data:", data[start:end])

		err := peer.TransportAPacket(rAddr, key, syn, data[start:end])
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	peer := q2p.NewPeer("127.0.0.1", 10000, nil, 0, callback, callbackFailed)

	err := peer.Run()
	if err != nil {
		log.Fatal(err)
	}

	// transport data
	remoteAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	data := []byte("Hello, P2P Network")
	key, err := !peer.Transport(remoteAddr, data)
	if err != nil {
		log.Fatal(err)
	}
	transmissionCache[key] = data
}
```

## Testing

This project provides some simple unit tests that cover the basic operation and data transmission functions of P2P nodes. You can run the tests with the following command:

```bash
go test -v -count 1 test.run TestQ2p 或 go test -v -count 1 test.run TestTransport // Start the initial seed node. The default host is 127.0.0.1:10000
go test -v -count 1 test.run TestTransport -remote_host 127.0.0.1:10000 -port 10001
go test -v -count 1 test.run TestTransport -remote_host 127.0.0.1:10000 -port 10002
```

### Test Files

 `q2p_test.go` file contains examples for testing node connections, event handling, and data transmission.

### Event Transmission Test

`TestQ2p` test function mainly reflects the testing of the protocol layer.

`TestTransport` test function simulates the whole process of data transmission, including data sending and receiving, and verifies the data fragmentation transmission and timeout control mechanism.

## Contribution

Submissions of issues and pull requests are welcome, and any suggestions for improvements and feature requests will be carefully considered. Please ensure that the code complies with Go language coding style and write necessary unit tests.

## License

This project is licensed under the **GNU General Public License v3.0**.

For detailed license content, please refer to the [`LICENSE`](LICENSE) file.

