package main

import (
	"log"

	appinit "sample_app/init"
)

func main() {
	if err := appinit.App().ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
