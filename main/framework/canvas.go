package framework

import (
	"github.com/koinuri/game-project/main/global"
)

type Canvas struct {
	Width  float32
	Height float32
	X      float32
	Y      float32
}

//Initiates a canvas to hold the text, sprite, or shape objects.
func InitCanvas(i ...float32) Canvas {
	//Set up the default values
	var width float32 = 1.0
	var height float32 = 1.0
	var x float32 = 0.0
	var y float32 = 0.0
	countFloat := 0

	//Assign values that has been given
	for _, item := range i {
		var f float32 = float32(item)

		if countFloat == 0 {
			width = f
			countFloat++
		} else if countFloat == 1 {
			height = f
			countFloat++
		} else if countFloat == 2 {
			x = f
			countFloat++
		} else if countFloat == 3 {
			y = f
			countFloat++
		} else {
			panic("Invalid argument.  Found more than 4 floats in the argument")
		}
	}

	if countFloat == 1 || countFloat == 3 {
		panic("Invalid argument.  Both the pairs (width, height) and (x, y) must be specified.")
	}

	return Canvas{
		width,
		height,
		x,
		y,
	}
}

func (c *Canvas) GetAspectRatio() float32 {
	return (float32(global.Width) * c.Width) / (float32(global.Height) * c.Height)
}
