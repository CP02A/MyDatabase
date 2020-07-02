package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

/*
 * CONFIG STRUCTS
 */
type config struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

// --------------

/*
 * TABLE STRUCTS
 */
var tables map[string]table

type table struct {
	Columns []tableColumn
}

type tableColumn struct {
	Name string
	Type string
}

// --------------

func main() {
	// -----------------
	// Startup Messaging
	// -----------------
	fmt.Print("---------------------------\n")
	fmt.Print("\tMy Database\n")
	fmt.Print("---------------------------\n")
	// -----------------

	// ----------------
	// Importing Config
	// ----------------
	var quit chan string
	var config config
	if _, err := os.Stat("config.yml"); err == nil {
		fmt.Print("Importing config.yml")
		quit = loading()
	} else if os.IsNotExist(err) {
		fmt.Print("config.yml does not exist! Creating the file")
		quit = loading()
		f, err := os.Create("config.yml")
		if err != nil {
			quit <- "Error!"
			fmt.Println(err)
			return
		}
		f.WriteString("address: 0.0.0.0\nport: 6666")
		f.Close()
	}
	config.load("config.yml")
	quit <- "Done!"
	// ----------------

	// --------------------
	// Starting Interpreter
	// --------------------
	interpret := make(chan string)
	go startInterpreter(interpret)
	// --------------------

	// ----------------
	// Importing Tables
	// ----------------

	if _, err := os.Stat("tables.json"); err == nil {
		fmt.Print("Importing tables.json")
		quit = loading()
	} else if os.IsNotExist(err) {
		fmt.Print("tables.json does not exist! Creating the file")
		quit = loading()
		f, err := os.Create("tables.json")
		if err != nil {
			quit <- "Error!"
			fmt.Println(err)
			return
		}
		f.WriteString("{}")
		f.Close()
	}
	jsonFile, err := ioutil.ReadFile("tables.json")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal([]byte(jsonFile), &tables)
	if err != nil {
		quit <- "Error!"
		fmt.Println(err)
		panic("failed to import table.json")
	}
	quit <- "Done!"
	// ----------------

	// ----------------
	// START TCP SERVER
	// ----------------
	fmt.Print("Starting Server on " + config.Address + ":" + strconv.Itoa(config.Port) + "")
	quit = loading()
	serverStatus := startServer(config.Address, config.Port, interpret)
	if i := <-serverStatus; i == 1 {
		quit <- "Done!"
	} else if i == 2 {
		quit <- "Error!"
	}
	// ----------------

	// keeps program running
	<-make(chan struct{})
	return
}

func loading() chan string {
	quit := make(chan string)
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			select {
			case msg := <-quit:
				fmt.Print(" " + msg + "\n")
				return
			default:
				fmt.Print(".")
			}
		}
	}()
	return quit
}

func startInterpreter(textChannel chan string) {
	for {
		select {
		case text := <-textChannel:
			fmt.Print(text)
		}
	}
}

func startServer(address string, port int, listener chan string) chan int {
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
					//fmt.Print("TCP Packet Received: ", string(netDataLine))
					listener <- string(netDataLine)
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
