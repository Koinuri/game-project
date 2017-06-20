package framework

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/koinuri/game-project/main/global"
	"github.com/go-gl/mathgl/mgl32"
)

type Sprite struct {
	X       float32
	Y       float32
	vao     uint32
	texture uint32
	translation mgl32.Mat4
	scale mgl32.Mat4
	rotation mgl32.Mat4
}

type origin uint32

const (
	_              = iota
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
//  * New(string)
// 	* New(string, Canvas)
//  * New(string, Origin)
//	* New(string, Canvas, Origin)
//Where:
//	* String is the location of the image, relative to the executable file.
//	* Canvas is the container in which the image will be stored in.  It will be defaulted to the container holding the entire window if it has not been specified.
//	* Origin is where the coordinate system for this image is based on.  For example, TopLeft places the origin on the top left corner of the image, so if the image is moved to (0, 0), the top left of the image will be placed at (0, 0).  It will be defaulted to the center of the image if it has not been specified.
func InitSprite(i ...interface{}) Sprite {
	var dir string
	var x float32
	var y float32
	var canvas Canvas

	//The string must be provided if it's more than one
	if len(i) >= 1 {
		test, succ := i[0].(string)
		if !succ {
			panic(fmt.Sprintf("Invalid argument.  Expected string, got %T", i[0]))
		}
		dir = test
	}

	//If it's two arguments, it can either be origin or canvas
	if len(i) == 2 {
		_, succ1 := i[1].(origin)
		_, succ2 := i[1].(Canvas)
		if !succ1 && !succ2 {
			panic(fmt.Sprintf("Invalid argument.  Expected Origin or Canvas, got %T", i[1]))
		}
		if succ1 {
			loc, _ := i[1].(int32)
			x = float32(math.Mod(float64(loc), 3.0) - 2)
			y = float32(loc)/3.0 - 2
		} else {
			canvas = i[1].(Canvas)
		}
	}

	//If it's three arguments, all string, origin, and canvas are present
	if len(i) == 3 {
		_, succ2 := i[2].(origin)
		_, succ1 := i[1].(Canvas)

		if succ2 {
			loc, _ := i[2].(int32)
			x = float32(math.Mod(float64(loc), 3.0) - 2)
			y = float32(loc)/3.0 - 2
		} else {
			panic(fmt.Sprintf("Invalid argument.  Expected Origin, got %T", i[2]))
		}

		if succ1 {
			canvas = i[1].(Canvas)
		} else {
			panic(fmt.Sprintf("Invalid argument.  Expected Canvas, got %T", i[1]))
		}
	}

	//If arguments are not within intervals [1, 3], it's invalid numbers of arguments
	if len(i) == 0 || len(i) > 3 {
		panic("Invalid number of arguments.  Could not match with any of the possible argument numbers")
	}

	//create an image, then create vao and texture based on that image
	img, err := createImage(dir)
	if err != nil {
		panic(fmt.Sprintf("Could not load the file \"%v\".\nDoes it exist?  If so, is it in .png format?", path.Join(global.Directory, dir), err))
	}
	vao := createVao(img, &canvas)
	texture := createTexture(img)

	//the transformation matrix should be identity matrix at the beginning
	tmat := mgl32.Ident4().Add(mgl32.Mat4FromRows(
		mgl32.Vec4 { 0, 0, 0, canvas.X },
		mgl32.Vec4 { 0, 0, 0, canvas.Y },
		mgl32.Vec4 { 0, 0, 0, 0 },
		mgl32.Vec4 { 0, 0, 0, 0 },
	))
	smat := mgl32.Ident4()
	rmat := mgl32.Ident4()

	return Sprite{
		X:       x,
		Y:       y,
		vao:     vao,
		texture: texture,
		translation: tmat,
		scale: smat,
		rotation: rmat,
	}
}

func createImage(dir string) (*image.RGBA, error) {
	//read the image file, return empty rgba if error
	imgFile, err := os.Open(path.Join(global.Directory, dir))
	if err != nil {
		a := image.Rectangle{image.Point{0, 0}, image.Point{0, 0}}
		return image.NewRGBA(a), err
	}
	defer imgFile.Close()

	//decode the image file to use, return empty rgba is error
	img, err := png.Decode(imgFile)
	if err != nil {
		a := image.Rectangle{image.Point{0, 0}, image.Point{0, 0}}
		return image.NewRGBA(a), err
	}

	//create a rgba canvas to draw the stuff on
	rgba := image.NewRGBA(img.Bounds())

	//draw the image file based on image file
	draw.Draw(rgba, rgba.Bounds(), img, image.ZP, draw.Over)

	return rgba, nil
}

func createVao(img *image.RGBA, canvas *Canvas) uint32 {
	//calculate the image's x and y depending on image aspect ratio
	var x float32
	var y float32

	//The aspect ratio of the image and canvas
	ratio := float32(img.Rect.Size().X) / float32(img.Rect.Size().Y)
	cratio := canvas.GetAspectRatio()

	//the bigger side will cover the entire canvas's side, and the other side will get the size equal to the ratio of both ratios 
	//to ensure that the canvas shape will not affect the image ratio
	if ratio > 1.0 {
		ratio = 1 / ratio
		x = canvas.Width / 2
		y = canvas.Height * (ratio / cratio) / 2
	} else {
		x = canvas.Width * (ratio / cratio) / 2
		y = canvas.Height / 2
	}

	//create vertices based on the calculated x's and y's and the coordinate of image each vertices should be associated to
	//translate the vectors based on the canvas x's and y's
	var vec []float32 = []float32{
		x * -1, y, 0, 0.0, 0.0, //top left
		x, y, 0, 1.0, 0.0, //top right
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
	gl.BufferData(gl.ARRAY_BUFFER, len(vec)*4, gl.Ptr(vec), gl.STATIC_DRAW)

	//create ebo that will be used for the currently binded vao
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(ind)*4, gl.Ptr(ind), gl.STATIC_DRAW)

	//set the attributes
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
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

func (s *Sprite) GetTransformation() mgl32.Mat4 {
	ortho := mgl32.Ortho2D(-800, 800, -450, 450)

	return ortho.Mul4(s.translation.Mul4(s.rotation.Mul4(s.scale)))
}

func (s *Sprite) Move(x, y float64) {
	s.translation = mgl32.Translate3D(float32(x), float32(y), 0)
}

func (s *Sprite) Scale(v ...float64) {
	hasArgs := false
	var x float32
	var y float32
	for i, sc := range v {
		switch i {
		case 0:
			x = float32(sc)
			y = float32(sc)
			hasArgs = true
		case 1:
			y = float32(sc)
		default:
			panic("Invalid number of arguments.  Expected either 1 or 2 float32.  Found more")
		}
	}

	if !hasArgs {
		panic("Invalid number of arguments.  Could not find any arguments passed in")
	}

	s.scale = mgl32.Diag4(mgl32.Vec4{x, y, 1, 1})
}

func (s *Sprite) RadianRotate(angle float64) {
	s.rotation = mgl32.HomogRotate3DZ(float32(angle))
}

func (s *Sprite) AngleRotate(angle float64) {
	a := angle * (math.Pi / 180.0)
	s.rotation = mgl32.HomogRotate3DZ(float32(a))
}
