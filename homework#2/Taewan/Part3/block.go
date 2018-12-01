//PART 3: 영속성 및 CLI

//블록체인을 데이터베이스에 저장, 블록체인 연산을 위한 커맨드라인툴
//BoltDB : 단순하다, GO 언어, 서버가 필요 X, 데이터 설계구조 자유롭다
//키, 값 스토리지로 테이블과 행 및 열이 필요 X
//데이터 타입 X

package main

import (
  "bytes"
  "encoding/gob"
  "log"
  "time"
)

//트랜잭션 정보, 버전, 타임스탬프, 이전 블록의 해시값, 논스값 포함
type Block struct {
  Timestamp int64
  Data []byte
  PrevBlockHash []byte
  Hash []byte
  Nonce int
}

//직렬화
//[]byte 타입 값만 사용 가능 -- DB에 Block구조체를 저장하려고 한다.
//Block의 Serialize 메서드 구현

func (b *Block) serialize() []byte {
  var result bytes.Buffer
  encoder := gob.NewEncoder(&result)

  err := encoder.Encode(b)
  if err != nil {
    log.Panic(err)
  }

  return result.Bytes()
}


//새로운 블록 생성과 리턴
func NewBlock(data string, prevBlockHash []byte) *Block {
  block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
  pow := NewProofOfWork(block)
  nonce, hash := pow.Run()

  block.Hash = hash[:]
  block.Nonce = nonce

  return block
}


//제네시스블록 생성과 리턴
func NewGenesisBlock() *Block {
  return NewBlock("Genesis Block", []byte{})
}


//바이트 배열을 받아 Block을 반환하는 역직렬화 함수 필요
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
