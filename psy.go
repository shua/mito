package main

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
    conn *connection
    bodies map[*body] int
    name string
    back string

    die chan *body
    move chan *body
}



