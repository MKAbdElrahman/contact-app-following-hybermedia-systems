package main

import (
	"app/protogen/tutorialpb"
	"log"
	"os"

	"google.golang.org/protobuf/proto"
)

func main() {
	msg := &tutorialpb.Hello{
		Name: "Hello MK!",
	}

	out, err := proto.Marshal(msg)

	if err != nil {
		log.Fatalln("Failed to encode  message:", err)
	}

	if err := os.WriteFile("test", out, 0644); err != nil {
		log.Fatalln("Failed to write message:", err)
	}
}
