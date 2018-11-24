package main

import(
	"fmt"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}


func (cli *CLI) printUsage(){
	fmt.Println("Usage")
	fmt.Println(" addblock -data BLOCK_DATA - add a block to the blockchaib")
	fmt.Println("   printchain - print all the blocks of the blockchain")
}


func (cli *CLI) validateArgs(){
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}


func (cli *CLI) addblock(data string){
	cli.bc.AddBlock(data)
	fmt.Println("Success")
}

func (cli *CLI) printChan(){
	bci := cli.bc.Iterator()


	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n",block.PrevHashBlock)
		fmt.Printf("Data %s\n",block.Data)
		fmt.Printf("Hash: %x\n",block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("Pow: %s \n",strconv.FormatBool(pow.Validate()))
		fmt.Println

		if len(block.PrevHashBlock) == 0 {
			break
		}

	}
}


func 