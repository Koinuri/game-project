package framework
/*
import (
	"github.com/go-gl/mathgl/mgl32"
)
*/

type Object struct {
	canvas Canvas
	artists map[string]Artist
	x float32
	y float32
	ox float32
	oy float32
	scalex float32
	scaley float32
	angle float32
}

//Initializes the object object that stores artist objects.
//	*InitObject(*Canvas, 
/*
func InitObject(i ...interface{}) Object {
	var width float32
	var height float32
	var x float32
	var y float32

	c := InitCanvas(width, height, x, y)
	a := make(map[string]*Artist)

	return Object {
		c,
		a,
	}
}

func (o *Object) CreateSprite(i ...interface{}) {
	var name string
	var dir string
	var origin origin

	sprite := InitSprite(dir, o.canvas, origin)

	artists[name] = sprite
}
*/

func (o *Object) GetArtists() []Artist {
	artists := make([]Artist, len(o.artists))
	i := 0

	for _, artist := range o.artists {
		artists[i] = artist.applyTransformations(o.x, o.y, o.scalex, o.scaley, o.angle)
		i++
	}

	return artists
}
