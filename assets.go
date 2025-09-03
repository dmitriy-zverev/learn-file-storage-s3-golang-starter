package main

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"os/exec"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		filePath,
	)

	reader := bytes.Buffer{}

	cmd.Stdout = &reader
	if err := cmd.Run(); err != nil {
		return "", err
	}

	type stream struct {
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	}

	type response struct {
		Streams []stream `json:"streams"`
	}

	resp := response{}
	if err := json.NewDecoder(&reader).Decode(&resp); err != nil {
		return "", err
	}

	remainder := resp.Streams[0].Width / resp.Streams[0].Height

	if math.Abs(remainder-1.77777) < 0.001 {
		return "16:9", nil
	}

	if math.Abs(remainder-0.5625) < 0.001 {
		return "9:16", nil
	}

	return "other", nil
}
