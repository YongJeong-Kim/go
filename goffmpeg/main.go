package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"goffmpeg/util"
	"log"
	"os/exec"
	"strconv"
)

func handler(ctx context.Context, s3Event events.S3Event) error {
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load default config: %s", err)
		return err
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	var b string
	var k string
	var region string
	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		b = bucket
		key := record.S3.Object.URLDecodedKey
		k = key
		region = record.AWSRegion
		getOutput, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err != nil {
			log.Printf("error getting of object %s/%s: %s", bucket, key, err)
			return err
		}
		log.Printf("successfully retrieved %s/%s of type %s", bucket, key, *getOutput.ContentType)

		path := util.GetS3URL(region, b, k)
		resol := util.GetResolution(path)
		resolution, err := util.Resolution2WidthHeight(resol)
		if err != nil {
			log.Fatal(err)
		}

		tmpEncodePath := util.GetTmpEncodePath()
		util.ExecEncode(resolution, &util.ExecEncodeParam{
			Input:       path,
			ReduceScale: util.ReduceScale,
			Destination: tmpEncodePath,
		})

		duration := util.GetRandomDuration(tmpEncodePath)

		tmpThumbPath := util.GetTmpThumbPath()
		util.ExtractThumbnail(resolution, &util.ExtractThumbnailParam{
			Input:       tmpEncodePath,
			Duration:    strconv.Itoa(duration),
			ReduceScale: util.ReduceScale,
			Destination: tmpThumbPath,
		})

		folderName := uuid.NewString()
		newThumbKey := util.NewThumbKey(folderName)
		fileThumb := util.GetFile(tmpThumbPath)
		thumbContentType := util.GetThumbContentType()
		util.Upload(ctx, s3Client, &s3.PutObjectInput{
			Bucket:      aws.String(b),
			Key:         aws.String(newThumbKey),
			Body:        fileThumb,
			ACL:         types.ObjectCannedACLPublicRead,
			ContentType: aws.String(thumbContentType),
		})

		newEncodeKey := util.NewEncodeKey(folderName)
		fileEncoded := util.GetFile(tmpEncodePath)
		encodeContentType := util.GetEncodeContentType()
		util.Upload(ctx, s3Client, &s3.PutObjectInput{
			Bucket:      aws.String(b),
			Key:         aws.String(newEncodeKey),
			Body:        fileEncoded,
			ACL:         types.ObjectCannedACLPublicRead,
			ContentType: aws.String(encodeContentType),
		})

		err = exec.Command("rm", "-f", "/tmp/output.mp4").Run()
		if err != nil {
			log.Fatal("rm mp4 failed: ", err)
		}
		err = exec.Command("rm", "-f", "/tmp/output.jpg").Run()
		if err != nil {
			log.Fatal("rm jpg failed: ", err)
		}
	}

	//var c *exec.Cmd

	/*output, err := exec.Command("/opt/bin/ffmpeg", "-version").Output()
	//c.Stdout = os.Stdout
	//output, err := c.Output()
	if err != nil {
		log.Printf("ffmpeg version failed: %v\n", err)
	}

	log.Println(string(output))*/

	/*log.Println("path: ", path)
	o, err := exec.Command("/opt/bin/ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", path).Output()
	if err != nil {
		log.Printf("ffprobe failed: %v", err)
	}
	resolution := string(o)
	log.Println(string(o))

	wh := strings.Split(resolution, "x")
	width := wh[0]
	height := wh[1]

	log.Println("width:", width)
	log.Println("height:", height)*/

	/*d, err := exec.Command("/opt/bin/ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=nw=1:nk=1", path).Output()
	if err != nil {
		log.Println("ffprobe duration failed: ", err)
	}
	durationStr := string(d)
	log.Println("duration:", durationStr)

	duration, err := strconv.ParseFloat(strings.TrimSpace(durationStr), 64)
	if err != nil {
		log.Println("duration failed: ", err)
	}
	log.Println("duration: ", duration)

	durationInt := int(duration)
	randomDur := rand.Intn(durationInt)
	log.Println("random duration: ", randomDur)*/

	/*err = exec.Command(
		"/opt/bin/ffmpeg",
		"-i",
		path,
		"-f",
		"mp4",
		"-movflags",
		"frag_keyframe+empty_moov",
		"/tmp/output.mp4",
	).Run()
	if err != nil {
		log.Println("encode failed:", err)
	}
	log.Println("vanilla encode completed")*/

	/*to, err := exec.Command(
		"/opt/bin/ffmpeg",
		"-ss",
		strconv.Itoa(randomDur),
		"-i",
		path,
		"-vf",
		"thumbnail,scale=160:90",
		"-vframes",
		"1",
		"/tmp/screenshot.jpg").Output()
	if err != nil {
		log.Println("thumbnail failed:", err)
	}
	log.Println(string(to))*/

	/*file2, err := os.Open("/tmp/output.mp4")
	if err != nil {
		log.Println("open file mp4 failed:", err)
	}
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("ppppppnnnnnn"),
		Key:         aws.String("output.mp4"),
		Body:        file2,
		ACL:         types.ObjectCannedACLPublicRead,
		ContentType: aws.String("video/mp4"),
	})
	if err != nil {
		log.Println("put object failed:", err)
	}
	log.Println("thumb upload success")*/

	/*file, err := os.Open("/tmp/screenshot.jpg")
	if err != nil {
		log.Println("open file failed:", err)
	}

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("ppppppnnnnnn"),
		Key:         aws.String("screenshot1.jpg"),
		Body:        file,
		ACL:         types.ObjectCannedACLPublicRead,
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		log.Println("put object failed:", err)
	}
	log.Println("thumb upload success")*/

	return nil
}

func main() {
	lambda.Start(handler)
}

/*func main() {
	var c *exec.Cmd

		switch runtime.GOOS {
		case "windows":
		default:
		}

	//c = exec.Command("ffmpeg-6.0-essentials_build\\ffmpeg-6.0-essentials_build\\bin\\ffmpeg.exe", "-version")
	c = exec.Command("ffmpeg", "-version")
	//path := "ffmpeg-6.0-essentials_build\\ffmpeg-6.0-essentials_build\\bin\\ffprobe.exe"
	//c = exec.Command(path, "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", "example_480_1_5MG.mp4")
	//err := c.Run()
	output, err := c.Output()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(output))
}*/

/*
extract thumbnail(-ss 뒤에는 초)
.\ffmpeg.exe -ss 16 -i C:\Users\yongyong\Downloads\file_example_MP4_480_1_5MG.mp4 -vf thumbnail,scale=160:90 -vframes 1 mm.jpg


video duration(seconds, float)
.\ffprobe.exe -v error -show_entries format=duration -of default=nw=1:nk=1 C:\Users\yongyong\Downloads\file_example_MP4_480_1_5MG.mp4

video resolution(e.g. 480x270)
.\ffprobe.exe -v error -select_streams v:0 -show_entries stream=width,height -of csv=s=x:p=0 C:\Users\yongyong\Downloads\file_example_MP4_480_1_5MG.mp4

video resolution(e.g. json)
.\ffprobe.exe -v error -select_streams v:0 -show_entries stream=width,height -of json C:\Users\yongyong\Downloads\file_example_MP4_480_1_5MG.mp4

*/

/*

encode to s3
maybe add s3 CLI role
https://stackoverflow.com/questions/54736246/stream-ffmpeg-transcoding-result-to-s3

ffmpeg -i https://notreal-bucket.s3-us-west-1.amazonaws.com/video/video.mp4 -f mp4 -movflags frag_keyframe+empty_moov pipe:1 | aws s3 cp - s3://notreal-bucket/video/output.mp4
aws s3 cp s3://source-bucket/source.mp4 - | ffmpeg -i - -f matroska - | aws s3 cp - s3://dest-bucket/output.mkv
*/
