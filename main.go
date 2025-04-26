package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"rpg-tutorial/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	player               *entities.Player
	enemies              []*entities.Enemy
	potions              []*entities.Potion
	tileMapJSON          *TileMapJSON
	tileSets             []TileSet
	connectedControllers []ebiten.GamepadID
	GamepadID            ebiten.GamepadID
	camera               *Camera
	colliders            []image.Rectangle
}

const (
	KB_SPEED  = 2
	GP_SPEED  = 2.5
	DEAD_ZONE = 0.1
	TILE_SIZE = 16
	ENEMY_SPEED = 0.8
)

var (
	// Gamepad Button Constants - for readability
	GamepadButtonRight = ebiten.StandardGamepadButtonLeftRight
	GamepadButtonLeft  = ebiten.StandardGamepadButtonLeftLeft
	GamepadButtonUp    = ebiten.StandardGamepadButtonLeftTop
	GamepadButtonDown  = ebiten.StandardGamepadButtonLeftBottom

	// Fallback Gamepad Buttons (if standard layout is not available) - use with caution, may vary by controller
	FallbackButtonRight = ebiten.GamepadButton(15)
	FallbackButtonLeft  = ebiten.GamepadButton(14)
	FallbackButtonUp    = ebiten.GamepadButton(12)
	FallbackButtonDown  = ebiten.GamepadButton(13)
)

func (g *Game) detectAndSelectGamepad() {
	// Detect connected gamepads and store their IDs.
	g.connectedControllers = g.connectedControllers[:0] // Reset the slice
	g.connectedControllers = ebiten.AppendGamepadIDs(g.connectedControllers)

	// Use the first available gamepad if any are connected.
	if len(g.connectedControllers) > 0 {
		g.GamepadID = g.connectedControllers[0]
		// Optional: Print gamepad information for debugging (uncomment to use).
		// fmt.Println("Connected Gamepad SDL ID:", ebiten.GamepadSDLID(g.GamepadID))
		// fmt.Println("Connected Gamepad Name:", ebiten.GamepadName(g.GamepadID))

	} else {
		g.GamepadID = -1 // No controller connected
	}
}

func (g *Game) handleGamepadInput() {
	if g.GamepadID >= 0 {
		//handle buttons
		if ebiten.IsStandardGamepadLayoutAvailable(g.GamepadID) {
			g.handleStandardGamepadButtons()
		} else {
			g.handleNonStandardGamepadButtons()
		}
		// Analog Stick Movement (Left Stick)
		g.handleAnalogStickMovement()
	}
}

func (g *Game) handleStandardGamepadButtons() {
	if ebiten.IsStandardGamepadButtonPressed(g.GamepadID, GamepadButtonRight) {
		g.player.Dx = GP_SPEED
	}
	if ebiten.IsStandardGamepadButtonPressed(g.GamepadID, GamepadButtonLeft) {
		g.player.Dx = -GP_SPEED
	}
	if ebiten.IsStandardGamepadButtonPressed(g.GamepadID, GamepadButtonUp) {
		g.player.Dy = -GP_SPEED
	}
	if ebiten.IsStandardGamepadButtonPressed(g.GamepadID, GamepadButtonDown) {
		g.player.Dy = GP_SPEED
	}
}

func (g *Game) handleNonStandardGamepadButtons() {
	if ebiten.IsGamepadButtonPressed(g.GamepadID, FallbackButtonRight) { // Right
		g.player.Dx = GP_SPEED
	}
	if ebiten.IsGamepadButtonPressed(g.GamepadID, FallbackButtonLeft) { // Left
		g.player.Dx = -GP_SPEED
	}
	if ebiten.IsGamepadButtonPressed(g.GamepadID, FallbackButtonUp) { // Up
		g.player.Dy = -GP_SPEED
	}
	if ebiten.IsGamepadButtonPressed(g.GamepadID, FallbackButtonDown) { // Down
		g.player.Dy = GP_SPEED
	}
}

func isOutsideDeadZone(value float64, deadZone float64) bool {
	return value > deadZone || value < -deadZone
}

func (g *Game) handleAnalogStickMovement() {
	xAxis := ebiten.GamepadAxisValue(g.GamepadID, 0)
	yAxis := ebiten.GamepadAxisValue(g.GamepadID, 1)

	if isOutsideDeadZone(xAxis, DEAD_ZONE) {
		g.player.Dx = xAxis * GP_SPEED
	}
	if isOutsideDeadZone(yAxis, DEAD_ZONE) {
		g.player.Dy = yAxis * GP_SPEED
	}
}

func (g *Game) handleKeyboardControls() {
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.Dx = KB_SPEED
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.Dx = -KB_SPEED
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Dy = -KB_SPEED
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Dy = KB_SPEED
	}
}

func (g *Game) handleMovement() {
	// Reset movement
	g.player.Dx = 0
	g.player.Dy = 0

	// Handle input
	g.handleKeyboardControls()
	g.handleGamepadInput()

	// Move player
	g.player.X += g.player.Dx
	g.player.Y += g.player.Dy
}

func CheckCollisionHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+16.0,
				int(sprite.Y)+16.0,
			),
		) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - 16.0
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+16.0,
				int(sprite.Y)+16.0,
			),
		) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - 16.0
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.X)
			}
		}
	}
}

func (g *Game) Update() error {

	// detect controllers
	g.detectAndSelectGamepad()
	g.handleMovement()

	//collision
	CheckCollisionHorizontal(g.player.Sprite, g.colliders)
	CheckCollisionVertical(g.player.Sprite, g.colliders)

	//spawning enemy
	for _, enemy := range g.enemies {
		enemy.Dx = 0.00
		enemy.Dy = 0.00
		if enemy.FollowsPlayer {
			if enemy.X < g.player.X {
				enemy.Dx += 0.8
			} else if enemy.X > g.player.X {
				enemy.Dx -= 0.8
			}
			if enemy.Y < g.player.Y {
				enemy.Dy += 0.8
			} else if enemy.Y > g.player.Y {
				enemy.Dy -= 0.8
			}
		}
		enemy.X += enemy.Dx
		CheckCollisionHorizontal(enemy.Sprite, g.colliders)
		enemy.Y += enemy.Dy
		CheckCollisionVertical(enemy.Sprite, g.colliders)
	}

	for _, potion := range g.potions {
		if g.player.X > potion.X {
			g.player.Health += potion.HealAmount
		}
	}

	g.camera.FollowTarget(g.player.X+TILE_SIZE/2, g.player.Y+TILE_SIZE/2, 320, 240)
	g.camera.Constrain(
		float64(g.tileMapJSON.Layers[0].Width)*TILE_SIZE,
		float64(g.tileMapJSON.Layers[0].Height)*TILE_SIZE,
		320,
		240,
	)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	screen.Fill(color.RGBA{120, 180, 255, 255})
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %.2f Y: %.2f", g.player.X, g.player.Y))
	//start drawing map
	opts := ebiten.DrawImageOptions{}

	// loop over the layers
	for layerIndex, layer := range g.tileMapJSON.Layers {
		// loop over the tiles in the layer data
		for index, id := range layer.Data {

			if id == 0 {
				continue
			}

			// get the tile position of the tile
			x := index % layer.Width
			y := index / layer.Width

			// convert the tile position to pixel position
			x *= TILE_SIZE
			y *= TILE_SIZE

			img := g.tileSets[layerIndex].Img(id)

			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + TILE_SIZE))

			opts.GeoM.Translate(g.camera.X, g.camera.Y)

			screen.DrawImage(img, &opts)

			// reset the opts for the next tile
			opts.GeoM.Reset()
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("HP: %d ", g.player.Health))

	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.camera.X, g.camera.Y)
	screen.DrawImage(g.player.Img.SubImage(
		image.Rect(0, 0, 16, 16),
	).(*ebiten.Image), &opts)

	opts.GeoM.Reset()
	for _, enemy := range g.enemies {
		opts.GeoM.Translate(enemy.X, enemy.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)
		screen.DrawImage(enemy.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image), &opts)
		opts.GeoM.Reset()
	}
	opts.GeoM.Reset()
	for _, potion := range g.potions {
		opts.GeoM.Translate(potion.X, potion.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)
		screen.DrawImage(potion.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image), &opts)
		opts.GeoM.Reset()
	}
	opts.GeoM.Reset()

	for _, collider := range g.colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)+float32(g.camera.X), //for making rectangle not moving with camera
			float32(collider.Min.Y)+float32(g.camera.Y),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			true,
		)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//return ebiten.WindowSize()
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	// ebiten.SetTPS(ebiten.SyncWithFPS);
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/NinjaSpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/images/SkeletonSpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	potionImg, _, err := ebitenutil.NewImageFromFile("assets/images/HealingPotion.png")
	if err != nil {
		log.Fatal(err)
	}
	tileMapJSON, err := NewTileMapJSON("assets/maps/spawn-map.json")
	if err != nil {
		log.Fatal(err)
	}
	tileSets, err := tileMapJSON.GenerateTileSets()
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   17,
				Y:   15,
			},
			Health: 100,
		},
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   50,
					Y:   55,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   170,
					Y:   180,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   100,
					Y:   155,
				},
				FollowsPlayer: false,
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   120,
					Y:   128,
				},
				HealAmount: 10,
			},

			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   190,
					Y:   128,
				},
				HealAmount: 10,
			},
		},
		tileMapJSON: tileMapJSON,
		tileSets:    tileSets,
		camera:      NewCamera(0.0, 0.0),
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
