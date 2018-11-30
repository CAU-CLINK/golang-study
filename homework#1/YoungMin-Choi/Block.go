package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevHashBlock []byte
	Hash          []byte
	Nonce         int
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()

}

func NewBlock(data string, prevHashBlock []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevHashBlock, []byte{}, 0}
	pow := NewProofOfWork(block)
	//구조체를 포인터로 받아와서 설정한다.
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	//해쉬값을 가지고 돌린다
	return block
	//해쉬값을 반환
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
