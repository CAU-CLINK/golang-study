//블록체인에 트랜잭션 저장하기

package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

//트랜잭션 정보, 버전, 타임스탬프, 이전 블록의 해시값, 논스값 포함
//블록 채굴 - 모든 블록은 적어도 하나의 트랜잭션을 가짐
//Data 필드 제거 후 트랜잭션에 저장해야함을 의미
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

//직렬화
//[]byte 타입 값만 사용 가능 -- DB에 Block 구조체를 저장
//Block의 Serialize 메서드 구현

//우선 직렬화된 데이터를 저장할 버퍼를 선언하고 gob 인코더를 초기화한 뒤 블록을 인코딩하면 바이트 배열 반환
func (b *Block) serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

//새로운 블록 생성과 리턴 - 수정 진행
func NewBlock(transaction []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

//제네시스 블록 생성과 리턴
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

//바이트 배열을 받아 Block을 반환하는 역질렬화 함수 필요
//메서드가 아닌 독립적인 함수로 구현
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
