package profile

import (
	"log"
	"time"
)

func Duration(start time.Time,desc string)time.Duration {
	elapsed := time.Since(start)
	log.Printf("%s lasted %s", desc, elapsed)
	return elapsed
}