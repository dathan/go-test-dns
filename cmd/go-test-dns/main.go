package main

import (
	"context"
	"log"
	"net"
	"time"
)

func main() {
	for {
		func() {
			conn, err := net.Dial("udp", "8.8.8.8:53")
			if err != nil {
				log.Printf("Error dialing UDP: %v", err)
				return
			}
			defer conn.Close()

			udpConn, ok := conn.(*net.UDPConn)
			if !ok {
				log.Println("Error: Not a UDP connection")
				return
			}

			localAddr := udpConn.LocalAddr().(*net.UDPAddr)
			log.Printf("Source port for DNS lookup: %d", localAddr.Port)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var r net.Resolver
			_, err = r.LookupHost(ctx, "google.com")
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Timeout() {
					log.Println("DNS lookup timed out")
				} else {
					log.Printf("DNS lookup error: %v", err)
				}
			} else {
				log.Println("DNS lookup successful")
			}
		}()

		time.Sleep(1 * time.Second)
	}
}
