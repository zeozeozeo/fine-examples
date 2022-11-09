// A simple example of integrating physics into your game.
// This uses the Go port of Chipmunk2D: https://github.com/jakecoffman/cp.
// Hold the left mouse button to spawn new bodies.
package main

import (
	"image/color"
	"math/rand"

	"github.com/jakecoffman/cp"
	"github.com/zeozeozeo/fine"
)

var (
	space       *cp.Space
	bodies      []*PhysicsBody
	spawnAmount int     = 2
	ballRadius  float64 = 5
)

type PhysicsBody struct {
	Body   *cp.Body
	Shape  *cp.Shape
	Radius float64
	Entity *fine.Entity
}

func main() {
	app := fine.NewApp("Physics", 1280, 720).
		SetUpdateFunc(update)
	initSpace()
	makeBounds(app)

	// Start the application
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func update(dt float64, app *fine.App) {
	if app.IsMouseButtonDown(fine.MBUTTON_LEFT) {
		mX, mY := app.GetMousePos()

		for i := 0; i < spawnAmount; i++ {
			jitter := fine.NewVec2(rand.Float64()*10, rand.Float64()*10)
			worldPos := screenToWorld(mX, mY, app)
			spawnBall(app, worldPos.Add(jitter), ballRadius)
		}
	}

	space.Step(dt)
	draw(app)
}

func draw(app *fine.App) {
	for idx, body := range bodies {
		bodyPos := body.Body.Position()
		screenX, screenY := app.Camera.WorldToScreen(fine.NewVec2(bodyPos.X, -bodyPos.Y))

		// Destroy off-screen bodies
		if float64(screenX)+body.Radius > float64(app.Width) || float64(screenY)+body.Radius > float64(app.Height) {
			space.RemoveBody(body.Body)
			space.RemoveShape(body.Shape)
			body.Entity.Destroy()

			if idx < len(bodies) {
				bodies[idx] = bodies[len(bodies)-1]
				bodies = bodies[:len(bodies)-1]
			}
			continue
		}

		body.Entity.Position.X, body.Entity.Position.Y = bodyPos.X, -bodyPos.Y
	}
}

func initSpace() {
	space = cp.NewSpace()
	space.Iterations = 1
	space.SetGravity(cp.Vector{X: 0, Y: -9.81 * 100})

	// Generally, you wouldn't be using a spatial hash, but if
	// you have a lot of physics bodies (like here), it is a great
	// performance improvement.
	space.UseSpatialHash(2.0, 10000)
}

func makeBounds(app *fine.App) {
	lineStart1, lineEnd1 := fine.NewVec2(-400, 0), fine.NewVec2(-20, 200)
	lineStart2, lineEnd2 := fine.NewVec2(400, 0), fine.NewVec2(20, 200)
	app.Line(lineStart1, lineEnd1, color.RGBA{255, 255, 255, 255}, true)
	app.Line(lineStart2, lineEnd2, color.RGBA{255, 255, 255, 255}, true)

	space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: lineStart1.X, Y: -lineStart1.Y}, cp.Vector{X: lineEnd1.X, Y: -lineEnd1.Y}, 0))
	space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: lineStart2.X, Y: -lineStart2.Y}, cp.Vector{X: lineEnd2.X, Y: -lineEnd2.Y}, 0))
}

func spawnBall(app *fine.App, position fine.Vec2, radius float64) {
	// Create the physics body
	body := cp.NewBody(1.0, cp.INFINITY)
	body.SetPosition(cp.Vector{X: position.X, Y: -position.Y})

	shape := cp.NewCircle(body, radius, cp.Vector{})
	shape.SetElasticity(0)
	shape.SetFriction(0)

	space.AddBody(shape.Body())
	space.AddShape(shape)

	// Create the entity
	// TODO: Draw circles instead of rectangles
	entity := app.Rect(
		position,
		radius*2,
		radius*2,
		color.RGBA{255, 255, 0, 255},
		true,
	)
	pBody := &PhysicsBody{
		Body:   body,
		Shape:  shape,
		Radius: radius,
		Entity: entity,
	}

	bodies = append(bodies, pBody)
}

func screenToWorld(x, y int, app *fine.App) fine.Vec2 {
	zx := (float64(x) * app.Camera.Zoom) - (float64(app.Width) / 2)
	zy := (float64(y) * app.Camera.Zoom) - (float64(app.Height) / 2)
	return fine.NewVec2(zx-app.Camera.Position.X, zy-app.Camera.Position.Y)
}
