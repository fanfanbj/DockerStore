package main

import (
	"fmt"
	"registry/client"
)

func main () {
	result,_, err := registry.RegistryAPI("GET","/v2/_catalog","admin","badmin","")
	if err != nil {
		fmt.Println(err)
		return 
	}
	fmt.Println(string(result))
}
