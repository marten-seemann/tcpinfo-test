package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/mikioh/tcp"
	"github.com/mikioh/tcpinfo"
)

func main() {
	connChan := make(chan net.Conn, 1)
	r := http.Transport{
		DialTLSContext: func(_ context.Context, network, addr string) (net.Conn, error) {
			c, err := net.Dial(network, addr)
			if err != nil {
				return nil, err
			}
			conn := tls.Client(c, &tls.Config{InsecureSkipVerify: true})
			conn.Handshake()
			if err != nil {
			  return nil, err
			}
			connChan <- c
			return conn, nil
		},
	}

	req, err := http.NewRequest(http.MethodGet, "https://dl.google.com/go/go1.15.7.src.tar.gz", nil)
	if err != nil {
		log.Fatal(err)
	}
	rsp, err := r.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
	  log.Fatal(err)
	}
	log.Printf("Read %d kb.\n", len(data)/1024)

	c := <-connChan
	tc, err := tcp.NewConn(c)
	if err != nil {
		log.Fatal(err)
	}
	printInfo(tc)
}

func printInfo(tc *tcp.Conn) {
	var o tcpinfo.Info
	var b [256]byte
	i, err := tc.Option(o.Level(), o.Name(), b[:])
	if err != nil {
		log.Fatal(err)
	}
	txt, err := json.Marshal(i)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(txt))
}
