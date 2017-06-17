package framework

import (
	"image"
	"os"
	"fmt"
	"image/draw"
	"image/png"
	"math"
	"path"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/koinuri/game-project/main/global"
)
type Sprite struct {
	X     float32
	Y     float32
	vao   uint32
	texture uint32
}

type origin uint32

const (
	_ = iota
	TopLeft origin = iota
	TopCenter
	TopRight
	CenterLeft
	Center
	CenterRight
	BottomLeft
	BottomCenter
	BottomRight
)

//Constructor for Sprite struct.  There are several possible arguments for the sprite:
//  * New(string) - Creates a sprite with the image with string showing the relative location of the image.  Origin will be the center.
//  * New(string, origin) - Creates a sprite with the image using the image given by string.  Origin places the origin at the specified place.  Top left corner for TopLeft, for example.
func InitSprite(i ...interface{}) Sprite {
	var dir string
	var x float32
	var y float32
	//The string must be provided if it's more than one
	if len(i) >= 1 {
		test, succ := i[0].(string)
		if !succ {
			panic(fmt.Sprintf("Invalid argument.  Expected string, got %T", i[0]))
		}
		dir = test
	}

	if len(i) == 2 {
		_, succ := i[1].(origin)
		if !succ {
			panic(fmt.Sprintf("Invalid argument.  Expected Origin, got %T", i[1]))
		}
		loc, _ := i[1].(int32)
		x = float32(math.Mod(float64(loc), 3.0) - 2)
		y = float32(loc) / 3.0 - 2
	}

	if len(i) == 0 || len(i) > 3 {
		panic("Invalid number of arguments.  Could not match with any of the possible argument numbers")
	}

	img, err := createImage(dir)
	if err != nil {
		panic(fmt.Sprintf("Could not find file at |%v|\n%v", dir, err))
	}
	vao := createVao(img)
	texture := createTexture(img)
	return Sprite{
		X:     x,
		Y:     y,
		vao: vao,
		texture: texture,
	}
}

func createImage(dir string) (*image.RGBA, error) {
	//read the image file, return empty rgba if error
	imgFile, err := os.Open(path.Join(global.Directory, dir))
	if err != nil {
		a := image.Rectangle{image.Point {0,0}, image.Point {0,0}}
		return image.NewRGBA(a), err
	}
	defer imgFile.Close()

	//decode the image file to use, return empty rgba is error
	img, err := png.Decode(imgFile)
	if err != nil {
		a := image.Rectangle{image.Point {0,0}, image.Point {0,0}}
		return image.NewRGBA(a), err
	}

	//create a rgba canvas to draw the stuff on
	rgba := image.NewRGBA(img.Bounds())

	//draw the image file based on image file
	draw.Draw(rgba, rgba.Bounds(), img, image.ZP, draw.Over)

	return rgba, nil
}

func createVao(img *image.RGBA) uint32 {
	//calculate the image's x and y depending on image aspect ratio
	rec := img.Rect
	var x float32
	var y float32
	if ratio := float32(rec.Max.X) / float32(rec.Max.Y); ratio > 1.0 {
		ratio = 1 / ratio
		x = 1.0
		y = ratio / 2
	} else {
		x = ratio / 2
		y = 1.0
	}

	//create vertices based on the calculated x's and y's and the coordinate of image each vertices should be associated to
	var vec []float32 = []float32{ 
		x * -1, y, 0, 0.0, 0.0, //top left
		x, y, 0, 1.0, 0.0,               //top right
		x * -1, y * -1, 0, 0.0, 1.0, //bottom left
		x, y * -1, 0, 1.0, 1.0, //bottom right
	}

	//the indices to create rectangles using the vectors
	var ind []uint32 = []uint32{
		0, 1, 2,
		1, 2, 3,
	}

	//initialize vao then binding it to vertex array
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	//create vbo that will be used for the currently binded vao
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vec) * 4, gl.Ptr(vec), gl.STATIC_DRAW)

	//create ebo that will be used for the currently binded vao
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(ind) * 4, gl.Ptr(ind), gl.STATIC_DRAW)

	//set the attributes
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5 * 4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5 * 4, gl.PtrOffset(3 * 4))
	gl.EnableVertexAttribArray(1)

	return vao
}

func createTexture(img *image.RGBA) uint32 {
	//initiate texture
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	//set how the texture behaves when the shape created by vao is larger than texture.  It shouldn't happen so just setting it to repeat
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	//set how the texture should guess the pixels when shrinking/enlargening an image.  linear makes an average, blurry, nearest takes the nearest color, blocky
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	//bind the image to this texture
	gl.TexImage2D(
		gl.TEXTURE_2D, 
		0, 
		gl.RGBA, 
		int32(img.Rect.Size().X), 
		int32(img.Rect.Size().Y), 
		0, 
		gl.RGBA, 
		gl.UNSIGNED_BYTE, 
		gl.Ptr(img.Pix))

	//create mipmap, which makes smaller image crispy
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return texture
}

//Returns the vao and texture held by the sprite to draw
func (s *Sprite) GetDrawInfo() (uint32, uint32) {
	return s.vao, s.texture
}