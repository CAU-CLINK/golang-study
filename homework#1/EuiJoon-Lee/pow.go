package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const targetBits = 5

// 목표 난이도를 상수값으로 정해준다.

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// ProofOfWork 구조체는 블록 포인터와 타겟 포인터를 가진다.

//블록 포인터를 인풋값으로 가지고, 새로운 ProofOfWork 구조체 포인터를 결과값으로 반환하는 NewProofOfWork 함수를 만든다.
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	//타겟 1로 초기화00
	target.Lsh(target, uint(256-targetBits))
	// Lsh(left shift?) target인 1을 255만큼 좌측으로 이동시키기
	pow := &ProofOfWork{b, target}
	// b와 lsh된 타겟값을 가지는 ProofOfWork구조체의 pow라는 포인터 변수 만들기
	return pow
	//pow는 ProofOfWork라는 구조체의 포인터이기 때문에 리턴값이 될 수 있음.
}

//데이터를 준비하는 메소드를 만들어보자. nonce를 인풋값으로 가지고 바이트 슬라이스를 아웃풋으로 가지는 POW구조체 메소드이다.
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash, // 오타있다. chian -> Hash
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
	//블록 구조체와 targetvBits, nonce값을 합쳐서 data 라는 바이트 슬라이스를 만든다.
}

//실행시키려는데 IntToHex가 undefine 되었다고 뜸, https://stackoverflow.com/questions/47302402/how-to-convert-int-to-hex 참고
func IntToHex(n int64) []byte {
	return []byte(strconv.FormatInt(n, 16))
}

//실제로 POW를 가동시키는 함수 Run을 만들어 보자. POW 구조체 포인터인 pow를 넣으면 정수와 바이트 슬라이스를 반환한다.
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	maxNonce := math.MaxInt64 // 이부분 빠진것 같습니다.

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for nonce < maxNonce { // maxNonce==math.MaxInt64, -9223372036854775808 to 9223372036854775807
		data := pow.prepareData(nonce) // prepareData 메소드를 통해 바이트슬라이스 만들기
		hash = sha256.Sum256(data)     // data를 sha256해싱
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 { //cmp함수는 양 값을 비교함. 이 값이 -1이라는 뜻은 hashint < pow.target라는 뜻 즉 블록해시가 난이도 타겟값보다 작을때를 의미.
			break
		} else {
			nonce++ // 아닌 경우에 논스값에 계속 1을더함
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

//pow 검증 로직
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	var isValid bool // isValid bool타입으로 선언

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	if hashInt.Cmp(pow.target) == -1 {
		isValid = true
	} else {
		isValid = false
	} // 이 부분은 제 나름대로 수정해 봤습니다. 그대로 클론코딩 하니까 에러가 나서요

	return isValid
}
