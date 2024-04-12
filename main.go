package main

import (
	oak "github.com/oakmound/oak/v4"
	"os"
	"github.com/oakmound/oak/v4/event"
	"strconv"
	random "math/rand"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/collision"
	"image/color"
	"image"
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

var score string = "0"
var lives string = "3"
var diff int = 1

func playerLoop(pl *entities.Entity, _ event.EnterPayload) event.Response {
	if oak.IsDown(key.D) {
		pl.ShiftPos(float64(4 * ((diff + 150) / 150)), 0)
	}
	if oak.IsDown(key.A) {
		pl.ShiftPos(-float64(4 * ((diff + 150) / 150)), 0)
	}
	return 0
}

func asteroidLoop(ast *entities.Entity, _ event.EnterPayload) event.Response {
	ast.Delta[1] = float64(1.2 * float64((diff + 150) / 150))
	if collision.HitLabel(ast.Space, Bullet) != nil {
		go func() {
			spawn := floatgeom.Point2{float64(random.Intn(600)), float64(random.Intn(50) - 150)}
			time.Sleep(time.Millisecond * 20);
			ast.SetPos(spawn)
			int_score, _ := strconv.Atoi(score)
			int_score += 5 * (diff + 50) / 50
			score = strconv.Itoa(int_score)
			diff++
		}()
	}
	if ast.Y() > 360 {
		go func() {
			spawn := floatgeom.Point2{float64(random.Intn(600)), float64(random.Intn(50) - 150)}
			ast.SetPos(spawn)
			int_lives, _ := strconv.Atoi(lives)
			int_lives--
			lives = strconv.Itoa(int_lives)
			if int_lives < 0 {
				os.Exit(228)
			}
		}()
	}
	ast.ShiftDelta()

	return 0
}

func getImageFromFilePath(filePath string) image.Image {
    f, _ := os.Open(filePath)
    defer f.Close()
    image, _, _ := image.Decode(f)
    return image
}

func startScene(ctx *scene.Context) {
	oak.UpdateViewSize(640, 360)
	globalCtx = ctx
	ctx.Draw(render.NewStrPtrText(&score, 10, 10))
	ctx.Draw(render.NewStrPtrText(&lives, 10, 25))
	player := entities.New(ctx, entities.WithColor(color.White), entities.WithRect(floatgeom.NewRect2WH(0, 0, 20, 20)), entities.WithPosition(floatgeom.Point2{250, 340}))
	event.Bind(ctx, event.Enter, player, playerLoop)
	event.Bind(ctx, key.Down(key.Spacebar), player, shoot)
	for range 5 {
		ast := entities.New(ctx, entities.WithLabel(Asteroid), entities.WithRect(floatgeom.NewRect2WH(0, 0, 64, 64)), entities.WithColor(color.White), entities.WithPosition(floatgeom.Point2{float64(random.Intn(600)), float64(random.Intn(50) - 150)}))
		ast.Delta[1] = 2
		event.Bind(ctx, event.Enter, ast, asteroidLoop)
	}
}

func main() {
	mainScene := scene.Scene{Start: startScene}
	oak.AddScene("Main", mainScene)
	oak.Init("Main")
}
