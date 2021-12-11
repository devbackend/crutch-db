package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conf := flag.String("c", "localhost:2357", "running config")

	flag.Parse()

	conn, err := net.Dial("tcp", *conf)
	if err != nil {
		log.Fatal("connect err", err)
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("> ")

		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("err on read", err)
			continue
		}

		text = strings.TrimSpace(text)

		if strings.ToLower(text) == "exit" {
			fmt.Println("stopping...")
			break
		}

		_, err = fmt.Fprintf(conn, text+"\n")
		if err != nil {
			log.Println("err on send", err)
			continue
		}

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print(message)
	}
}
