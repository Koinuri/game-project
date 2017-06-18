package main

import (
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

	//Drawing logic, only draws cute Akane right in the middle of the window for testing purposes
	sprtopleft := framework.InitSprite("kotonoha-7.png", framework.InitCanvas(.5, .5, -.5, .5))
	sprtopright := framework.InitSprite("kotonoha-7.png", framework.InitCanvas(.5, .5, .5, .5))
	sprbotleft := framework.InitSprite("kotonoha-7.png", framework.InitCanvas(.5, .5, -.5, -.5))
	sprbotright := framework.InitSprite("kotonoha-7.png", framework.InitCanvas(.5, .5, .5, -.5))

	//Main loop to draw the drawing logic created
	for !window.ShouldClose() {
		framework.InitFrame(program)

		framework.Draw(&sprtopleft)
		framework.Draw(&sprtopright)
		framework.Draw(&sprbotleft)
		framework.Draw(&sprbotright)

		framework.SwapWindowAndPollEvents(window)
	}
}
