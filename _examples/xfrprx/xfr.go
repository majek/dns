package main

// Xfrprx is a proxy that intercepts notify messages
// and then performs a ixfr/axfr to get the new 
// zone contents. 
// This zone is then checked cryptographically is
// everything is correct.
// If a new DNSKEY record is seen for the apex and
// it validates it writes this record to disk and
// this new key will be used in future validations.

import (
        "os"
        "os/signal"
	"fmt"
	"dns"
)

func reply(d *dns.Conn, i *dns.Msg) []byte {
        return nil
}

func handle(d *dns.Conn, i *dns.Msg) {
        if i.MsgHdr.Response == true {
                return
        }
        answer := reply(d, i)
        d.Write(answer)
}

func listen(addr string, e chan os.Error, tcp string) {
        switch tcp {
        case "tcp":
                err := dns.ListenAndServeTCP(addr, handle)
                e <- err
        case "udp":
                err := dns.ListenAndServeUDP(addr, handle)
                e <- err
        }
        return
}

func main() {
	err := make(chan os.Error)
	go listen("127.0.0.1:8053", err, "tcp")
	go listen("[::1]:8053", err, "udp")
	go listen("127.0.0.1:8053", err, "tcp")
	go listen("[::1]:8053", err, "udp")

forever:
	for {
		select {
		case e := <-err:
			fmt.Printf("Error received, stopping: %s\n", e.String())
			break forever
		case <-signal.Incoming:
			fmt.Printf("Signal received, stopping")
			break forever
		}
	}
	close(err)

}