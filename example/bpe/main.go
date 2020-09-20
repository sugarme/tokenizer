package main

import (
	"flag"
)

var mode string

func init() {
	flag.StringVar(&mode, "mode", "test", "run train or test mode")
}

func main() {
	flag.Parse()
	switch mode {
	case "test":
		runTest()
	case "train":
		runTrain()
	}

}
