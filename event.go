package q2p // github.com/mosalut/q2p

// JOIN: When a node requests to join the network, the current node, as a seed node, will send a TOUCHREQUEST message to the already connected nodes to help the new node punch holes.
// TOUCHREQUEST: When a node requests to establish a connection, the seed node will send the UDP address of the target node, and the current node will then punch holes to that address and send a TOUCHED message to confirm to the seed node.
// TOUCH: This protocol mainly helps new nodes punch holes to establish connections through seed nodes, although the first punch hole may not receive a confirmation message due to route interception.
// TOUCHED: When a node receives a TOUCHED event, it will send a CONNECTREQUEST message to the punching target node to request a connection to the punching node.
// CONNECTREQUEST: After receiving the network event, the current node learns that the punching node has opened a port, and then it should send a CONNECT message to the punching node to establish a connection.
// CONNECT: When the punching node receives the connection request from the target node, it will reply with a CONNECTED message and add the target node to the connected list.
// CONNECTED: After receiving the CONNECTED message, the node will add the other party node to the connected list.
// TRANSPORT: Node starts data transmission.
// TRANSPORT_FAILED: Data transmission failed.
const (
	JOIN = iota
	TOUCHREQUEST
	TOUCH
	TOUCHED
	CONNECTREQUEST
	CONNECT
	CONNECTED
	CONNECT_FAILED
	TRANSPORT
	TRANSPORT_FAILED
	TEST
)
