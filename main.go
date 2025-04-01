package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	X, Y float64
	Img  *ebiten.Image
}

type Player struct {
	*Sprite
	Health uint16
}

type Enemy struct {
	*Sprite
	FollowsPlayer bool
}

type Potion struct {
	*Sprite
	HealAmount uint16
}

type Game struct {
	player  *Player
	enemies []*Enemy
	potions []*Potion
	// connectedControllers holds a slice of GamepadIDs representing the
	// currently connected game controllers. Each GamepadID corresponds
	// to a unique identifier for a connected gamepad device.
	connectedControllers []ebiten.GamepadID
	GamepadID            ebiten.GamepadID
}

const (
	KB_SPEED  = 2
	GP_SPEED  = 2.5
	DEAD_ZONE = 0.1
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
	

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	screen.Fill(color.RGBA{120, 180, 255, 255})
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %.2f Y: %.2f", g.player.X, g.player.Y))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("HP: %d ", g.player.Health))
	//draw the player
	// screen.DrawImage(g.PlayerImage, &ebiten.DrawImageOptions{})
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.player.X, g.player.Y)
	screen.DrawImage(g.player.Img.SubImage(
		image.Rect(0, 0, 16, 16),
	).(*ebiten.Image), &opts)

	opts.GeoM.Reset()
	for _, enemy := range g.enemies {
		opts.GeoM.Translate(enemy.X, enemy.Y)
		screen.DrawImage(enemy.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image), &opts)
		opts.GeoM.Reset()
	}
	opts.GeoM.Reset()
	for _, potion := range g.potions {
		opts.GeoM.Translate(potion.X, potion.Y)
		screen.DrawImage(potion.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image), &opts)
		opts.GeoM.Reset()
	}
	opts.GeoM.Reset()

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// return ebiten.WindowSize()
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

	if err := ebiten.RunGame(&Game{
		player: &Player{
			&Sprite{
				Img: playerImg,
				X:   100,
				Y:   100,
			},
			100,
		},
		enemies: []*Enemy{
			{
				&Sprite{
					Img: skeletonImg,
					X:   50,
					Y:   55,
				},
				true,
			},
			{
				&Sprite{
					Img: skeletonImg,
					X:   170,
					Y:   180,
				},
				true,
			},
			{
				&Sprite{
					Img: skeletonImg,
					X:   100,
					Y:   155,
				},
				false,
			},
		},
		potions: []*Potion{
			{
				&Sprite{
					Img: potionImg,
					X: 120,
					Y: 128,
				},
				10,
			},

			{
				&Sprite{
					Img: potionImg,
					X: 190,
					Y: 128,
				},
				10,
			},
		},
	}); err != nil {
		log.Fatal(err)
	}
}
