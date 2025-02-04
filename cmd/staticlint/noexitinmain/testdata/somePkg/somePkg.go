package somePkg

import "os"

// не должно зафейлиться, не main.main
func someFunc() {
	go os.Exit(255)
	defer os.Exit(255)
	os.Exit(255)
}

// не должно зафейлиться, не main.main
func thisIsFine() {
	os.Exit(322)
}
