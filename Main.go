package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

/*
 * CONFIG VARIABLES
 */
type config struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

func main() {
	// ----------------
	// Importing Config
	// ----------------
	fmt.Print("Importing config.yml")
	quit := loading()
	var config config
	config.load("config.yml")
	quit <- "Done!"
	// ----------------

	// ----------------
	// START TCP SERVER
	// ----------------
	fmt.Print("Starting Server on " + address + ":" + strconv.Itoa(port) + "")
	quit = loading()
	serverStatus := startServer()
	if i := <-serverStatus; i == 1 {
		quit <- "Done!"
	} else if i == 2 {
		quit <- "An error during server startup has occured!"
	}
	// ----------------

	// --------------------
	// Starting Interpreter
	// --------------------
	receive := make(chan string)
	go startInterpreter(receive)
	receive <- "test"
	// --------------------

	<-make(chan struct{})
	return
}

func loading() chan string {
	quit := make(chan string)
	go func() {
		for {
			select {
			case msg := <-quit:
				fmt.Print(" " + msg + "\n")
				return
			default:
				fmt.Print(".")
				time.Sleep(1000 * time.Millisecond)
			}
		}
	}()
	return quit
}

func startInterpreter(textChannel chan string) {
	for {
		select {
		case text := <-textChannel:
			fmt.Print(text + "\n")
		}
	}
}

func startServer() chan int {
	status := make(chan int)
	/*
	 * 0 -> Offline
	 * 1 -> Online
	 * 2 -> Error
	 */
	go func() {
		l, err := net.Listen("tcp", address+":"+strconv.Itoa(port))
		if err != nil {
			fmt.Println(err)
			status <- 2
			return
		}
		status <- 1
		for {
			c, err := l.Accept()
			if err != nil {
				fmt.Println(err)
				status <- 2
				return
			}

			netData := bufio.NewReader(c)
			netDataLine, err := netData.ReadString('\n')
			if err != nil {
				fmt.Println(err)
			} else {
				if strings.TrimSpace(string(netDataLine)) == "STOP" {
					fmt.Println("TCP Stream Closed!")
				} else {
					fmt.Print("TCP Packet Received!: ", string(netDataLine))
				}
			}
			//t := time.Now()
			//myTime := t.Format(time.RFC3339) + "\n"
			//c.Write([]byte(myTime))
		}
	}()
	return status
}

func (c *config) load(path string) *config {

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	return c
}
