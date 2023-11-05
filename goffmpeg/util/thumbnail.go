package util

import (
	"fmt"
	"log"
	"os/exec"
)

type ExtractThumbnailParam struct {
	Input       string
	Duration    string
	ReduceScale int
	Destination string
}

func ExtractThumbnail(resolution *Resolution, param *ExtractThumbnailParam) {
	reduceW := resolution.Width / param.ReduceScale
	reduceH := resolution.Height / param.ReduceScale
	scaleOpt := fmt.Sprintf("thumbnail,scale=%d:%d", reduceW, reduceH)
	output, err := exec.Command(
		Ffmpeg,
		"-ss",
		param.Duration,
		"-i",
		param.Input,
		"-vf",
		scaleOpt,
		"-vframes",
		"1",
		param.Destination).Output()
	if err != nil {
		log.Fatal("extract thumbnail failed:", err)
	}
	log.Println(string(output))
}

func NewThumbKey(folderName string) string {
	return fmt.Sprintf("%s/%s/%s.%s", keyPrefix, folderName, tmpEncodeName, tmpThumbFormat)
}

func GetTmpThumbPath() string {
	return fmt.Sprintf("/tmp/%s.%s", tmpThumbName, tmpThumbFormat)
}

func GetThumbContentType() string {
	return fmt.Sprintf("image/%s", tmpThumbFormat)
}
