package q2p

import (
	"testing"
	"syscall"
	"os"
	"os/signal"
	"flag"
)

type cmdFlag_T struct {
	IP string
	Port int
}

var cmdFlag *cmdFlag_T

func init() {
	cmdFlag = &cmdFlag_T{}
	readFlags(cmdFlag)
}

func TestQ2P(t *testing.T) {
	t.Log(*cmdFlag)

	err := Q2P(cmdFlag.IP, cmdFlag.Port)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	t.Log("Received signal, shutting down...")
}

func readFlags(cmdFlag *cmdFlag_T) {
        flag.StringVar(&cmdFlag.IP, "ip", "0.0.0.0", "UDP host IP")
	flag.IntVar(&cmdFlag.Port, "port", 10000, "UDP host Port")
}
