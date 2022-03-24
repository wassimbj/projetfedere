package utils

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logsFilePath string = filepath.Join(RootPath(), "logs/logsfile.txt")
var Logger zerolog.Logger

func init() {
	logsFile, logsFileErr := os.OpenFile(logsFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if logsFileErr != nil {
		log.Fatal().Err(logsFileErr)
	} else {
		log.Logger = log.Output(logsFile)
	}

	// causes an error, just let the GC handle it
	// defer logsFile.Close()

	logsFileStat, _ := logsFile.Stat()
	// log.Printf("LOGS FILE SIZE: %dKb", logsFileStat.Size()/1024)

	// max logs file size is 1MB
	if logsFileStat.Size() >= (1024 * 1024) {
		os.Remove(logsFilePath)
	}

	Logger = log.Logger
}
