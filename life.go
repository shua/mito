package main

type life struct {
    connections map[*mind] bool
    broadcast chan []byte
    join chan *mind
    leave chan *mind
}

var l = life{
    connections:    make(map[*mind]bool),
    broadcast:      make(chan []byte),
    join:           make(chan *mind),
    leave:          make(chan *mind),
}

func (l *life) run() {
    for {
        select {
        case m := <-l.join:
            l.connections[m] = true
        case m := <-l.leave:
            if _, ok := l.connections[m]; ok {
                delete(l.connections, m)
                close(m.send)
            }
        case s := <-l.broadcast:
            for m := range l.connections {
                select {
                case m.send <- s:
                default:
                    close(m.send)
                    delete(l.connections, m)
                }
            }
        }
    }
}
