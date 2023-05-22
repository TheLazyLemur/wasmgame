package main

import "syscall/js"

type Player struct {
	Height float32
	Width  float32
	X, Y   float32
	Color  string
}

func (p *Player) Draw() {
	Cvs.DrawRect(player.Color, player.X, player.Y, player.Width, player.Height)
}

type Canvas struct {
	cvs js.Value
	ctx js.Value
}

func (c *Canvas) DrawRect(color string, x, y, width, height float32) {
	c.ctx.Set("fillStyle", color)
	c.ctx.Call("fillRect", x, y, width, height)
}

func (c *Canvas) Clear() {
	c.ctx.Set("fillStyle", Screen_Background)
	c.ctx.Call("fillRect", 0, 0, Screen_Width, Screen_Height)
}

type Bullet struct {
	Height float32
	Width  float32
	X, Y   float32
	Color  string
}

func (b *Bullet) Draw() {
	Cvs.DrawRect(b.Color, b.X, b.Y, b.Width, b.Height)
}

func (b *Bullet) Move() {
	b.X++
}
