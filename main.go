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

	obstancle = &Obstacle{
		Height: 50,
		Width:  50,
		X:      Screen_Width / 2,
		Y:      Screen_Height / 2,
		Color:  "magenta",
		Alive:  true,
	}

	keys = make(map[string]struct{})
)

func initGame(this js.Value, args []js.Value) interface{} {
	document := args[0]
	canvas := document.Call("getElementById", "myCanvas")
	ctx := canvas.Call("getContext", "2d")

	Cvs.cvs = canvas
	Cvs.ctx = ctx

	return nil
}

func gameInput() interface{} {
	for k := range keys {
		switch k {
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
			fmt.Println("Unknown input:", k)
		}
	}
	return nil
}

func gameUpdate(this js.Value, args []js.Value) interface{} {
	gameInput()
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
				player.Alive = false
			}
		}
	}

	if obstancle.Alive {
		for b := range bullets {
			if obstancle.X+obstancle.Width > b.X && obstancle.X < b.X+b.Width {
				if obstancle.Y+obstancle.Height > b.Y && obstancle.Y < b.Y+b.Height {
					obstancle.Alive = false
					b.LifeTime = 0
				}
			}
		}
	}

	js.Global().Call("requestAnimationFrame", js.FuncOf(gameUpdate))
	return nil
}

func HandleKeys(this js.Value, args []js.Value) interface{} {
	event := args[0]
	if event.Get("repeat").Bool() {
		return nil
	}

	keys[event.Get("key").String()] = struct{}{}
	return nil
}

func HandleKeysUp(this js.Value, args []js.Value) interface{} {
	delete(keys, args[0].Get("key").String())
	return nil
}

func main() {
	js.Global().Set("initGame", js.FuncOf(initGame))

	js.Global().Call("requestAnimationFrame", js.FuncOf(gameUpdate))
	js.Global().Call("addEventListener", "keydown", js.FuncOf(HandleKeys))
	js.Global().Call("addEventListener", "keyup", js.FuncOf(HandleKeysUp))

	select {}
}
