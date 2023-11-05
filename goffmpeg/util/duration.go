package util

import (
	"log"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
)

func getDuration(key string) string {
	d, err := exec.Command(
		Ffprobe,
		"-v",
		"error",
		"-show_entries",
		"format=duration",
		"-of",
		"default=nw=1:nk=1",
		key,
	).Output()
	if err != nil {
		log.Fatal("ffprobe duration failed: ", err)
	}

	durationStr := string(d)
	log.Println("duration str:", durationStr)

	return durationStr
}

func duration2SecondInt(durationStr string) int {
	duration, err := strconv.ParseFloat(strings.TrimSpace(durationStr), 64)
	if err != nil {
		log.Fatal("parse float duration failed: ", err)
	}
	log.Println("duration float:", duration)

	return int(duration)
}

func GetRandomDuration(tmpEncodePath string) int {
	duration := getDuration(tmpEncodePath)
	d := duration2SecondInt(duration)

	randomDur := rand.Intn(d)
	log.Println("random duration:", randomDur)

	return randomDur
}
