package main

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain) AddNewBlock(data string) {
	// blockchian이라는 구조체에 새로운 block이라는 구조체 타입의 요소를 추가하는 blockchain을 메소드로 가지는 함수.
	prevBlock := bc.blocks[len(bc.blocks)-1]
	//Blockchain 구조체의 마지막 블록을 prevBlock이라고 한다.
	newBlock := NewBlock(data, prevBlock.Hash)
	// newBlock은 data와 prevBlock 구조체의 Hash라는 요소를 NewBlock 함수에 넣어서 만든 새로운 Block 타입의 구조체이다.
	bc.blocks = append(bc.blocks, newBlock)
	// Blockchain구조체에 newBlock 구조체를 추가한다.
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
	// 빈 PrevBlockHash 값과 "Hello, Bitcoin!"라는 string 타입의 data 변수를 가지는 제네시스 블록을 만든다.
}
