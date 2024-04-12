package main

import (
	oak "github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/event"
	random "math/rand"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/scene"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/collision"
	"image/color"
	"time"
)

var globalCtx *scene.Context

const (
	Asteroid collision.Label = 2
	Bullet collision.Label = 1
)

func shoot(pl *entities.Entity, _ key.Event) event.Response {
	bul := entities.New(globalCtx, entities.WithLabel(Bullet), entities.WithColor(color.White), entities.WithRect(floatgeom.NewRect2WH(0, 0, 5, 10)), entities.WithPosition(floatgeom.Point2{pl.X(), pl.Y()}))
	bul.Delta[1] = -4
	event.Bind(globalCtx, event.Enter, bul, bulletLoop)

	return 0
}

func bulletLoop(bul *entities.Entity, _ event.EnterPayload) event.Response {
	if collision.HitLabel(bul.Space, Asteroid) != nil {
		go func() {
			time.Sleep(time.Millisecond * 20);
			bul.Destroy()
		}()
	}
	bul.ShiftDelta()

	return 0
}

var score int = 0

func playerLoop(pl *entities.Entity, _ event.EnterPayload) event.Response {
	if oak.IsDown(key.D) {
		pl.ShiftPos(4, 0)
	}
	if oak.IsDown(key.A) {
		pl.ShiftPos(-4, 0)
	}
	return 0
}

func asteroidLoop(ast *entities.Entity, _ event.EnterPayload) event.Response {
	if collision.HitLabel(ast.Space, Bullet) != nil || ast.Y() > 500 {
		go func() {
			spawn := floatgeom.Point2{float64(random.Intn(550)), float64(random.Intn(50) - 150)}
			time.Sleep(time.Millisecond * 20);
			ast.SetPos(spawn)
			score += 10
		}()
	}
	ast.ShiftDelta()

	return 0
}

func startScene(ctx *scene.Context) {
	globalCtx = ctx
	player := entities.New(ctx, entities.WithColor(color.White), entities.WithRect(floatgeom.NewRect2WH(0, 0, 20, 20)), entities.WithPosition(floatgeom.Point2{250, 460}))
	event.Bind(ctx, event.Enter, player, playerLoop)
	event.Bind(ctx, key.Down(key.Spacebar), player, shoot)
	for range 5 {
		ast := entities.New(ctx, entities.WithLabel(Asteroid), entities.WithRect(floatgeom.NewRect2WH(0, 0, 64, 64)), entities.WithColor(color.White), entities.WithPosition(floatgeom.Point2{float64(random.Intn(550)), float64(random.Intn(50) - 150)}))
		ast.Delta[1] = 2
		event.Bind(ctx, event.Enter, ast, asteroidLoop)
	}
}

func main() {
	mainScene := scene.Scene{Start: startScene}
	oak.AddScene("Main", mainScene)
	oak.Init("Main")
}
