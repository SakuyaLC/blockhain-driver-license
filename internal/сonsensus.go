package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"blockchain/main/internal/model"
)

// pickWinner creates a lottery pool of validators and chooses the validator who gets to forge a block to the blockchain
// by random selecting from the pool, weighted by amount of tokens staked
func PickWinner(Blockchain []model.Block, tempBlocks []model.Block, candidateBlocks chan model.Block, validators map[string]int) []model.Block {
	time.Sleep(20 * time.Second)

	temp := tempBlocks

	lotteryPool := []string{}
	if len(temp) > 0 {

		// slightly modified traditional proof of stake algorithm
		// from all validators who submitted a block, weight them by the number of staked tokens
		// in traditional proof of stake, validators can participate without submitting a block to be forged
	OUTER:
		for _, block := range temp {
			// if already in lottery pool, skip
			for _, node := range lotteryPool {
				if block.Validator == node {
					continue OUTER
				}
			}

			// lock list of validators to prevent data race
			setValidators := validators

			k, ok := setValidators[block.Validator]
			if ok {
				for i := 0; i < k; i++ {
					lotteryPool = append(lotteryPool, block.Validator)
				}
			}
		}

		// randomly pick winner from lottery pool
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]

		// add block of winner to blockchain and let all the other nodes know
		for _, block := range temp {
			if block.Validator == lotteryWinner {
				Blockchain = append(Blockchain, block)
				break
			}
		}
	}

	tempBlocks = []model.Block{}

	return Blockchain
}

// isBlockValid makes sure block is valid by checking index
// and comparing the hash of the previous block
func IsBlockValid(newBlock, oldBlock model.Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if CalculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// SHA256 hasing
// calculateHash is a simple SHA256 hashing function
func CalculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// calculateBlockHash returns the hash of all block information
func CalculateBlockHash(block model.Block) string {
	record := string(block.Index) + block.Timestamp + block.LicenseInfo + block.PrevHash
	return CalculateHash(record)
}

// generateBlock creates a new block using previous block's hash
func GenerateBlock(oldBlock model.Block, licenseInfo string, address string) (model.Block, error) {

	var newBlock model.Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.LicenseInfo = licenseInfo
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = CalculateBlockHash(newBlock)
	newBlock.Validator = address

	return newBlock, nil
}

func ShowBlockchain(Blockchain []model.Block) {
	time.Sleep(22 * time.Second)
	fmt.Println("New Blockchain is:")
	fmt.Println(Blockchain)
	fmt.Println("")
}
