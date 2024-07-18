package q2p

import (
	"testing"
	"syscall"
	"os"
	"os/signal"
	"flag"
)

type cmdFlag_T struct {
	ip string
	port int
	remoteHost string
	networkIdentify string
}

var cmdFlag *cmdFlag_T

func init() {
	cmdFlag = &cmdFlag_T{}
	readFlags(cmdFlag)
	cmdFlag.networkIdentify = "Hello"
}

func TestQ2P(t *testing.T) {
	t.Log(*cmdFlag)

	params := &q2pParams_T {
		cmdFlag.ip,
		cmdFlag.port,
		cmdFlag.remoteHost,
		cmdFlag.networkIdentify,
	}
	err := Q2P(params)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	t.Log("Received signal, shutting down...")
}

func readFlags(cmdFlag *cmdFlag_T) {
	flag.StringVar(&cmdFlag.ip, "ip", "0.0.0.0", "UDP host IP")
	flag.IntVar(&cmdFlag.port, "port", 10000, "UDP host Port")
	flag.StringVar(&cmdFlag.remoteHost, "remote_host", "", "remote host address")
}
