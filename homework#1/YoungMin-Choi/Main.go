package main

import (
	"strconv"
	"fmt"
)

func main() {
	bc := NewBlockchain()

	bc.AddBlock("1")
	bc.AddBlock("2")

	for _, block := range bc.blocks {
		fmt.Printf("Prev hash %x\n", block.PrevHashBlock)
		fmt.Printf("Data %x\n", block.Data)
		fmt.Printf("Hash %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("Pow: %s\n",strconv.FormatBoola(pow.Validate()))
		fmt.Println()
	}
}
