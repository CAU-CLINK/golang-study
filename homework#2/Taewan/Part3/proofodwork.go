package main

import(
  "bytes"
  "crypto/sha256"
  "fmt"
  "math"
  "math/big"
)

var (
  maxNonce = math.MaxInt64
)


// 채굴 난이도 정의
const targetBits = 24


// ProofOfWork 구조체 구성 - 블록 포인터, 타켓 포인터
type ProofOfWork struct {
  block *Block
  target *big.Int
}


//NewProofOfWork 함수에서 bit.Int을 1로 초기화
//256 - targetBits 비트만큼 좌측 시프트 연산
func NewProofOfWork(b *Block) *ProofOfWork {
  target := big.NewInt(1)
  target.Lsh(target, uint(256-targetBits))

  pow := &ProofOfWork{b, target}

  return pow
}


//해시 계산을 위한 데이터가 필요 -> 데이터 준비
//블록의 필드값들과 타겟 및 논스값을 병합하는 직관적인 코드다.
func (pow *ProofOfWork) prepareData(nonce int) []byte {
  data := bytes.Join(
    [][]byte{
      pow.block.PrevBlockHash,
      pow.block.Data,
      IntToHex(pow.block.Timestamp)
      IntToHex(int64(targetBits))
      IntToHex(int64(nonce))
    },
    []byte{}
  )

  return data
}

//proof-of-work 알고리즘 핵심 코드 구현
func (pow *ProofOfWork) Run() (int, []byte) {
  var hashInt big.Int
  var hash [32]byte
  nonce := 0

  fmt.Println("Mining the block containing \%s\"\n", pow.block.Data)
  for nonce < maxNonce {
    data := pow.prepareDate(nonce)

    hash = sha256.Sum256(data)
    fmt.Printf("\r%x", hash)
    hashInt.SetBytes(hash[:])

    if hashInt.Cmp(pow.target) == -1 {
      break
    } else {
      nonce++
    }
  }
  fmt.Print("\n\n")

  return nonce, hash[:]
}


//작업 증명을 검증할 수 있는 기능
func (pow *ProofOfWork) Validate() bool {
  var hashInt big.Int

  data := pow.prepareData(pow.block.Nonce)
  hash := sha256.Sum256(data)
  hashInt.SetBytes(hash[:])

  isValid := hashInt.Cmp(pow.target) == -1

  return isValid
}
