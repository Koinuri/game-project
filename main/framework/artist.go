package framework

import "github.com/go-gl/mathgl/mgl32"

//Interface that is required
type Artist interface {
	GetDrawInfo() (uint32, uint32)
	GetTransformation() mgl32.Mat4
}
