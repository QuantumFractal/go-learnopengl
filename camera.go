package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// YAW thing.
const YAW float32 = -90.0

// PITCH thing.
const PITCH = 0.0

// SPEED thing.
const SPEED = 2.5

const SENSITIVTY = 0.1
const ZOOM = 45.0

const FORWARD = 0
const BACKWARD = 1
const LEFT = 2
const RIGHT = 3
const UP = 4
const DOWN = 5

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

// CreateCamera thing.
func CreateCamera(position mgl32.Vec3, up mgl32.Vec3, yaw float32, pitch float32) *Camera {
	camera := Camera{}
	camera.Position = position
	camera.WorldUp = up
	camera.Yaw = yaw
	camera.Pitch = pitch
	camera.MovementSpeed = SPEED
	camera.MouseSensitivity = SENSITIVTY
	camera.Zoom = ZOOM
	camera.updateCameraVectors()
	return &camera
}

// GetViewMatrix thing.
func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}

// ProcessKeyboard thing.
func (c *Camera) ProcessKeyboard(direction int, deltaTime float32) {
	velocity := c.MovementSpeed * deltaTime
	if direction == UP {
		c.Position = c.Position.Add(c.WorldUp.Mul(velocity))
	}
	if direction == DOWN {
		c.Position = c.Position.Sub(c.WorldUp.Mul(velocity))
	}
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

// ProcessMouseMovement thing.
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
	c.updateCameraVectors()
}

// ProcessMouseScroll thing.
func (c *Camera) ProcessMouseScroll(yoffset float32) {
	if c.Zoom >= 1.0 && c.Zoom <= 45.0 {
		c.Zoom -= yoffset
	}
	if c.Zoom <= 1.0 {
		c.Zoom = 1.0
	}
	if c.Zoom >= 45.0 {
		c.Zoom = 45.0
	}
}

func (c *Camera) updateCameraVectors() {
	frontX := math.Cos(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch)))
	frontY := math.Sin(float64(mgl32.DegToRad(c.Pitch)))
	frontZ := math.Sin(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch)))
	c.Front = mgl32.Vec3{float32(frontX), float32(frontY), float32(frontZ)}.Normalize()
	c.Right = c.Front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}
