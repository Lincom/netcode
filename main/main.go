package main

import (
	"fmt"
	"github.com/Lincom/netcode"
	"github.com/ogier/pflag"
	"os"
	"os/signal"
	"syscall"
)

const (
	VERSION = 0
)

var helpFlag bool
var locFlag string
var typeFlag string

func init() {
	pflag.BoolVar(&helpFlag, "help", false, "Show the help menu.")
	pflag.StringVar(&locFlag, "loc", ":8080", "Listen from the specified location (default port 8080 all interfaces).")
	pflag.StringVar(&typeFlag, "type", "tcp", "Listen on the specified connection type (default is tcp).")
}

func main() {
	pflag.Parse()

	if helpFlag {
		fmt.Printf("netcode v%v : Apache License, Robert Xu\n", VERSION)
		pflag.PrintDefaults()
		fmt.Printf("\nGood connection types: http://golang.org/pkg/net/#Listen\n")
		return
	}

	fmt.Printf("netcode v%v : starting on %s port %s...\n", VERSION, typeFlag, locFlag)
	go netcode.Listen(typeFlag, locFlag)

	terminateSignal := make(chan os.Signal)
	signal.Notify(terminateSignal, syscall.SIGINT)

	fmt.Printf("ctrl-c to terminate.\n")

	<-terminateSignal
	fmt.Printf("caught ctrl-c, stopping listener...\n")
	fmt.Printf("this program will fully terminate upon disconnection of remaining clients.\n\n")
	netcode.Stop()

	fmt.Printf("\nnetcode: closing and exiting\n")
}
