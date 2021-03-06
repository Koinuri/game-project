package framework

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	vertexShaderSource = `
        #version 400
        in vec3 vp;
		in vec2 tx;

		uniform mat4 transformation;
		uniform mat4 projection;

		out vec2 TexCoord;

        void main() {
            gl_Position = projection * transformation * vec4(vp, 1.0);
			TexCoord = tx;
        }
    ` + "\x00"
	fragmentShaderSource = `
        #version 400
		in vec2 TexCoord;

		out vec4 frag_colour;
		
		uniform sampler2D ourTexture;

        void main() {
            frag_colour = texture(ourTexture, TexCoord);
        }
    ` + "\x00"
)

func Init(width int, height int) (*glfw.Window, uint32) {
	window := initGlfw(width, height)
	prog := initOpenGL()

	gl.UseProgram(prog)

	ortho := mgl32.Ortho2D(-800, 800, -450, 450)
	orthoUniform := gl.GetUniformLocation(prog, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(orthoUniform, 1, false, &ortho[0])

	return window, prog
}

func initGlfw(width int, height int) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
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

func InitFrame() {
	clear()
}

func SwapWindowAndPollEvents(window *glfw.Window) {
	glfw.PollEvents()
	window.SwapBuffers()
}

func Draw(objects []Artist, prog uint32) {
	gl.UseProgram(prog)

	for _, obj := range objects {
		vao, texture := obj.GetDrawInfo()
		transformation := obj.GetTransformation()

		tUniform := gl.GetUniformLocation(prog, gl.Str("transformation\x00"))
		gl.UniformMatrix4fv(tUniform, 1, false, &transformation[0])

		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	}
}

func clear() {
	gl.ClearColor(0.1, 0.2, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
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

func Clean() {
	glfw.Terminate()
}
