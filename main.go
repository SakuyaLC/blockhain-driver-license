package main

import (
	controller "blockchain/main/controller"
	"blockchain/main/internal/model"
	"sync"
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

var mutex = &sync.Mutex{}

// validators keeps track of open validators and balances
var validators = make(map[string]int)

func main() {
	controller.HandleConnection()
	/*err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	// create genesis block
	t := time.Now()
	genesisBlock := model.Block{}
	genesisBlock = model.Block{Index: 0, Timestamp: t.String(), LicenseInfo: "genesisLicenseInfo", Hash: lib.CalculateBlockHash(genesisBlock), PrevHash: "", Validator: ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	tcpPort := os.Getenv("PORT")

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("TCP Server Listening on port :", tcpPort)
	defer server.Close()

	go func() {
		for candidate := range candidateBlocks {
			mutex.Lock()
			tempBlocks = append(tempBlocks, candidate)
			mutex.Unlock()
		}
	}()

	go func() {
		for {
			Blockchain = lib.PickWinner(Blockchain, tempBlocks, candidateBlocks, validators)
			lib.ShowBlockchain(Blockchain)
		}
	}()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go controller.HandleConn(conn, Blockchain, validators, candidateBlocks)
	}*/

}
