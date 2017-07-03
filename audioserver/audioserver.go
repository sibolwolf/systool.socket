package main

import (
	"encoding/json"
	"log"

	"smartconn.cc/tosone/audioctrller"

	"github.com/systool.socket/socket/server"
	//"smartconn.cc/tosone/audioctrller"
)

type audioType map[string]string

func main() {
	IP := "localhost:1024"
	server.ServerSetup(IP, audiohandle)
}

func audiohandle(buffer []byte) {
	var audioOnce audioType
	log.Println(len(buffer))
	if err := json.Unmarshal(buffer, &audioOnce); err != nil {
		log.Println("Unmarshal err:", err)
		return
	}
	for k, v := range audioOnce {
		switch k {
		case "Play":
			audioctrller.Play(v)
			/*
				cmd := "play " + v
				log.Println(cmd)
				rs := exec.Command("/bin/bash", "-c", cmd)
				if _, err := rs.Output(); err != nil {
					log.Println(err)
				}
			*/
		}
	}

}
