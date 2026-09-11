package main

import (
	"fmt"
	"property/src/config"
	"property/src/read"
)

func main() {
	println("Hello, World!")
	
	config.InitConfig()
	atts := read.ReadAtt()

	fmt.Printf(" %v ", atts)
}
