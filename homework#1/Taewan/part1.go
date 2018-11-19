//Part 1. 블록 (Block)

//블록 : 가치있는 정보를 저장하는 데이터 구조
//트랜잭션 정보, 버전, 타임스탬프, 이전 블록의 해시값 포함

package main

func main() {

	type Block struct {
		Timestamp     int64
		Data          []byte
		PrevBlockHash []byte
		Hash          []byte
	}

}
