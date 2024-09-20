package main

import (
	"context"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	// Create a custom resolver that doesn't use the OS's DNS cache
	resolver := &net.Resolver{
		// PreferGo controls whether Go's built-in DNS resolver is preferred
		// on platforms where it's available. It is equivalent to setting
		// GODEBUG=netdns=go, but scoped to just this resolver.
		PreferGo: true,
		// PreferGo controls whether Go's built-in DNS resolver is preferred
		// on platforms where it's available. It is equivalent to setting
		// GODEBUG=netdns=go, but scoped to just this resolver.
		StrictErrors: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			// Create a dialer with a timeout
			d := net.Dialer{
				Timeout: 5 * time.Second,
			}
			// Dial the DNS server and return the connection
			conn, err := d.DialContext(ctx, "udp", "8.8.8.8:53")
			if err != nil {
				return nil, err
			}

			// Log the source port of the connection
			localAddr := conn.LocalAddr().(*net.UDPAddr)
			log.Printf("Source port for DNS lookup: %d", localAddr.Port)

			return conn, nil
		},
	}

	for {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Perform the DNS lookup using our custom resolver
			addr, err := resolver.LookupHost(ctx, "google.com")
			if err == nil {
				log.Println("DNS lookup successful: " + strings.Join(addr, ","))
			} else {
				if err, ok := err.(net.Error); ok && err.Timeout() {
					log.Println("DNS lookup timed out")
				} else {
					log.Printf("DNS lookup error: %v", err)
				}
			}
		}()

		time.Sleep(1 * time.Second)
	}
}
