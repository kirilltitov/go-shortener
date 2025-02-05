package main

import "os"

// должно зафейлиться, потому что main.main
func main() {
	go os.Exit(255)    // want "os.Exit in main.main function is forbidden"
	defer os.Exit(255) // want "os.Exit in main.main function is forbidden"
	os.Exit(255)       // want "os.Exit in main.main function is forbidden"
}

// не должно зафейлиться, не main.main
func thisIsFine() {
	os.Exit(322)
}
