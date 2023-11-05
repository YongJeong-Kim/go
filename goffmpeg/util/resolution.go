package util

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Resolution struct {
	Width  int
	Height int
}

func GetResolution(path string) string {
	resolution, err := exec.Command(
		Ffprobe,
		"-v",
		"error",
		"-select_streams",
		"v:0",
		"-show_entries",
		"stream=width,height",
		"-of",
		"csv=s=x:p=0",
		path,
	).Output()
	if err != nil {
		log.Fatal("get resolution failed.", err)
	}
	log.Println("resolution:", string(resolution))
	return strings.TrimSpace(string(resolution))
}

func Resolution2WidthHeight(resolution string) (*Resolution, error) {
	r := strings.Split(resolution, "x")
	w, err := strconv.Atoi(r[0])
	if err != nil {
		return nil, fmt.Errorf("resolution width split failed. %v", err)
	}

	h, err := strconv.Atoi(r[1])
	if err != nil {
		return nil, fmt.Errorf("resolution height split failed. %v", err)
	}

	log.Printf("convert resolution to width and height: %d, %d\n", w, h)
	return &Resolution{
		Width:  w,
		Height: h,
	}, nil
}

func GetS3URL(region, bucket, key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
}
