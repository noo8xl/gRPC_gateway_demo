package utils

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/noo8xl/anvil-common/exceptions"
)

type LoggerDto struct {
	Ip        string
	Method    string
	Url       string
	UserAgent string
	Referer   string
	Host      string
}

// RequestLoggerUtil logs the request
func RequestLoggerUtil(dto LoggerDto) {

	var file *os.File
	var err error
	fileName := "err.log"

	curTime := time.Now().Format("2006-01-02")
	dirName := strings.Join([]string{"LOG", curTime}, "")

	// Create logs directory if it doesn't exist
	logsDir := "../LOGS"
	if _, err = os.Stat(logsDir); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(logsDir, 0700)
			if err != nil {
				log.Printf("cannot create logs directory: %v", err)
				exceptions.HandleAnException(err)
			}
		}
	}

	// Create date-specific directory inside logs
	fullPath := logsDir + "/" + dirName
	if _, err = os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(fullPath, 0700)
			if err != nil {
				log.Printf("cannot create error log directory: %v", err)
				exceptions.HandleAnException(err)
			}
		}
	}

	// Create or verify log file exists
	filePath := fullPath + "/" + fileName
	if _, err = os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filePath)
			if err != nil {
				log.Printf("cannot create error log file: %v", err)
				exceptions.HandleAnException(err)
			}
		}
	}

	file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Cannot open a file.")
		exceptions.HandleAnException(err)
	}

	// close a file and terminate the process
	defer func() {
		file.Close()
	}()

	str := curTime + " :> " + dto.Method + " " + dto.Url + " " + dto.Ip + " " + dto.UserAgent + " " + dto.Referer + " " + dto.Host + " \n"
	_, err = file.WriteString(str)
	if err != nil {
		log.Printf("Cannot write to the file.")
		exceptions.HandleAnException(err)
	}

	return
}
