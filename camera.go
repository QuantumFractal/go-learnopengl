package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

const YAW float32 = -90.0
const PITCH = 0.0
const SPEED = 2.5
const SENSITIVTY = 0.1
const ZOOM = 45.0

type Camera_Movement int

const (
	FORWARD Camera_Movement = iota
	BACKWARD
	LEFT
	RIGHT
)

// Camera structure.
type Camera struct {
	Position mgl32.Vec3
	Front    mgl32.Vec3
	Up       mgl32.Vec3
	Right    mgl32.Vec3
	WorldUp  mgl32.Vec3

	Yaw              float32
	Pitch            float32
	MovementSpeed    float32
	MouseSensitivity float32
	Zoom             float32
}

func CreateCamera(position mgl32.Vec3, up mgl32.Vec3, yaw float32, pitch float32) *Camera {
	camera := Camera{}
	camera.Position = position
	camera.WorldUp = up
	camera.Pitch = pitch

	return &camera
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}

func (c *Camera) ProcessKeyboard(direction Camera_Movement, deltaTime float32) {
	velocity := c.MovementSpeed * deltaTime
	if direction == FORWARD {
		c.Position = c.Position.Add(c.Front.Mul(velocity))
	}
	if direction == BACKWARD {
		c.Position = c.Position.Sub(c.Front.Mul(velocity))
	}
	if direction == LEFT {
		c.Position = c.Position.Sub(c.Right.Mul(velocity))
	}
	if direction == RIGHT {
		c.Position = c.Position.Add(c.Right.Mul(velocity))
	}
}

func (c *Camera) ProcessMouseMovement(xoffset float32, yoffset float32, constrainPitch bool) {
	xoffset *= c.MouseSensitivity
	yoffset *= c.MouseSensitivity

	c.Yaw += xoffset
	c.Pitch += yoffset

	if constrainPitch {
		if c.Pitch > 89.0 {
			c.Pitch = 89.0
		}
		if c.Pitch < -89.0 {
			c.Pitch = -89
		}
	}
}

func (c *Camera) updateCameraVectors() {
	var front mgl32.Vec3
	front.
}
