package framework

import (
	"math"
	"fmt"
)

type Object struct {
	artists []Artist
	names   []string
	x       float32
	y       float32
	width   float32
	height  float32
	ox      float32
	oy      float32
	scalex  float32
	scaley  float32
	angle   float32
	bounds	Bound
	origin origin
}

//Initializes the object object that stores artist objects.
//	*InitObject()
// 	*InitObject(origin)
//	*InitObject(width, height)
//	*InitObject(width, height, origin)
//	*InitObject(width, height, x, y)
//	*InitObject(width, height, x, y, origin)
//Where:
//	width and height are the width and height of the 
func InitObject(i ...interface{}) Object {
	var width float32 = 1600
	var height float32 = 900
	var x float32
	var y float32
	var or origin = 4

	a := make([]Artist, 0)
	n := make([]string, 0)

	countFloats := 0
	for _, item := range i {
		u, isOrigin := item.(origin)
		if isOrigin {
			or = u
		}

		r, isFloat := item.(float32)
		if isFloat {
			switch (countFloats) {
			case 0:
				width = r
				countFloats++
			case 1:
				height = r
				countFloats++
			case 2:
				x = r
				countFloats++
			case 3:
				y = r
				countFloats++
			default:
				panic(fmt.Sprintf("Invalid number of arguments.  Expected 2 or 4 float32, found %v.", countFloats))
			}
		}
	}

	if countFloats == 1 || countFloats == 3 {
		panic(fmt.Sprintf("Invalid number of arguments.  Expected 2 or 4 float32, found %v.", countFloats))
	}

	obj := Object {
		a,
		n,
		x,
		y,
		width,
		height,
		0,
		0,
		1,
		1,
		0,
		Bound {0, 0, 0, 0},
		or,
	}

	obj.updateOrigin()

	obj.setBounds()

	return obj
}

func (o *Object) updateOrigin() {
	sx := o.width / 2
	sy := o.height / 2

	//based on the origin, get the constant to place origin (negative, 0, or positive)
	detx := math.Mod(float64(o.origin), 3.0) - 1
	dety := math.Floor(float64(o.origin)/3.0) - 1

	//origin at angle 0
	ox := o.x - sx*float32(detx)
	oy := o.y - sy*float32(dety)

	//difference between (x,y) and (ox,oy) and distance based on the numbers calculated
	dx := o.x - ox
	dy := o.y - oy
	d := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	//The real origin using the angle provided
	if o.angle != 0 {
		ox = d*float32(math.Cos(float64(o.angle))) + o.x
		oy = d*float32(math.Sin(float64(o.angle))) + o.y
	}

	o.ox = ox
	o.oy = oy
}

func (o *Object) setBounds() {
	left := o.ox - (o.width / 2)
	right := o.ox + (o.width / 2)
	up := o.oy + (o.height / 2)
	down := o.oy + (o.height / 2)

	o.bounds = Bound {left, right, up, down}
}

//Creates a sprite to be added to the object then returns the reference to that sprite
//	*CreateSprite(name, directory)
//	*CreateSprite(name, directory, origin)
//Where
//	name is the name of the sprite that can be referred to
//	directory is the location of the image to create sprites with
//	origin is the origin the sprite will be based on.
func (o *Object) CreateSprite(i ...interface{}) *Sprite {
	var name string
	var dir string
	var or origin = 4

	if len(i) > 0 {
		test, succ := i[0].(string)
		if !succ {
			panic(fmt.Sprintf("Invalid argument.  Expected string, found %T", i[0]))
		} else {
			name = test
		}

		for _, n := range o.names {
			if name == n {
				panic(fmt.Sprintf("Invalid argument.  The name \"%v\" already exists", name))
			}
		}
	}

	if len(i) > 1 {
		test, succ := i[1].(string)
		if !succ {
			panic(fmt.Sprintf("Invalid argument.  Expected string, found %T", i[1]))
		} else {
			dir = test
		}
	}

	if len(i) > 2 {
		test, succ := i[2].(origin)
		if !succ {
			panic(fmt.Sprintf("Invalid argument.  Expected origin, found %T", i[2]))
		} else {
			or = test
		}
	}

	if len(i) == 0 || len(i) > 3 {
		panic(fmt.Sprintf("Invalid number of arguments.  The method only accepts one to three arguments, found %v instead.", len(i)))
	}

	sprite := InitSprite(dir, InitCanvas(o.width, o.height, o.ox, o.oy), or)

	o.artists = append(o.artists, &sprite)
	o.names = append(o.names, name)

	return &sprite
}

func (o *Object) GetArtists() []Artist {
	artists := make([]Artist, len(o.artists))

	for i, artist := range o.artists {
		artists[i] = artist.applyTransformations(o.ox, o.oy, o.scalex, o.scaley, o.angle)
	}

	return artists
}

func (o *Object) GetArtist(name string) Artist {
	for i, n := range o.names {
		if n == name {
			return o.artists[i]
		}
	}
	panic(fmt.Sprintf("Invalid name.  Could not find any artist with the name \"%v\".", name))
}

func (o *Object) Move(x, y float64) {
	o.x = float32(x)
	o.y = float32(y)

	o.updateOrigin()
}

func (o *Object) Scale(i ...float64) {
	var x float32 = 1.0
	var y float32 = 1.0

	for a, b := range i {
		switch (a) {
		case 0:
			x = float32(b)
			y = float32(b)
		case 1:
			y = float32(b)
		default:
			panic(fmt.Sprintf("Invalid argument.  Expected less than two arguments, found more"))
		}
	}

	if len(i) == 0 {
		panic(fmt.Sprintf("Invalid argument.  Did not find any float value"))
	}

	o.scalex = x
	o.scaley = y
}

func (o *Object) AngleRotate(angle float64) {
	a := angle * (math.Pi / 180)
	o.angle = float32(a)
}

func (o *Object) RadianRotate(angle float64) {
	o.angle = float32(angle)
}
