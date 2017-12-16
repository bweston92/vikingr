// vickginr is a server that listens to new repositories and commits and
// will ensure the master branch is protected.
package main

import (
	"log"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
