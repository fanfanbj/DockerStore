package main

import (
	"fmt"
	"registry/client"
)

func main () {
	result,_, err := client.RegistryAPI("GET","/v2/_catalog","")
	if err != nil {
		fmt.Println(err)
		return 
	}
	fmt.Println(string(result))
}
