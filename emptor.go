package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// Emptor (aka Client)
// Conn : TCP connection struct
// Topics : Map of topics the client is currently subscribed to
// Commands : Channel to recieve the commands and send it to broker
type Emptor struct {
	conn     net.Conn
	topics   map[string]*Topic
	commands chan<- Command
}

// keep listening for messages on the TCP connection and send it to the broker
func (e *Emptor) readMessage() {
	for {
		msg, err := bufio.NewReader(e.conn).ReadString('\n')

		// Client disconnected or some error in TCP connection.
		// Either way delete the client info from its currently subbed topics
		if err != nil {
			e.commands <- Command{
				id:     QUIT,
				emptor: e,
				args:   []string{},
			}
			log.Println("error reading message from client :", err)
			return
		}
		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])
		if len(args) == 0 || args[0] == "" || args[0] == "\n" || args[0][0] == '\t' {
			continue
		}
		switch cmd {
		case "pub":
			e.commands <- Command{
				id:     PUB,
				emptor: e,
				args:   args,
			}
		case "sub":
			e.commands <- Command{
				id:     SUB,
				emptor: e,
				args:   args,
			}
		case "unsub":
			e.commands <- Command{
				id:     UNSUB,
				emptor: e,
				args:   args,
			}
		case "ls":
			e.commands <- Command{
				id:     LIST,
				emptor: e,
				args:   args,
			}
		case "quit":
			e.commands <- Command{
				id:     QUIT,
				emptor: e,
				args:   args,
			}
		default:
			e.Err(fmt.Errorf("invalid operation"))
		}
	}
}

// send an error message to user. 0 indicates error.
func (e *Emptor) Err(err error) {
	e.conn.Write([]byte("0 " + err.Error() + "\n"))
}

// send an success message to user. 1 indicates success.
func (e *Emptor) Msg(m string) {
	e.conn.Write([]byte("1 " + m + "\n"))
}

// send an warning message to user. -1 indicates warning.
func (e *Emptor) Warn(m string) {
	e.conn.Write([]byte("-1 " + m + "\n"))
}
