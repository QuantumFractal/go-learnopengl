package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// FloatSize represents how large a float is.
const FloatSize = 4

func init() {
	runtime.LockOSThread()
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

func compileShader(shaderFile string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	data, err := ioutil.ReadFile(shaderFile)
	if err != nil {
		panic(err)
	}

	csources, free := gl.Strs(string(data) + "\x00")
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

		return 0, fmt.Errorf("failed to compile shader\n %v \n because:\n %v", string(data), log)
	}

	return shader, nil
}

func createProgram(vertexShader uint32, fragmentShader uint32) (uint32, error) {
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}
	return shaderProgram, nil
}

func createTexture(textureFile string) (uint32, error) {
	imgFile, err := os.Open(textureFile)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found: %v", textureFile, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TextureParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TextureParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TextureParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TextureParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	gl.GenerateMipmap(gl.TEXTURE_2D)

	return texture, nil
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "triangles", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	var vertices = []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}

	var indices = []int32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}

	// VBO, EBO, VAO creation
	var (
		VBO uint32
		EBO uint32
		VAO uint32
	)

	// Generate the buffers
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)
	gl.GenVertexArrays(1, &VAO)

	// Bind Vertex Array Object
	gl.BindVertexArray(VAO)

	// Copy vertices array into vertex buffer object
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	// Copy index array into element buffer object
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	// position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// texture attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteBuffers(1, &VBO)

	// Load Texture
	texture, err := createTexture("tarp.png")
	if err != nil {
		panic(err)
	}

	ball, err := createTexture("ball.png")
	if err != nil {
		panic(err)
	}

	// Load Shaders
	vertexShader, err := compileShader("vertex.glsl", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := compileShader("fragment.glsl", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	defer gl.DeleteShader(fragmentShader)

	shaderProgram, err := createProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	defer gl.DeleteProgram(shaderProgram)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.MULTISAMPLE)

	gl.UseProgram(shaderProgram)
	gl.Uniform1i(gl.GetUniformLocation(shaderProgram, gl.Str("texture1\x00")), 0)
	gl.Uniform1i(gl.GetUniformLocation(shaderProgram, gl.Str("texture2\x00")), 1)

	var (
		model      mgl32.Mat4
		view       mgl32.Mat4
		projection mgl32.Mat4
	)

	model = mgl32.HomogRotate3D(mgl32.DegToRad(-55.0), mgl32.Vec3{1.0, 0.0, 0.0})
	view = mgl32.Translate3D(0.0, 0.0, -3.0)
	projection = mgl32.Perspective(45.0, 800/600, 1.0, 100.0)

	for !window.ShouldClose() {
		// Process input
		processInput(window)

		// Render stuff
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)

		model = mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{1, 1, 0.0}.Normalize())

		gl.UseProgram(shaderProgram)
		modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model\x00"))
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		viewUniform := gl.GetUniformLocation(shaderProgram, gl.Str("view\x00"))
		gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

		projectionUniform := gl.GetUniformLocation(shaderProgram, gl.Str("projection\x00"))
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

		gl.BindVertexArray(VAO)

		// Draw triangles
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, ball)

		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		//gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		// Swap buffers
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
