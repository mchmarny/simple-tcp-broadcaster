package client

import (
	"bufio"
	"log"
	"os"

	"github.com/mchmarny/simple-tcp-broadcaster/commons"
)

var (
	agent *commons.Agent
)

// StopClient stops client
func StopClient() error {
	if agent != nil {
		return agent.Stop()
	}
	return nil
}

// StartClient starts a client and connects to server
func StartClient(serverAddress string) error {

	var err error
	agent, err = commons.NewClientAgent(serverAddress)
	if err != nil {
		return err
	}

	log.Printf("Client ID: %s", agent.GetLocalClientID())
	go agent.Read()

	for {
		message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		msg := commons.NewMessage(agent.GetLocalClientID())
		msg.Data = []byte(message)
		if err := agent.Write(msg); err != nil {
			log.Fatalf("Error on write: %v", err)
		}

	}

}
