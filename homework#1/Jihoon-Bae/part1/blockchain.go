package main

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
} //[]*Block{} 1. struct 초기화시 인자전달, 또는 포인터값 전달가능. 2. 스트럭트 배열 생성시에도 그냥 스트럭트 초기화 방식으로 가능.
