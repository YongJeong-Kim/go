package util

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

type ExecEncodeParam struct {
	Input       string
	ReduceScale int
	Destination string
}

func ExecEncode(resolution *Resolution, param *ExecEncodeParam) {
	reduceW := resolution.Width / param.ReduceScale
	reduceH := resolution.Height / param.ReduceScale
	scaleOpt := fmt.Sprintf("scale=%s:%s", strconv.Itoa(reduceW), strconv.Itoa(reduceH))
	log.Println(param.Input)
	log.Println(scaleOpt)
	log.Println(param.Destination)
	err := exec.Command(
		Ffmpeg,
		"-i",
		param.Input,
		//"-vf",
		//"scale=240:135",
		"-f",
		tmpEncodeFormat,
		"-movflags",
		"frag_keyframe+empty_moov",
		param.Destination,
	).Run()
	if err != nil {
		log.Fatal("encode failed: ", err)
	}
	log.Println("encode completed")
}

func GetTmpEncodePath() string {
	return fmt.Sprintf("/tmp/%s.%s", tmpEncodeName, tmpEncodeFormat)
}

func NewEncodeKey(folderName string) string {
	return fmt.Sprintf("%s/%s/%s.%s", keyPrefix, folderName, tmpEncodeName, tmpEncodeFormat)
}

func GetEncodeContentType() string {
	return fmt.Sprintf("video/%s", tmpEncodeFormat)
}
