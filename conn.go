package main

import (
    "github.com/gorilla/websocket"
    "time"
)

const (
    writeWait = 10 * time.Second
    pongWait = 60 * time.Second
    pingPeriod = (pongWait * 9) / 10
    maxMessage = 512
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type connection struct {
    ws *websocket.Conn
    send chan []byte
}

func (c *connection) readPump() {
    defer func() {
        l.leave <- c
        c.ws.Close()
    }()
    c.ws.SetReadLimit(maxMessage)
    c.ws.SetReadDeadline(time.Now().Add(pongWait))
    c.ws.SetPongHandler(func(string) error {
        c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil
    })
    for {
        _, message, err := c.ws.ReadMessage()
        if err != nil {
            break
        }
        l.broadcast <- message
    }
}

func (c *connection) write(mt int, payload []byte) error {
    c.ws.SetWriteDeadline(time.Now().Add(writeWait))
    return c.ws.WriteMessage(mt, payload)
}

func (c *connection) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.ws.Close()
    }()
    for {
        select {
        case message, ok := <-c.send:
            if !ok {
                c.write(websocket.CloseMessage, []byte{})
                return
            }
            if err := c.write(websocket.TextMessage, message); err != nil {
                return
            }
        case <-ticker.C:
            if err := c.write(websocket.PingMessage, []byte{}); err != nil {
                return
            }
        }
    }
}

