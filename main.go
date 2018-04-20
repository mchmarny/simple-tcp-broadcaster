package main

import (
	"errors"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"os"

	"github.com/mchmarny/simple-tcp-broadcaster/client"
	"github.com/mchmarny/simple-tcp-broadcaster/server"
	"github.com/urfave/cli"
)

const (
	appName    = "simple"
	appVersion = "0.2.0"

	undefinedCLIMode cliMode = 0
	clientCLIMode    cliMode = 1
	serverCLIMode    cliMode = 2
)

type cliMode int

var (
	serverCommand = cli.Command{
		Name:   "server",
		Usage:  "Simple Server",
		Action: printCommandHelp,
		Subcommands: []cli.Command{
			{
				Name:   "start",
				Usage:  "Starts server",
				Action: startServer,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "port",
						Usage: "Server port on which the server should listen",
					},
				},
			},
		},
	}

	clientCommand = cli.Command{
		Name:   "client",
		Usage:  "Simple Server Client",
		Action: printCommandHelp,
		Subcommands: []cli.Command{
			{
				Name:   "connect",
				Usage:  "Starts client and connects to target server",
				Action: startClient,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "address",
						Usage: "Server address",
					},
					cli.IntFlag{
						Name:  "port",
						Usage: "Server port on which the server listens",
					},
				},
			},
		},
	}
)

func main() {

	cmd := cli.NewApp()
	cmd.Name = appName
	cmd.Usage = fmt.Sprintf("Simple Server CLI (%s v%s)", appName, appVersion)
	cmd.Version = appVersion
	cmd.Commands = []cli.Command{
		serverCommand,
		clientCommand,
	}

	err := cmd.Run(os.Args)
	if err != nil {
		fmt.Printf("\nError: %v\n\n", err)
	} else {
		fmt.Println()
	}

}

func printCommandHelp(c *cli.Context) error {
	return cli.ShowSubcommandHelp(c)
}

func startServer(c *cli.Context) error {

	port := c.Int("port")
	if port < 1024 {
		return errors.New("Server port must be above 1024")
	}

	go handleConsoleSignal(serverCLIMode)

	return server.StartServer(port)

}

func startClient(c *cli.Context) error {

	port := c.Int("port")
	if port < 1024 {
		return errors.New("Server port must be above 1024")
	}

	// if "" then localhost
	address := c.String("address")

	go handleConsoleSignal(clientCLIMode)

	serverAddress := fmt.Sprintf("%s:%d", address, port)

	return client.StartClient(serverAddress)
}

// handleConsoleSignal Waits for SIGINT and SIGTERM (HIT CTRL-C)
func handleConsoleSignal(mode cliMode) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	if mode == clientCLIMode {
		log.Println("Shutting down client...")
		// TODO: what on client?
	} else if mode == serverCLIMode {
		log.Println("Shutting down server...")
		server.StopServer()
	} else {
		log.Println("Error, CLI mode not set")
	}
	os.Exit(0)
}
