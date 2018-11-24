// 블록체인 (Blockchain) 구현

package main

//배열과 맵 활영
//배열은 정렬된 해시를 유지하고 맵은 해시-블록쌍을 유지
//지금은 해시 검색 기능이 필요x - 프로토 타입 구현에서는 배열만 사용

type Blockchain struct {
	blocks []*Block
}

// 블록 추가 기능 만들기

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

// 제너시스 블록을 가지고 블록체인 생성하는 함수 구현

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
