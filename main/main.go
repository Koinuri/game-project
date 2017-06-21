package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"math"

	"github.com/koinuri/game-project/main/framework"
	"github.com/koinuri/game-project/main/global"
)

const (
	width  = 1366
	height = 768
)

//Initializes the the program.
func Init() {
	runtime.LockOSThread()
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	global.Directory = path.Dir(ex)
	global.Width = width
	global.Height = height
}

func main() {
	//Initialize the global variables and main thread
	Init()

	//Initialize the window and OpenGL program to draw, as well as closing it when the job is done
	window, program := framework.Init(width, height)
	defer framework.Clean()

	//Drawing logic, only draws cute Akane right in the middle of the window for testing purposes
	spr1 := framework.InitSprite("kotonoha-7.png", framework.InitCanvas(1600, 900, 0, 0))
	spr2 := framework.InitSprite("kotonoha-7.png", framework.InitCanvas(1600, 900, 0, 0))

	fmt.Println(spr1)

	spr2.Scale(.5)

	angle := float64(0.0)
	//Main loop to draw the drawing logic created
	for !window.ShouldClose() {
		framework.InitFrame()
		rad := angle * (math.Pi / 180)

		spr2.AngleRotate(angle + 90)
		spr2.Move(450 * math.Cos(float64(rad)) * .5, 450 * math.Sin(float64(rad)) * .5)

		sprites := make([]framework.Artist, 0)

		sprites = append(sprites, &spr1)
		sprites = append(sprites, &spr2)

		framework.Draw(sprites, program)

		framework.SwapWindowAndPollEvents(window)
		angle += 1
	}
}
