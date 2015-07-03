package main

import (
    "math"
)


/* linear algebra functions */
type vec2 struct {
    x int
    y int
}

func v2add(a vec2, b vec2) vec2 {
    return vec2{ x: a.x+b.x, y: a.y+b.y }
}

func v2scale(a vec2, s int) vec2 {
    return vec2{ x: a.x*s, y: a.y*s }
}

func v2neg(a vec2) vec2 {
    return v2scale(a, -1)
}

func v2sqHy(a vec2) int {
    return a.x*a.x + a.y*a.y
}

func v2norm(a vec2) vec2 {
    hy := math.Sqrt(float64(v2sqHy(a)))
    return vec2{ x: int(float64(a.x)/hy), y: int(float64(a.y)/hy) }
}

/* physics bodies */
type body struct {
    pos vec2
    size int
    density int
}

func (b *body) checkColl(c *body) bool {
    var minrad int
    if(b.size < c.size) {
        minrad = b.size
    } else {
        minrad = c.size
    }
    sqminrad := minrad * minrad
    sqhy := v2sqHy(v2add(v2neg(b.pos), c.pos))

    return sqhy <= sqminrad
}

/* world should eventually be split into quadtrees */
type world struct {
    bodies map[*body] *mind
    bounds vec2
}

