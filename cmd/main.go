package main

import (
	"server/boot"

	"log"
	"os"
	"runtime"
	"time"
)

func main() {
	setupEnvironment()

	if err := boot.Boot(); err != nil {
		log.Fatalln("start failure: ", err.Error())
	}
}

func setupEnvironment() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	timeZone := os.Getenv("TIMEZONE")
	env := os.Getenv("ENV")

	location, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Fatal(err)
	}

	time.Local = location

	log.Printf("------ ENV:%s TIMEZONE:%s CurrentTime:%s\n", env, timeZone, time.Now().Format(time.RFC3339))
}
