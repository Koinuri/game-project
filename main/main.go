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
	Init()

	window, program := framework.Init(width, height)
	defer framework.Clean()

	spr := framework.InitSprite("kotonoha-7.png", framework.InitCanvas(1.0, 1.0))

	for !window.ShouldClose() {
		framework.Draw(&spr, window, program)
	}
}
