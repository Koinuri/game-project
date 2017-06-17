package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"os"
	"path"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/koinuri/game-project/main/framework"
	"github.com/koinuri/game-project/main/global"
)

const (
	width              = 1366
	height             = 768
	vertexShaderSource = `
        #version 450
        in vec3 vp;
		in vec2 tx;

		out vec2 TexCoord;

        void main() {
            gl_Position = vec4(vp, 1.0);
			TexCoord = tx;
        }
    ` + "\x00"
	fragmentShaderSource = `
        #version 450
		in vec2 TexCoord;

		out vec4 frag_colour;
		
		uniform sampler2D ourTexture;

        void main() {
            frag_colour = texture(ourTexture, TexCoord);
        }
    ` + "\x00"
)

func Init() {
	runtime.LockOSThread()

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	global.Directory = path.Dir(ex)
}

func main() {
	Init()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	spr := framework.InitSprite("kotonoha-7.png", framework.TopLeft)

	for !window.ShouldClose() {
		draw(spr, window, program)
	}
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Game Project", nil, nil)

	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	return window
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func draw(s framework.Sprite, window *glfw.Window, prog uint32) {
	vao, texture := s.GetDrawInfo()
	gl.ClearColor(0.1, 0.2, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(prog)

	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	glfw.PollEvents()
	window.SwapBuffers()
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
