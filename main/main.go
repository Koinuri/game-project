package main

import (
	//	"fmt"
	"math"
	"os"
	"path"
	"runtime"

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

	//Drawing logic
	obj := framework.InitObject()
	spr1 := obj.CreateSprite("first akane chan", "kotonoha-7.png", framework.BottomCenter)
	spr2 := obj.CreateSprite("second akane chan", "kotonoha-7.png", framework.BottomCenter)
	spr3 := obj.CreateSprite("third akane chan", "kotonoha-7.png", framework.BottomCenter)

	spr1.AngleRotate(90)
	spr2.AngleRotate(210)
	spr3.AngleRotate(330)

	angle := float64(0.0)
	//Main loop to draw the drawing logic created
	for !window.ShouldClose() {
		framework.InitFrame()

		rad := angle * (math.Pi / 180)

		obj.AngleRotate(angle * 3)

		obj.Move(450*math.Cos(rad)*.5, 450*math.Sin(rad)*.5)

		obj.Scale(.3)

		framework.Draw(obj.GetArtists(), program)

		framework.SwapWindowAndPollEvents(window)
		angle += 1
	}
}
