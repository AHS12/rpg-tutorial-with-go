package main

import (
	"encoding/json"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TileSet interface {
	Img(id int) *ebiten.Image
}

type UniformTileSet struct {
	img *ebiten.Image
	gid int
}

type UniformTileSetJSON struct {
	Path string `json:"image"`
}

func (u *UniformTileSet) Img(id int) *ebiten.Image {
	id -= u.gid
	//get the position of the image where the tile id is
	srcX := id % 22
	srcY := id / 22

	//convert the src tile position to pixel src position
	srcX *= TILE_SIZE
	srcY *= TILE_SIZE

	return u.img.SubImage(
		image.Rect(
			srcX,
			srcY,
			srcX+TILE_SIZE,
			srcY+TILE_SIZE,
		),
	).(*ebiten.Image)
}

type DynamicTileSet struct {
	images []*ebiten.Image
	gid    int
}

type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Height int    `json:"imageheight"`
	Width  int    `json:"imagewidth"`
}

type DynamicTileSetJSON struct {
	Tiles []TileJSON `json:"tiles"`
}

func (d DynamicTileSet) Img(id int) *ebiten.Image {
	id -= d.gid

	return d.images[id]
}

func getPathString(path string) string {
	cleanedPath := filepath.Clean(path)
	cleanedPath = strings.ReplaceAll(cleanedPath, "\\", "/")
	cleanedPath = strings.TrimPrefix(cleanedPath, "../")
	cleanedPath = strings.TrimPrefix(cleanedPath, "../")
	cleanedPath = filepath.Join("assets/", cleanedPath)

	return cleanedPath
}

func NewTileSet(path string, gid int) (TileSet, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if strings.Contains(path, "buildings") {
		//return dynamic tileset
		var dynTileSetJSON DynamicTileSetJSON
		err = json.Unmarshal(contents, &dynTileSetJSON)
		if err != nil {
			return nil, err
		}

		dyneTileSet := DynamicTileSet{}
		dyneTileSet.gid = gid
		dyneTileSet.images = make([]*ebiten.Image, 0)

		for _, tileJSON := range dynTileSetJSON.Tiles {
			img, _, err := ebitenutil.NewImageFromFile(getPathString(tileJSON.Path))
			if err != nil {
				return nil, err
			}

			dyneTileSet.images = append(dyneTileSet.images, img)
		}
		return dyneTileSet, nil
	}
	//return Uniform tileset
	var uniformTileSetJSON UniformTileSetJSON
	err = json.Unmarshal(contents, &uniformTileSetJSON)
	if err != nil {
		return nil, err
	}

	uniformTileSet := UniformTileSet{}
	img, _, err := ebitenutil.NewImageFromFile(getPathString(uniformTileSetJSON.Path))
	if err != nil {
		return nil, err
	}
	uniformTileSet.img = img
	uniformTileSet.gid = gid

	return &uniformTileSet, nil

}
