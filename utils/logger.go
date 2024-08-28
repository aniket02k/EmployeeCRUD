package utils

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger() {
    Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Println("Logger Initialized")
}
