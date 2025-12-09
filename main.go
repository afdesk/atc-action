package main

import (
	"log"

	"github.com/smartforce-io/atc/githubservice/push"
)

func main() {
	log.Println("Automated Tag Creator")
	err := push.AddTag()
	if err != nil {
		log.Fatalf("error creating tag %v", err)
	}
}
