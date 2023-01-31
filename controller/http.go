package controller

import (
	lib "blockchain/main/internal"
	"blockchain/main/internal/model"
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

	// create genesis block
	t := time.Now()
	genesisBlock := model.Block{}
	genesisBlock = model.Block{Index: 0, Timestamp: t.String(), LicenseInfo: "genesisLicenseInfo", Hash: lib.CalculateBlockHash(genesisBlock), PrevHash: "", Validator: ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Success")
	})
	app.Get("/account", GetAccountInfo)
	app.Get("/lottery", StartLottery)
	app.Get("/blockchain", GetBlockchain)

	app.Post("/create-account", CreateAccount)

	app.Listen(":80")
}

func GetAccountInfo(ctx *fiber.Ctx) error {
	user := model.Block{
		Index:       1,
		Timestamp:   time.Now().String(),
		LicenseInfo: "License Info",
		Hash:        "Some hash",
		PrevHash:    "Some prev hash",
		Validator:   "Validator",
	}
	return ctx.Status(fiber.StatusOK).JSON(user)
}

func CreateAccount(ctx *fiber.Ctx) error {

	body := new(model.Account)
	err := ctx.BodyParser(body)

	if err != nil {
		ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		return err
	}

	// validator address
	var address string = lib.CalculateHash(time.Now().String())

	account := model.Account{
		Name:        body.Name,
		Password:    body.Password,
		LicenseInfo: body.LicenseInfo,
		Address:     address,
		Tokens:      body.Tokens,
	}

	validators[address] = account.Tokens

	oldLastIndex := Blockchain[len(Blockchain)-1]

	// create newBlock for consideration to be forged
	newBlock, err := lib.GenerateBlock(oldLastIndex, account.LicenseInfo, address)
	if lib.IsBlockValid(newBlock, oldLastIndex) {
		candidateBlocks = append(candidateBlocks, newBlock)
	}

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
