package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
)

// TileMapLayerJSON represents a single tile layer in the map.
type TileMapLayerJSON struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Data   []int  `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type TileSetJSON struct {
	FirstGID int    `json:"firstgid"`
	Source   string `json:"source"`
}

// TileMapJSON represents the entire tile map.
type TileMapJSON struct {
	Layers     []TileMapLayerJSON `json:"layers"`
	TileSets   []TileSetJSON      `json:"tilesets"`
	TileWidth  int                `json:"tilewidth"`
	TileHeight int                `json:"tileheight"`
}

func (t *TileMapJSON) GenerateTileSets() ([]TileSet, error) {
	tileSetList := make([]TileSet, 0)
	//generate a tile map based on
	for _, tileSetData := range t.TileSets {
		tileSetPath := path.Join("assets/maps/", tileSetData.Source)
		tileset, err := NewTileSet(tileSetPath, tileSetData.FirstGID)
		if err != nil {
			return nil, err
		}

		tileSetList = append(tileSetList, tileset)

	}

	return tileSetList, nil
}

// NewTileMapJSON loads and parses a TileMapJSON from a file.
func NewTileMapJSON(path string) (*TileMapJSON, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var tileMap TileMapJSON
	if err := json.Unmarshal(contents, &tileMap); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	sort.Slice(tileMap.Layers, func(i, j int) bool {
		// Sort in ascending order of ID
		return tileMap.Layers[i].Id < tileMap.Layers[j].Id
	})

	return &tileMap, nil
}
