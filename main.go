package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Hash         []byte
	Data         []byte
	PreviousHash []byte
}

type BlockChain struct {
	blocks []*Block
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PreviousHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func createGenesisBlock() *Block {
	return CreateBlock("Genesis Block", []byte{})
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, new)
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{createGenesisBlock()}}
}

func main() {
	fmt.Println("Starting...")
	chain := InitBlockChain()
	chain.AddBlock("First")
	chain.AddBlock("Second")
	chain.AddBlock("Third")

	for _, block := range chain.blocks {
		fmt.Printf("PreviousHash: %x\n", block.PreviousHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n\n", block.Hash)
	}

}
