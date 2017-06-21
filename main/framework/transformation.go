package framework

import "github.com/go-gl/mathgl/mgl32"

type transformation struct {
	translation mgl32.Mat4
	scale       mgl32.Mat4
	rotation    mgl32.Mat4
}

func InitTransformation() transformation {
	return transformation{
		mgl32.Ident4(),
		mgl32.Ident4(),
		mgl32.Ident4(),
	}
}
