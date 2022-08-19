package main

import (
	"log"
	"math"
	"net"
	"time"
)

type Relay struct {
	bindAddr   *net.UDPAddr
	serverAddr *net.UDPAddr
	listener   *net.UDPConn
}

func NewRelay(bindAddr, serverAddr string) *Relay {
	bind, err := net.ResolveUDPAddr("udp", bindAddr)
	if err != nil {
		return nil
	}

	server, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil
	}

	return &Relay{bindAddr: bind, serverAddr: server}
}

const (
	maxPacketLength = math.MaxUint16 - 28
	timeoutDuration = 15 * time.Second
)

func (r *Relay) Serve() error {
	listener, err := net.ListenUDP("udp", r.bindAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	r.listener = listener
	log.Printf("Serving on %v\n", r.bindAddr)
	log.Printf("the provided server address is %v\n", r.serverAddr)

	var buf [maxPacketLength]byte
	for {
		n, senderAddr, err := listener.ReadFromUDP(buf[:])
		if err != nil {
			return err
		}
		packet := make([]byte, n)
		copy(packet, buf[:n])
		go r.handleClient(senderAddr, packet)
	}
}

func (r *Relay) handleClient(clientAddr *net.UDPAddr, packet []byte) {
	serverConn, err := net.Dial("udp", r.serverAddr.String())
	if err != nil {
		log.Println("could not connect to the dns server")
		return
	}
	defer serverConn.Close()

	_, err = serverConn.Write(packet)
	if err != nil {
		log.Println("server write packet failed:", err)
		return
	}

	var buf [maxPacketLength]byte
	serverConn.SetReadDeadline(time.Now().Add(timeoutDuration))
	n, err := serverConn.Read(buf[:])
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			log.Println("timeout:", err)
			return
		}
		log.Println("could not read")
		return
	}

	_, err = r.listener.WriteToUDP(buf[:n], clientAddr)
	if err != nil {
		log.Println("could not write packet:", err)
	}
}
