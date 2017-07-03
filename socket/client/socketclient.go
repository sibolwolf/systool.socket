package main

import (
	"log"
	"net"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:1024")

	if err != nil {
		log.Println(err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Connet to server success")
	//str := "{\"Play\": \"/home/jinlai/0.ogg\"}"
	str := "{\"Play\": \"/home/root/0.ogg\"}"
	sender(conn, str)
}

func sender(conn net.Conn, str string) {
	conn.Write([]byte(str))
	log.Println("send over")
}
