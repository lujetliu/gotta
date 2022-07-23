package main

import (
	"log"
	"net"
	"time"
)

/*


 */

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("error listen:", err)
		return
	}

	defer l.Close()
	log.Println("listen ok")

	var i int
	for {
		time.Sleep(10 * time.Second)
		if _, err := l.Accept(); err != nil {
			log.Println("accept error:", err)
			break
		}
		i++
		log.Printf("%d: accept a new connection\n", i)
	}
}
