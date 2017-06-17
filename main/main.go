package main

import (
	"runtime"
	"os"
	"path"

	"github.com/koinuri/game-project/main/framework"
	"github.com/koinuri/game-project/main/global"
)

const (
	width              = 1366
	height             = 768
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
	spr := framework.InitSprite("kotonoha-7.png", framework.InitCanvas(1.0, 1.0))

	//Main loop to draw the drawing logic created
	for !window.ShouldClose() {
		framework.Draw(&spr, window, program)
	}
}
