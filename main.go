package main

import (
	"fmt"
	"syscall/js"
)

const (
	Screen_Width      = 400
	Screen_Height     = 400
	Screen_Background = "white"
)

var (
	Cvs *Canvas = &Canvas{}

	player *Player = &Player{
		Height: 10,
		Width:  10,
		X:      0,
		Y:      0,
		Color:  "red",
		Alive:  true,
	}

	bullets = make(map[*Bullet]struct{})

	bulletPool = NewBulletPool(20)
)

func initGame(this js.Value, args []js.Value) interface{} {
	canvas := args[0]
	ctx := canvas.Call("getContext", "2d")

	Cvs.cvs = canvas
	Cvs.ctx = ctx

	return nil
}

func gameInput(this js.Value, args []js.Value) interface{} {
	switch args[0].String() {
	case "ArrowLeft":
		player.X--
	case "ArrowRight":
		player.X++
	case "ArrowDown":
		player.Y++
	case "ArrowUp":
		player.Y--
	case " ":
		b := GetFromPool()
		b.X = player.X
		b.Y = player.Y

		bullets[b] = struct{}{}
	default:
		fmt.Println("Unknown input:", args[0])
	}
	return nil
}

type Obstacle struct {
	Height float32
	Width  float32
	X, Y   float32
	Color  string
	Alive  bool
}

var obstancle = &Obstacle{
	Height: 50,
	Width:  50,
	X:      Screen_Width / 2,
	Y:      Screen_Height / 2,
	Color:  "magenta",
	Alive:  true,
}

func gameUpdate(this js.Value, args []js.Value) interface{} {
	Cvs.Clear()
	if player.Alive {
		player.Draw()
	}
	for b := range bullets {
		b.Move()
		b.Draw()
		b.LifeTime--
		if b.LifeTime <= 0 {
			fmt.Println("bullet die")
			AddToPool(b)
			delete(bullets, b)
		}
	}

	if obstancle.Alive {
		Cvs.DrawRect(obstancle.Color, obstancle.X, obstancle.Y, obstancle.Width, obstancle.Height)
	}

	if player.Alive {
		if player.X+player.Width > obstancle.X && player.X < obstancle.X+obstancle.Width {
			if player.Y+player.Height > obstancle.Y && player.Y < obstancle.Y+obstancle.Height {
				fmt.Println("hit")
				player.Alive = false
			}
		}
	}

	for b := range bullets {
		if obstancle.X+obstancle.Width > b.X && obstancle.X < b.X+b.Width {
			if obstancle.Y+obstancle.Height > b.Y && obstancle.Y < b.Y+b.Height {
				fmt.Println("hit")
				obstancle.Alive = false
				b.LifeTime = 0
			}
		}
	}

	return nil
}

func main() {
	js.Global().Set("gameUpdate", js.FuncOf(gameUpdate))
	js.Global().Set("gameInput", js.FuncOf(gameInput))
	js.Global().Set("initGame", js.FuncOf(initGame))

	select {}
}
