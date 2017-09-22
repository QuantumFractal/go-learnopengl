package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"
	"runtime"

	"github.com/disintegration/imaging"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// FloatSize represents how large a float is.
const FloatSize = 4

var (
	deltaTime float32
	lastFrame float32
)

var lastX = float32(800 / 2)
var lastY = float32(600 / 2)
var firstMouse = true
var camera = CreateCamera(mgl32.Vec3{0.0, 0.0, 3.0}, mgl32.Vec3{0.0, 1.0, 0.0}, YAW, PITCH)

func init() {
	runtime.LockOSThread()
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		camera.ProcessKeyboard(FORWARD, deltaTime)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camera.ProcessKeyboard(BACKWARD, deltaTime)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		camera.ProcessKeyboard(LEFT, deltaTime)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		camera.ProcessKeyboard(RIGHT, deltaTime)
	}
	if window.GetKey(glfw.KeyE) == glfw.Press {
		camera.ProcessKeyboard(UP, deltaTime)
	}
	if window.GetKey(glfw.KeyQ) == glfw.Press {
		camera.ProcessKeyboard(DOWN, deltaTime)
	}
}

func mouseCallback(window *glfw.Window, xpos float64, ypos float64) {
	if firstMouse {
		lastX = float32(xpos)
		lastY = float32(ypos)
		firstMouse = false
	}

	xoffset := float32(xpos) - lastX
	yoffset := float32(ypos) - lastY

	lastX = float32(xpos)
	lastY = float32(ypos)
	camera.ProcessMouseMovement(xoffset, -yoffset, true)
}

func scrollCallback(window *glfw.Window, xoffset float64, yoffset float64) {
	camera.ProcessMouseScroll(float32(yoffset))
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
	flipped := imaging.FlipV(img)
	draw.Draw(rgba, rgba.Bounds(), flipped, image.Point{0, 0}, draw.Src)

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

func getTime() float32 {
	return float32(glfw.GetTime())
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
	//window.SetPos(1980, 60)
	window.MakeContextCurrent()
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(mouseCallback)
	window.SetScrollCallback(scrollCallback)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	vertices := []float32{
		// positions          // normals           // texture coords
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,

		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, -0.5, -1.0, 0.0, 0.0, 1.0, 1.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, -1.0, 0.0, 0.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 1.0, 1.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
	}

	// VBO, EBO, VAO creation
	var (
		VBO      uint32
		VAO      uint32
		lightVAO uint32
	)

	// Generate the bufers
	gl.GenBuffers(1, &VBO)
	gl.GenVertexArrays(1, &VAO)

	// Copy vertices array into vertex buffer object
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 3*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	// Bind Vertex Array Object
	gl.BindVertexArray(VAO)

	// position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// normal attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// texture attribute
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteBuffers(1, &VBO)

	gl.GenVertexArrays(1, &lightVAO)
	gl.BindVertexArray(lightVAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	lightingShader, err := CreateShader("vertex.glsl", "fragment.glsl")
	if err != nil {
		panic(err)
	}

	lampShader, err := CreateShader("vertex.glsl", "lightfragment.glsl")
	if err != nil {
		panic(err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.MULTISAMPLE)

	var (
		view       mgl32.Mat4
		projection mgl32.Mat4
		model      mgl32.Mat4
	)

	model = mgl32.Ident4()
	lightPos := mgl32.Vec3{1.2, 1.0, 2.0}
	lightingShader.setInt("material.diffuse", 0)
	lightingShader.setInt("material.specular", 1)

	boxTexture, err := createTexture("diffuse_crate.png")
	if err != nil {
		panic(err)
	}

	boxSpecular, err := createTexture("specular_crate.png")
	if err != nil {
		panic(err)
	}

	// MAIN LOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOP
	for !window.ShouldClose() {
		currentFrame := float32(glfw.GetTime())
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		// Process input
		processInput(window)

		// Render stuff
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		lightingShader.Use()
		lightingShader.setVec3("light.position", lightPos)
		lightingShader.setVec3("viewPos", camera.Position)

		diffuseColor := mgl32.Vec3{0.5, 0.5, 0.5}
		ambientColor := mgl32.Vec3{0.2, 0.2, 0.2}

		lightingShader.setVec3("light.ambient", ambientColor)
		lightingShader.setVec3("light.diffuse", diffuseColor)
		lightingShader.setVec3("light.specular", mgl32.Vec3{1.0, 1.0, 1.0})

		lightingShader.setFloat("material.shininess", 64.0)

		projection = mgl32.Perspective(mgl32.DegToRad(camera.Zoom), float32(800)/float32(600), 1.0, 100.0)
		view = camera.GetViewMatrix()
		lightingShader.setMat4("projection", projection)
		lightingShader.setMat4("view", view)
		lightingShader.setMat4("model", model)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, boxTexture)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, boxSpecular)

		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		lampShader.Use()
		lampShader.setMat4("projection", projection)
		lampShader.setMat4("view", view)
		model := mgl32.Translate3D(lightPos.X(), lightPos.Y(), lightPos.Z())
		model = model.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
		lampShader.setMat4("model", model)

		gl.BindVertexArray(lightVAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		// Swap buffers
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
