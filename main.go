package main

import (
	"fmt"
	"sync"
	"syscall/js"
	"time"
)

const (
	Screen_Width      = 400
	Screen_Height     = 400
	Screen_Background = "white"
)

var (
	fps               float64 = 120
	frameDuration             = time.Second / time.Duration(fps)
	startTime                 = time.Now()
	previousFrameTime         = startTime
	dTime             float64 = 0
	deltaTime         float64 = 0
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

func gameInput() interface{} {
	for k := range keys {
		switch k {
		case "ArrowLeft":
			player.X -= 300 * float32(deltaTime)
		case "ArrowRight":
			player.X += 300 * float32(deltaTime)
		case "ArrowDown":
			player.Y += 300 * float32(deltaTime)
		case "ArrowUp":
			player.Y -= 300 * float32(deltaTime)
		case " ":
			if !player.Alive {
				continue
			}
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
	currntTime := time.Now()
	deltaTime = currntTime.Sub(previousFrameTime).Seconds()
	previousFrameTime = currntTime

	dTime -= deltaTime
	if dTime <= 0 {
		dTime = 1
		fps = 1 / deltaTime
	}

	gameInput()
	Cvs.Clear()

	Cvs.Text("FPS: "+fmt.Sprintf("%v", fps), 10, 10)

	if player.Alive {
		player.Draw()
	}
	for b := range bullets {
		b.Move(float32(deltaTime))
		b.Draw()
		b.LifeTime--
		if b.LifeTime <= 0 {
			AddToPool(b)
			delete(bullets, b)
		}
	}

	if obstancle.Alive {
		Cvs.DrawRect(obstancle.Color, obstancle.X, obstancle.Y, obstancle.Width, obstancle.Height)
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if player.Alive {
			if player.X+player.Width > obstancle.X && player.X < obstancle.X+obstancle.Width {
				if player.Y+player.Height > obstancle.Y && player.Y < obstancle.Y+obstancle.Height {
					player.Alive = false
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
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
	}()
	wg.Wait()

	elapsedTime := time.Since(startTime)
	remainingTime := frameDuration - elapsedTime
	if remainingTime > 0 {
		time.Sleep(time.Duration(remainingTime))
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

func initGame(this js.Value, args []js.Value) interface{} {
	document := args[0]
	canvas := document.Call("getElementById", "myCanvas")
	ctx := canvas.Call("getContext", "2d")

	Cvs.cvs = canvas
	Cvs.ctx = ctx

	return nil
}

func main() {
	js.Global().Set("initGame", js.FuncOf(initGame))

	js.Global().Call("requestAnimationFrame", js.FuncOf(gameUpdate))
	js.Global().Call("addEventListener", "keydown", js.FuncOf(HandleKeys))
	js.Global().Call("addEventListener", "keyup", js.FuncOf(HandleKeysUp))

	select {}
}
