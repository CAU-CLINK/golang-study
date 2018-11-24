package main

import (

	// sha256 단방향 해시 알고리즘 메소드 사용을 위해 import

	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0} //논스 추가(0)
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run() //Run 함수 돌린걸 논스(int 0부터)랑 hash값에 넣음

	block.Hash = hash[:] //해시를 처음부터 끝까지 블록구조체의 해시에 넣음
	block.Nonce = nonce  //논스를 이하동문하게 함
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Hello, Bitcoin!", []byte{})
}
