package logger

import (
	"log"
	"os"
)

func LogfileInit() {
	file, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("unable to set logfile: %v\n", err)
	}
	defer file.Close()

	log.SetOutput(file)
}
