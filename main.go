package main

import (
	controller "blockchain/main/controller"
	lib "blockchain/main/internal"
	"blockchain/main/internal/model"
	"time"

	"github.com/davecgh/go-spew/spew"
)

/**
* controller (HTTP CALL) -> Service -> lib -> model
Service -> lib - можно обьединить
*
*/
// Blockchain is a series of validated Blocks
var Blockchain []model.Block
var tempBlocks []model.Block

// candidateBlocks handles incoming blocks for validation
var candidateBlocks = make(chan model.Block)

// validators keeps track of open validators and balances
var validators = make(map[string]int)

func main() {
	controller.HandleConnection()

	// create genesis block
	t := time.Now()
	genesisBlock := model.Block{}
	genesisBlock = model.Block{Index: 0, Timestamp: t.String(), LicenseInfo: "genesisLicenseInfo", Hash: lib.CalculateBlockHash(genesisBlock), PrevHash: "", Validator: ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

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
