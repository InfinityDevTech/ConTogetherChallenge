// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

func handleErrs(vari interface{}, err error) interface{} {
	if err != nil {
		log.Println(err)
		return err
	}
	return vari 
}

func main() {
	exec.Command("cmd", "/C", "title", "The one true notepad!").Run()
	if _, err := os.Stat("readme.txt"); errors.Is(err, os.ErrNotExist) {
		os.Create("readme.txt")
	  }
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "b8ec-75-69-111-14.ngrok.io", Path: "/"}
	color.Red("Connecting to %s", u.Host)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			//write to file readme.txt
			f, err := os.OpenFile("readme.txt", os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}
			f.Write(message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	var last []byte = nil
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			type data struct {
				Fname string
				Fdata []byte
			}
			fData := handleErrs(ioutil.ReadFile("readme.txt"))
			//If its different just print in the console that it is
			if last != nil && string(last) != string(fData.([]byte)) {
				d := handleErrs(json.Marshal(data{"hello", fData.([]byte)}))
				err := c.WriteMessage(websocket.TextMessage, d.([]byte))
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
			last = fData.([]byte)
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
