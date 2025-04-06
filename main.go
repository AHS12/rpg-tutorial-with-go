package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"rpg-tutorial/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	player               *entities.Player
	enemies              []*entities.Enemy
	potions              []*entities.Potion
	tileMapJSON          *TileMapJSON
	tileMapFloorImg      *ebiten.Image
	tileMapObjectImg     *ebiten.Image
	connectedControllers []ebiten.GamepadID
	GamepadID            ebiten.GamepadID
	camera               *Camera
}

const (
	KB_SPEED  = 2
	GP_SPEED  = 2.5
	DEAD_ZONE = 0.1
	TILE_SIZE = 16
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
		g.player.X += GP_SPEED
	}
	if ebiten.IsStandardGamepadButtonPressed(g.GamepadID, GamepadButtonLeft) {
		g.player.X -= GP_SPEED
	}
	if ebiten.IsStandardGamepadButtonPressed(g.GamepadID, GamepadButtonUp) {
		g.player.Y -= GP_SPEED
	}
	if ebiten.IsStandardGamepadButtonPressed(g.GamepadID, GamepadButtonDown) {
		g.player.Y += GP_SPEED
	}
}

func (g *Game) handleNonStandardGamepadButtons() {
	if ebiten.IsGamepadButtonPressed(g.GamepadID, FallbackButtonRight) { // Right
		g.player.X += GP_SPEED
	}
	if ebiten.IsGamepadButtonPressed(g.GamepadID, FallbackButtonLeft) { // Left
		g.player.X -= GP_SPEED
	}
	if ebiten.IsGamepadButtonPressed(g.GamepadID, FallbackButtonUp) { // Up
		g.player.Y -= GP_SPEED
	}
	if ebiten.IsGamepadButtonPressed(g.GamepadID, FallbackButtonDown) { // Down
		g.player.Y += GP_SPEED
	}
}

func isOutsideDeadZone(value float64, deadZone float64) bool {
	return value > deadZone || value < -deadZone
}

func (g *Game) handleAnalogStickMovement() {
	xAxis := ebiten.GamepadAxisValue(g.GamepadID, 0)
	yAxis := ebiten.GamepadAxisValue(g.GamepadID, 1)

	if isOutsideDeadZone(xAxis, DEAD_ZONE) {
		g.player.X += xAxis * GP_SPEED
	}
	if isOutsideDeadZone(yAxis, DEAD_ZONE) {
		g.player.Y += yAxis * GP_SPEED
	}
}

func (g *Game) handleKeyboardControls() {
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.X += KB_SPEED
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.X -= KB_SPEED
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Y -= KB_SPEED
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Y += KB_SPEED
	}
}

func (g *Game) Update() error {

	// Keyboard Controls
	g.handleKeyboardControls()
	// detect controllers
	g.detectAndSelectGamepad()
	// Gamepad Controls
	g.handleGamepadInput()

	//spawning enemy
	for _, enemy := range g.enemies {
		if enemy.FollowsPlayer {
			if enemy.X < g.player.X {
				enemy.X += 0.8
			} else if enemy.X > g.player.X {
				enemy.X -= 0.8
			}
			if enemy.Y < g.player.Y {
				enemy.Y += 0.8
			} else if enemy.Y > g.player.Y {
				enemy.Y -= 0.8
			}
		}
	}

	for _, potion := range g.potions {
		if g.player.X > potion.X {
			g.player.Health += potion.HealAmount
		}
	}

	g.camera.FollowTarget(g.player.X + TILE_SIZE/2, g.player.Y+ TILE_SIZE/2, 320, 240)
	g.camera.Constrain(
		float64(g.tileMapJSON.Layers[0].Width * TILE_SIZE),
		float64(g.tileMapJSON.Layers[0].Height * TILE_SIZE),
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

	for _, layer := range g.tileMapJSON.Layers {
		// fmt.Println("Starting to process Layer:", layer.Name, "ID:", layer.Id)
		var tileset *ebiten.Image
		var tilesetWidth int
		var firstgid int

		if layer.Name == "Floor" {
			tileset = g.tileMapFloorImg
			tilesetWidth = tileset.Bounds().Dx() / TILE_SIZE
			firstgid = g.tileMapJSON.TileSets[0].FirstGID
		} else if layer.Name == "Object" {
			tileset = g.tileMapObjectImg
			tilesetWidth = tileset.Bounds().Dx() / TILE_SIZE
			firstgid = g.tileMapJSON.TileSets[1].FirstGID
		}

		// Loop over the tiles
		for index, tileID := range layer.Data {
			opts.GeoM.Reset()

			// Skip empty tiles
			if tileID == 0 {
				continue
			}

			// Calculate the tile position in the world
			x := (index % layer.Width) * TILE_SIZE
			y := (index / layer.Width) * TILE_SIZE

			// Adjust the tile ID based on the tileset's firstgid
			adjustedID := tileID - firstgid + 1

			// Calculate the position in the tileset image
			srcX := ((adjustedID - 1) % tilesetWidth) * TILE_SIZE
			srcY := ((adjustedID - 1) / tilesetWidth) * TILE_SIZE

			// Safety check to ensure we're not accessing outside the tileset bounds
			tilesetBounds := tileset.Bounds()
			if srcX < 0 || srcY < 0 || srcX+16 > tilesetBounds.Dx() || srcY+16 > tilesetBounds.Dy() {
				fmt.Printf("WARNING: Tile ID %d adjusted to %d gives invalid source rect (%d,%d,%d,%d) for tileset bounds (%d,%d)\n",
					tileID, adjustedID, srcX, srcY, srcX+TILE_SIZE, srcY+TILE_SIZE, tilesetBounds.Dx(), tilesetBounds.Dy())
				continue
			}

			// Draw the tile
			opts.GeoM.Translate(float64(x), float64(y))
			opts.GeoM.Translate(g.camera.X, g.camera.Y)
			screen.DrawImage(
				tileset.SubImage(image.Rect(srcX, srcY, srcX+TILE_SIZE, srcY+TILE_SIZE)).(*ebiten.Image),
				&opts,
			)
		}
		// fmt.Println("Processing Layer Completed:", layer.Name, "ID:", layer.Id)
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
	tileMapFloorImg, _, err := ebitenutil.NewImageFromFile("assets/images/TilesetFloor.png")
	if err != nil {
		log.Fatal(err)
	}
	tileMapObjectImg, _, err := ebitenutil.NewImageFromFile("assets/images/TilesetNature.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(&Game{
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
		tileMapJSON:      tileMapJSON,
		tileMapFloorImg:  tileMapFloorImg,
		tileMapObjectImg: tileMapObjectImg,
		camera: NewCamera(0.0, 0.0),
	}); err != nil {
		log.Fatal(err)
	}
}
