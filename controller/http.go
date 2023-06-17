package controller

import (
	"blockchain/main/db"
	lib "blockchain/main/internal"
	"blockchain/main/internal/model"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
)

// Blockchain is a series of validated Blocks
var Blockchain []model.Block

// candidateBlocks handles incoming blocks for validation
var candidateBlocks []model.Block

// validators keeps track of open validators and balances
var validators = make(map[string]int)

func HandleConnection() {
	app := fiber.New()

	Blockchain = db.UpdateBlockchain()

	hasRecords, err := db.CheckIfCollectionHasRecords(context.Background(), "blocks")
	if err != nil {
		// Обработка ошибки
		fmt.Println("Error while initiating blockchain!")
	}
	if hasRecords {
		// Если есть записи в коллекции
		Blockchain = db.UpdateBlockchain()
	} else {
		// Если нет записей в коллекции
		t := time.Now()
		genesisBlock := model.Block{}
		genesisBlock = model.Block{Index: 0, Timestamp: t.String(), Info: "genesisLicenseInfo", Hash: lib.CalculateBlockHash(genesisBlock), PrevHash: "", Validator: ""}
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)
		db.InsertBlock(genesisBlock)
	}

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Success")
	})
	app.Get("/lottery", StartLottery)
	app.Get("/blockchain", GetBlockchain)

	app.Post("/create-block", CreateBlock)

	app.Listen(":8080")
}

func CreateBlock(ctx *fiber.Ctx) error {

	startTime := time.Now()

	var message struct {
		Message string `json:"message"`
	}
	err := ctx.BodyParser(&message)

	if err != nil {
		ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		return err
	}

	// validator address
	var address string = lib.CalculateHash(time.Now().String())

	info := message.Message

	// Случайное число в диапазоне от 50 до 100
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(51) + 50

	validators[address] = randomInt

	oldLastIndex := Blockchain[len(Blockchain)-1]

	// create newBlock for consideration to be forged
	newBlock, err := lib.GenerateBlock(oldLastIndex, info, address)
	if lib.IsBlockValid(newBlock, oldLastIndex) {
		candidateBlocks = append(candidateBlocks, newBlock)
	}

	//fmt.Println("New block info: " + newBlock.Info)
	fmt.Println(newBlock.Index)

	endTime := time.Now()

	executionTime := endTime.Sub(startTime)

	fmt.Println("Execution time:", executionTime.Microseconds())

	return ctx.Status(fiber.StatusOK).JSON(candidateBlocks)
}

func StartLottery(ctx *fiber.Ctx) error {
	Blockchain = lib.PickWinner(Blockchain, candidateBlocks, validators)
	candidateBlocks = nil
	return ctx.Status(fiber.StatusOK).JSON("Lottery started")
}

func GetBlockchain(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(Blockchain)
}
