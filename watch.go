package dockerdiskwatcher

import (
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var sizeLimit = float64(650)

// getDirSize calculates the total size of all files in the directory.
func getDirSize(path string) (float64, error) {
	var size float64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += float64(info.Size())
		}
		return nil
	})
	return size / math.Pow(1024, 3), err
}

// pauseDockerCompose pauses the Docker Compose services.
func pauseDockerCompose(composeFile string) error {
	cmd := exec.Command("docker", "compose", "-f", composeFile, "pause")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Watch() {
	// Parse command-line arguments
	composeFile := flag.String("compose-file", "", "Path to the Docker Compose file")
	directory := flag.String("directory", "", "Path to the directory to watch")
	hardlimit := flag.Float64("hardlimit", 2, "amount to be set as gb hard limit")
	flag.Parse()
	sizeLimit = *hardlimit
	logrus.Infof("this is the current sizelimit : %f GB", sizeLimit)

	if *composeFile == "" || *directory == "" {
		logrus.Info("Both --compose-file and --directory arguments are required")
		flag.Usage()
		os.Exit(1)
	}
	SendTelegramMessage(fmt.Sprintf("watching %s and config file %s for limit %f is starting", *directory, *composeFile, sizeLimit))

	// Watch the directory size
	for {
		size, err := getDirSize(*directory)
		logrus.Infof("this is the current size of that dir: %f GB", size)
		if err != nil {
			logrus.Errorf("Error getting directory size: %v\n", err)
			SendTelegramMessage(fmt.Sprintf("Error getting directory size: %v\n", err))
			os.Exit(1)
		}

		if size > sizeLimit {
			logrus.Info("Directory size limit exceeded, pausing Docker Compose services...")
			if err := pauseDockerCompose(*composeFile); err != nil {
				logrus.Errorf("Error pausing Docker Compose services: %v\n", err)
				SendTelegramMessage(fmt.Sprintf("Error pausing Docker Compose services: %v\n", err))

				os.Exit(1)
			}
			SendTelegramMessage(fmt.Sprintf("watching %s and config file %s for limit %f -- compose has been paused", *directory, *composeFile, sizeLimit))

			logrus.Info("Docker Compose services paused.")
			break
		}

		// Sleep for a while before checking again
		time.Sleep(1 * time.Second)
	}

}
