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

func gameUpdate(this js.Value, args []js.Value) interface{} {
	Cvs.Clear()
	player.Draw()
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

	return nil
}

func main() {
	js.Global().Set("gameUpdate", js.FuncOf(gameUpdate))
	js.Global().Set("gameInput", js.FuncOf(gameInput))
	js.Global().Set("initGame", js.FuncOf(initGame))

	select {}
}
