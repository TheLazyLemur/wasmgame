package main

import (
	"fmt"
	"syscall/js"
)

type Player struct {
	Height float32
	Width  float32
	X, Y   float32
	Color  string
	Alive  bool
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

func (c *Canvas) Text(text string, x, y float32) {
	c.ctx.Set("fillStyle", "black")
	c.ctx.Call("fillText", text, x, y, Screen_Width)
}

type Bullet struct {
	Height   float32
	Width    float32
	X, Y     float32
	Color    string
	LifeTime float32
}

func (b *Bullet) Draw() {
	Cvs.DrawRect(b.Color, b.X, b.Y, b.Width, b.Height)
}

func (b *Bullet) Move(dt float32) {
	b.X += 300 * dt
}

type BulletPool struct {
	Bullets []*Bullet
}

func NewBulletPool(initAmount int) *BulletPool {
	p := &BulletPool{
		Bullets: make([]*Bullet, 0),
	}

	for i := 0; i < initAmount; i++ {
		p.Bullets = append(p.Bullets, &Bullet{
			Height:   3,
			Width:    3,
			X:        0,
			Y:        0,
			Color:    "blue",
			LifeTime: 500,
		})
	}

	return p
}

func GetFromPool() *Bullet {
	if len(bulletPool.Bullets) == 0 {
		fmt.Println("Bullet pool is empty")
		return &Bullet{
			Height:   3,
			Width:    3,
			X:        0,
			Y:        0,
			Color:    "blue",
			LifeTime: 500,
		}
	}

	bul := bulletPool.Bullets[0]
	bulletPool.Bullets = bulletPool.Bullets[1:]
	bul.LifeTime = 500
	return bul
}

func AddToPool(b *Bullet) {
	bulletPool.Bullets = append(bulletPool.Bullets, b)
	fmt.Println(len(bulletPool.Bullets))
}

type Obstacle struct {
	Height float32
	Width  float32
	X, Y   float32
	Color  string
	Alive  bool
}
