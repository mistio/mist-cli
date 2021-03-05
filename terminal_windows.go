// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	terminal "golang.org/x/term"
)

type terminalSize struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

func updateTerminalSize(c *websocket.Conn, writeMutex *sync.Mutex, writeWait time.Duration) error {
	width, height, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("Could not get terminal size %s\n", err)
	}
	resizeMessage := terminalSize{height, width}
	resizeMessageBinary, err := json.Marshal(&resizeMessage)
	if err != nil {
		return fmt.Errorf("Could not marshal resizeMessage %s\n", err)
	}
	writeMutex.Lock()
	c.SetWriteDeadline(time.Now().Add(writeWait))
	err = c.WriteMessage(websocket.BinaryMessage, append([]byte{1}, resizeMessageBinary...))
	writeMutex.Unlock()
	if err != nil {
		return fmt.Errorf("write: %s", err)
	}
	return nil
}

func handleTerminalResize(c *websocket.Conn, done *chan bool, writeMutex *sync.Mutex, writeWait time.Duration) {
	defer func() { *done <- true }()
	ticker := time.NewTicker(1000 * time.Millisecond)
	for {
		<-ticker.C
		err := updateTerminalSize(c, writeMutex, writeWait)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
