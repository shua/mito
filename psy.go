package main

import (
    "github.com/gorilla/websocket"
    "time"
    "encoding/json"
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

type EventID int
const (
    Conception EventID = iota
    Birth
    Intent
    Move
    Grow
    Split
    Die
)

type mind struct {
    ws *websocket.Conn
    send chan []byte
    //bodies map[*body] int
    name string
    background string

    die chan *body
    move chan *body
}

func (m *mind) readPump() {
    defer func() {
        l.leave <- m
        m.ws.Close()
    }()
    m.ws.SetReadLimit(maxMessage)
    m.ws.SetReadDeadline(time.Now().Add(pongWait))
    m.ws.SetPongHandler(func(string) error {
        m.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil
    })
    for {
        _, message, err := m.ws.ReadMessage()
        if err != nil {
            break
        }

        l.broadcast <- message
    }
}

func (m *mind) write(mt int, payload []byte) error {
    m.ws.SetWriteDeadline(time.Now().Add(writeWait))
    return m.ws.WriteMessage(mt, payload)
}

func (m *mind) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        m.ws.Close()
    }()
    for {
        select {
        case message, ok := <-m.send:
            if !ok {
                m.write(websocket.CloseMessage, []byte{})
                return
            }
            if err := m.write(websocket.TextMessage, message); err != nil {
                return
            }
        case <-ticker.C:
            if err := m.write(websocket.PingMessage, []byte{}); err != nil {
                return
            }
        }
    }
}

