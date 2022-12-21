package main

import (
	"fmt"
	"goblockchain/utils"
)

func main() {
	ip := utils.GetHost()
	fmt.Println(ip)
	//fmt.Println(utils.FindMyNeighbors("127.0.0.1", 5001, 0, 3, 5001, 5005))
}
