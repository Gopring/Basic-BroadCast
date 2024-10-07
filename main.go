package main

import (
	"WebRTC_POC/server"
	"log"
)

func main() {
	s := server.New()
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
