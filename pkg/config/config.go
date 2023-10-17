package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Window        WindowConfig
	UI            UIConfig
	Miscellaneous MiscellaneousConfig
}

type WindowConfig struct {
	Name      string
	Width     int
	Height    int
	FixedSize bool
}

type UIConfig struct {
	PointColorRGBA               [4]uint8
	SegmentColorRGBA             [4]uint8
	OffsetSegmentColorRGBA       [4]uint8
	SelectedSegmentColorRGBA     [4]uint8
	BackgroundColorRGBA          [4]uint8
	SecondaryBackgroundColorRGBA [4]uint8
	RasterBorderColorRGBA        [4]uint8
	ConstraintCharColorRGBA      [4]uint8
	RasterWidth                  int
	RasterHeight                 int
	PointRadius                  float64
	FillPoints                   bool
}

type MiscellaneousConfig struct {
	LineCatchError         float64
	MoveOverlapPointLength float64
	AllowMoveOverlapPoint  bool
	MaxSliderValue         float64
	OffsetAlgorithm        string
}

func Load(r io.Reader) (*Config, error) {
	var data Config
	_, err := toml.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func LoadStandard(dir string, filename string) (*Config, error) {
	path := filepath.Join(dir, filename)
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer r.Close()
	return Load(r)
}
