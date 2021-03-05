package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func readFromRemoteStdout(c *websocket.Conn, done *chan bool, pongWait time.Duration) {
	defer func() { *done <- true }()
	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mt, r, err := c.NextReader()
		if websocket.IsCloseError(err,
			websocket.CloseNormalClosure, // Normal.
		) {
			return
		}
		if err != nil {
			fmt.Printf("nextreader: %v\n", err)
			return
		}
		if mt != websocket.BinaryMessage {
			fmt.Println("binary message")
			return
		}
		if _, err := io.Copy(os.Stdout, r); err != nil {
			fmt.Printf("Reading from websocket: %v\n", err)
			return
		}
	}
}

func writeToRemoteStdin(c *websocket.Conn, done *chan bool, writeMutex *sync.Mutex, writeWait time.Duration) {
	defer func() { *done <- true }()
	for {
		var input []byte = make([]byte, 1)
		os.Stdin.Read(input)
		writeMutex.Lock()
		c.SetWriteDeadline(time.Now().Add(writeWait))
		err := c.WriteMessage(websocket.BinaryMessage, append([]byte{0}, input...))
		writeMutex.Unlock()
		if err != nil {
			fmt.Println("write:", err)
			return
		}
	}
}

func sendPingMessages(c *websocket.Conn, done *chan bool, writeWait time.Duration, pingPeriod time.Duration) {
	defer func() { *done <- true }()
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := c.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				fmt.Println("ping:", err)
				return
			}
		}
	}
}
