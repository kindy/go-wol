package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/sabhiram/go-wol/wol"
)

func main() {

	fmt.Printf("%+v", os.Args)

	if len(os.Args) != 5 {
		fmt.Printf("%s addr certFile keyfile macAddr\n", os.Args[0])
		os.Exit(1)
	}

	addr, certFile, keyFile, macAddr := os.Args[1], os.Args[2], os.Args[3], os.Args[4]

	if err := http.ListenAndServeTLS(addr, certFile, keyFile, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var err error

		defer (func() {
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Add("content-type", "text/plain")
				w.Write([]byte(fmt.Sprintf("fail: %s", err)))
			}
		})()

		if r.URL.Path != "/wake" {
			w.WriteHeader(http.StatusTeapot)
			return
		}

		// Build the magic packet.
		mp, err := wol.New(macAddr)
		if err != nil {
			return
		}

		// Grab a stream of bytes to send.
		bs, err := mp.Marshal()
		if err != nil {
			return
		}

		var localAddr *net.UDPAddr

		BroadcastIP := "255.255.255.255"
		UDPPort := "9"
		bcastAddr := fmt.Sprintf("%s:%s", BroadcastIP, UDPPort)
		udpAddr, err := net.ResolveUDPAddr("udp", bcastAddr)

		// Grab a UDP connection to send our packet of bytes.
		conn, err := net.DialUDP("udp", localAddr, udpAddr)
		if err != nil {
			return
		}
		defer conn.Close()

		fmt.Printf("Attempting to send a magic packet to MAC %s\n", macAddr)
		fmt.Printf("... Broadcasting to: %s\n", bcastAddr)
		n, err := conn.Write(bs)
		if err == nil && n != 102 {
			err = fmt.Errorf("magic packet sent was %d bytes (expected 102 bytes sent)", n)
		}
		if err != nil {
			return
		}

		fmt.Printf("Magic packet sent successfully to %s\n", macAddr)

	})); err != nil {
		panic(err)
	}

}
