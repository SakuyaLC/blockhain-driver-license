package controller

import (
	lib "blockchain/main/internal"
	"blockchain/main/internal/model"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

func HandleConnection() {
	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Success")
	})

	app.Get("/account", GetAccountInfo)
	app.Post("/create-account", CreateAccount)
	app.Get("/blockchain", GetBlockchain)

	app.Listen(":80")
}

// Переписать
func HandleConn(conn net.Conn, Blockchain []model.Block, validators map[string]int, candidateBlocks chan model.Block) {
	defer conn.Close()
	var mutex = &sync.Mutex{}

	// validator address
	var address string

	// allow user to allocate number of tokens to stake
	// the greater the number of tokens, the greater chance to forging a new block
	io.WriteString(conn, "Enter token balance:")
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number: %v", scanBalance.Text(), err)
			return
		}
		t := time.Now()
		address = lib.CalculateHash(t.String())
		validators[address] = balance
		fmt.Println(validators)
		break
	}

	io.WriteString(conn, "\nEnter a new License Info:")

	scanLicenseInfo := bufio.NewScanner(conn)

	go func() {
		for {
			// take in LicenseInfo from stdin and add it to blockchain after conducting necessary validation
			for scanLicenseInfo.Scan() {
				licenseInfo := scanLicenseInfo.Text()

				mutex.Lock()
				oldLastIndex := Blockchain[len(Blockchain)-1]
				mutex.Unlock()

				// create newBlock for consideration to be forged
				newBlock, err := lib.GenerateBlock(oldLastIndex, licenseInfo, address)
				if err != nil {
					log.Println(err)
					continue
				}
				if lib.IsBlockValid(newBlock, oldLastIndex) {
					candidateBlocks <- newBlock
				}
				io.WriteString(conn, "\nEnter a new License Info:")
			}
		}
	}()

	// simulate receiving broadcast
	for {
		time.Sleep(20 * time.Second)
		mutex.Lock()
		output, err := json.Marshal(Blockchain)
		mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(conn, string(output)+"\n")
	}

}

func GetAccountInfo(ctx *fiber.Ctx) error {
	user := model.Block{
		Index:       1,
		Timestamp:   "sdf",
		LicenseInfo: "sdf",
		Hash:        "sdf",
		PrevHash:    "sdf",
		Validator:   "sdf",
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

	account := model.Account{
		Name:     body.Name,
		Password: body.Password,
	}

	return ctx.Status(fiber.StatusOK).JSON(account)
}

func GetBlockchain(ctx *fiber.Ctx, Blockchain []model.Block) error {
	//Получить блокчейн
	return ctx.Status(fiber.StatusOK).JSON(Blockchain)
}
