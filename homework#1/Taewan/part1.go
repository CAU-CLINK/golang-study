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


  //해시 계산
  //블록을 구성하는 필드들을 하나로 이은 뒤 이어진 문자열에 대해 SHA-256 해시를 계산
  //SetHash 메서드 작성

  func (b *Block) SetHash() {
          timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
          headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
          hash := sha256.Sum256(headers)
          b.Hash = hash[:]
  }


  //Go의 컨벤션을 따라, 블록을 생성하는 함수 작성

    func NewBlock(data string, prevBlockHash []byte) *Block {
      block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
      block.SetHash()
      return block
    }


// 블록체인 (Blockchain) 구현

//배열과 맵 활영
//배열은 정렬된 해시를 유지하고 맵은 해시-블록쌍을 유지
//지금은 해시 검색 기능이 필요x - 프로토 타입 구현에서는 배열만 사용

type Blockchain struct {
        blocks []*Block
  }


   
}
