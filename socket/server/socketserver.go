package server

import (
	"log"
	"net"
)

// ServerSetup function
func ServerSetup(IP string, f func([]byte)) {
	// Setup socket, and create listen port
	// IP = "localhost:1024"
	netListen, err := net.Listen("tcp", IP)
	CheckErr(err)
	defer netListen.Close()

	Log("Waiting for clients!")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		Log(conn.RemoteAddr().String(), "tcp connect success")
		go handleConnection(conn, f)
	}

}

// CheckErr function
func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}

}

// Log function
func Log(v ...interface{}) {
	log.Println(v...)
}

func handleConnection(conn net.Conn, f func([]byte)) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		Log("n is ", n)

		if err != nil {
			Log(conn.RemoteAddr().String(), "connection err:", err)
			return
		}
		Log(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))
		f(buffer[:n])
		return
	}
}
