package main

type life struct {
    connections map[*connection] bool
    broadcast chan []byte
    join chan *connection
    leave chan *connection
}

var l = life{
    connections:    make(map[*connection]bool),
    broadcast:      make(chan []byte),
    join:           make(chan *connection),
    leave:          make(chan *connection),
}

func (l *life) run() {
    for {
        select {
        case c := <-l.join:
            l.connections[c] = true
        case c := <-l.leave:
            if _, ok := l.connections[c]; ok {
                delete(l.connections, c)
                close(c.send)
            }
        case m := <-l.broadcast:
            for c := range l.connections {
                select {
                case c.send <- m:
                default:
                    close(c.send)
                    delete(l.connections, c)
                }
            }
        }
    }
}
