package main

import (
	controller "blockchain/main/controller"
)

/**
* controller (HTTP CALL) -> Service -> lib -> model
Service -> lib - можно обьединить
*
*/

func main() {
	controller.HandleConnection()

	/*go func() {
		for candidate := range candidateBlocks {
			tempBlocks = append(tempBlocks, candidate)
		}
	}()

	go func() {
		for {
			Blockchain = lib.PickWinner(Blockchain, tempBlocks, candidateBlocks, validators)
		}
	}()*/

}
