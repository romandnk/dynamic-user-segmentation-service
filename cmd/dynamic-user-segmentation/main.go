package main

import (
	"flag"
	"fmt"
)

var configFile string

func main() {
	flag.StringVar(&configFile, "config", "./configs/dynamic-user-segmentation.toml", "Path to configuration file")
	fmt.Println("ready")
}
