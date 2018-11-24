package main

import (
	"bytes"
	"crypto/sha256" // sha256 단방향 해시 알고리즘 메소드 사용을 위해 import
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

func (b *Block) SetHash() { // 메소드 SetHash는 Block 타입을 리시버로 가진다. Block는 구조체로 pass by value이기 때문에 포인터를 사용해서 데이터에 접근해야 한다.
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	// 구조체 Block의 요소 Timestamp의 값을 10진수로 변환한 후, 이 10진수를 string 형으로 바꿔서 byte에 넣음
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	// 3가지 b.PrevBlockHash, b.Data, timestamp의 byte type을 공백 {}기준으로 합치고(Join함수사용) 이를 header라고 이름 붙이기
	hash := sha256.Sum256(headers)
	// func Sum256(data []byte) [Size]byte : 이 header를 sha256으로 해싱
	b.Hash = hash[:]
	// 해싱된 hash값 전체를 Block 구조체의 Hash라는 요소에 넣기
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	//Now() 메소드는 현재 로컬 시간을, Unix()메소드는 이 로컬 시간을 유닉스 타임으로 바꿔준다.
	// Block구조체 Timestamp 에는 time.Now().Unix(), Data에는 []byte(data), PrevBlockHash에는 prevBlockHash, Hash 에는 빈 바이트 타입을 넣는다.
	block.SetHash()
	// 이 세팅된 블록을 SetHash() 메소드를 통해 해싱
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Hello, Bitcoin!", []byte{})
}
