package main

import (
	"flag"
	"fmt"
	"os"
)

type videoSize struct {
	X int
	Y int
}

const (
	Standard = "SD"
	Medium   = "MD"
	High     = "HD"
)

var videoSizes = map[string]videoSize{
	Standard: {640, 360},
	Medium:   {854, 480},
	High:     {1280, 720},
}

func init() {}

func main() {
	//TODO: 1. wrap application into grpc service
	//2. Run service in docker container
	//3. Pass message with local file location
	//4. Stream encoding progress from app to subscriber/client
	//5. Limit memory consumption from Docker using env variables
	//6. Test with different file sizes
	inputFlag := flag.String("in", "", "folder with the files to process")
	outputFlag := flag.String("out", "", "folder for processed files")
	flag.Parse()

	//inputPath := "D:/Study/go/vidConv/files/input/"
	//outputPath := "D:/Study/go/vidConv/files/output/"
	inputPath := *inputFlag
	outputPath := *outputFlag

	if err := checkDirectories(inputPath, outputPath); err != nil {
		panic(err)
	}

	files, err := os.ReadDir(inputPath)
	if err != nil {
		panic(err)
	}
	runsCount := 0

	conv := ConverterImpl{}
	errChan := make(chan error)
	defer close(errChan)

	for sizeKey := range videoSizes {
		for _, file := range files {
			if !file.IsDir() {
				runsCount++
				go func(s string) {
					errChan <- conv.resizeWithErr(file.Name(), inputPath, outputPath, s)
				}(sizeKey)
			}
		}
	}

	// Selector approach
	for {
		select {
		case err = <-errChan:
			if err != nil {
				fmt.Print(err)
			}
		}
		return
	}
	fmt.Printf("All %v videos are converted successfully", runsCount)
}
