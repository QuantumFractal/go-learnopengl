package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v2.1/gl"
)

type Shader struct {
	id uint32
}

func CreateShader(vertexPath string, fragmentPath string) (*Shader, error) {
	shader := Shader{}

	vertexShader, err := compileShader(vertexPath, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fragmentShader, err := compileShader(fragmentPath, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	shader.id, err = createProgram(vertexShader, fragmentShader)
	if err != nil {
		return nil, err
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	return &shader, nil
}

func (s *Shader) Use() {
	gl.UseProgram(s.id)
}

func (s *Shader) setInt(name string, value int32) {
	gl.Uniform1i(gl.GetUniformLocation(s.id, gl.Str(name+"\x00")), value)
}

func (s *Shader) setFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(s.id, gl.Str(name+"\x00")), value)
}

func (s *Shader) setVec2(name string, value mgl32.Vec2) {
	gl.Uniform2fv(gl.GetUniformLocation(s.id, gl.Str(name+"\x00")), 1, &value[0])
}

func (s *Shader) setVec3(name string, value mgl32.Vec3) {
	gl.Uniform3fv(gl.GetUniformLocation(s.id, gl.Str(name+"\x00")), 1, &value[0])
}

func (s *Shader) setVec4(name string, value mgl32.Vec4) {
	gl.Uniform4fv(gl.GetUniformLocation(s.id, gl.Str(name+"\x00")), 1, &value[0])
}

func (s *Shader) setMat4(name string, value mgl32.Mat4) {
	gl.UniformMatrix4fv(gl.GetUniformLocation(s.id, gl.Str(name+"\x00")), 1, false, &value[0])
}

func (s *Shader) setBool(name string, value bool) {
	var intval int32
	if value == true {
		intval = 1
	} else {
		intval = 0
	}
	gl.Uniform1i(gl.GetUniformLocation(s.id, gl.Str(name+"\x00")), intval)
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
