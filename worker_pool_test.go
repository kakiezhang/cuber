package cuber

import (
	"log"
)

func myJob(message *Msg) {
	log.Println(message)
}

func main() {
	Configure(map[string]string{
		"server":   "localhost:6380",
		"database": "5",
		"pool":     "10",
	})

	Process("myQueue", myJob, 10)

	Run()
}
