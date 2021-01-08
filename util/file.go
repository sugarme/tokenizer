package util

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// CdToThis changes `working directory` to the current directory
// of function that calls it.
//
// A use case is when a function at coding file reads a relative file,
// however the end user, at run time, executes the coding file
// at a different directory to the coding file. As relative path has been
// used, one would change from any directory at run time to `current` directory
// where function is called to get the correct relative path.
//
// Example:
// There a Go package in a os directory as follow:
// + os-home
// 		|
//		+--go-package
//		|		|
//		|		+--data
//		|		|		|
//		|		|		+--data-file.txt
//		|		|
//		|		+--sub-package
//		|				|
//		|				+--go-file.go
//		|
//		+--other-folder
//
// There's a function `ReadData` in `go-file.go` that
// read a data file `data-file.txt` in other directory using
// a relative filepath to its - I.e: "../data/data-file.txt".
// Now, a end user run `go-file.go` from `other-folder` directory
// which is not at the same directory to `go-file.go` file. `ReadData`
// function would be panic if there was no helper to change directory
// from directory where end user executes the code ("os-home/other-folder")
// to directory of `go-file.go`. `CdToThis` does this job.
func CdToThis() {
	// pc, file, line, ok := runtime.Caller(1)
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatal("Runtime Caller error: cannot get information of current caller.")
	}
	currDir := filepath.Dir(file)

	if err := os.Chdir(currDir); err != nil {
		log.Fatal(err)
	}
}

// CdBack returns back to previous path.
func CdBack(backDir string) {
	if err := os.Chdir(backDir); err != nil {
		log.Fatal(err)
	}
}
