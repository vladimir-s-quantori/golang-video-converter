package main

import (
	"encoding/json"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"math/rand"
	"net"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type converter interface {
	resizeWithErr(fileName, inputDir, outputDir, videoQuality string) error
}

type ConverterImpl struct {
}

func checkDirectories(paths ...string) error {
	for _, dir := range paths {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

// Try this one too
func resize(input, outputDir, videoQuality string, errorsChanel chan error) {
	size := fmt.Sprintf("%dx%d", videoSizes[videoQuality].X, videoSizes[videoQuality].Y)
	fmt.Print("Video size is: " + size + "\n")
	outputFile := fmt.Sprintf("%sresize_%v_%v.mp4", outputDir, strings.Split(input, ".")[0], size)
	fmt.Printf("Output file is: %v\\%v\n", outputDir, outputFile)
	err := ffmpeg.Input(input).Output(outputFile, ffmpeg.KwArgs{"s": size}).OverWriteOutput().Run()
	if err != nil {
		errorsChanel <- err
	}
}

func (c ConverterImpl) resizeWithErr(fileName, inputDir, outputDir, videoQuality string) error {
	size := fmt.Sprintf("%dx%d", videoSizes[videoQuality].X, videoSizes[videoQuality].Y)
	fmt.Print("Video size is: " + size + "\n")
	inputFile := fmt.Sprintf("%s%s", inputDir, fileName)
	outputFile := fmt.Sprintf("%sresize_%s_%s.mp4", outputDir, fileName, size)
	fmt.Printf("Out file is: %v\n", outputFile)
	//Use channel for progress streaming smh
	probe, err := ffmpeg.Probe(inputFile)
	if err != nil {
		return err
	}
	total, err := probeDuration(probe)
	if err != nil {
		return err
	}

	fs := ffmpeg.Input(inputFile).
		Output(outputFile, ffmpeg.KwArgs{"s": size}).
		GlobalArgs("-progress", "unix://"+TempSock(total))

	return err
}

func TempSock(totalDuration float64) string {
	// serve
	rand.New(rand.NewSource(time.Now().Unix()))
	sockFileName := path.Join(os.TempDir(), fmt.Sprintf("%d_sock", rand.Int()))
	l, err := net.Listen("unix", sockFileName)
	if err != nil {
		panic(err)
	}

	go func() {
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		buf := make([]byte, 16)
		data := ""
		progress := ""
		for {
			_, err := fd.Read(buf)
			if err != nil {
				return
			}
			data += string(buf)
			a := re.FindAllStringSubmatch(data, -1)
			cp := ""
			if len(a) > 0 && len(a[len(a)-1]) > 0 {
				c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
				cp = fmt.Sprintf("%.2f", float64(c)/totalDuration/1000000)
			}
			if strings.Contains(data, "progress=end") {
				cp = "done"
			}
			if cp == "" {
				cp = ".0"
			}
			if cp != progress {
				progress = cp
				fmt.Println("progress: ", progress)
			}
		}
	}()

	return sockFileName
}

type probeFormat struct {
	Duration string `json:"duration"`
}

type probeData struct {
	Format probeFormat `json:"format"`
}

func probeDuration(input string) (float64, error) {
	pd := probeData{}
	err := json.Unmarshal([]byte(input), &pd)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(pd.Format.Duration, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}
