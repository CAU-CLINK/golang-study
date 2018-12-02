package main

improt (
  "fmt"
  "log"

  "github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blockBusket = "blocks"


//프로그램이 실행되는 동안 한 번 열어둔 데이터베이스를 유지하기 위해 DB 커넥션을 저장
//Blockchain 구조체
type Blockchain struct {
  tip []byte
  db *bolt.DB
}


//블록체인 반복자
type BlockchainIterator struct {
  currentHash []byte
  db *bolt.DB
}


//AddBlock 메서드
//블록을 DB에 저장
func (bc *Blockchain) AddBlock(data string) {
  var lastHash []byte

  err := bc.db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))
    lastHash = b.Get([]byte("1"))

    return nill
  })

  if err != nil {
    log.Panic(err)
  }
}

newBlock := NewBlock(data, lastHash)

err = bc.db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))
    err := b.Put(newBlock.Hash, newBlock.Serialize())
    if err != nil {
      log.Panic(err)
    }

    bc.tip = newBlock.Hash

    return nil
  })
}


//반복자는 블록체인의 블록을 반복할 때마다 만들어지며, 현재 반복의 블록 해시와 DB커넥션을 저장
//블록의 해시와 DB 커넥션이 필요하기 때문에 반복자는 논리적으로 블록체인과 연결되어 있으며
//이는 Blockchain의 메서드로 만든다.
func (bc *Blockchain) Iterator() *BlockchainIterator {
  bci := &BlockchainIterator{bc.tip, bc.db}
  return bci
}


//BlockchainIterator
//블록체인으로부터 다음 블록을 반환하는 일 수행
func (i *BlockchainIterator) Next() *Block {
  var block *Block

  err := i.db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blockBucket))
    encodedBlock := b.Get(i.currentHash)
    block = DeserializeBlock(encodeBlock)

    return nil
  })

  if err != nil {
    log.Panic(err)
  }

  i.currentHash = block.PrevBlockHash

  return block
}


//영속성

//NewBlockchain - 새로운  Blockchain 인스턴스 생성, 여기에 제네시스 블록 추가

//DB 파일을 연다.
//저장된 블록체인이 있는지 확인
//블록체인 존재 O - 새로운 Blockchain 인스턴스 생성, Blockchain 인스턴스의 끝부분을 DB에 저장된 마지막 블록의 해시로 설정
//블록체인 존재 X - 제네시스 블록 생성, DB에 저장, 제네시스 블록의 해시를 마지막 블록의 해시로 저장, 제네시스 블록을 끝부분으로 하는 새로운 Blockchain인스턴스 생성
func NewBlockchain() *Blockchain {
  var tip []byte
  db, err := bolt.Open(dbFile, 0600, nil)
  if err != nil {
    log.Panic(err)
  }

  err = db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))

    if b == nil {
      fmt.Println("No existing blockchain found. Creating a new one...")
      genesis := NewFenesisBlock()

      b, err := tx.CreateBucket([]byte(blocksBucket))
      if err != nil {
        log.Panic(err)
      }

      err = b.Put(genesis.Hash, genesis.Serialize())
      if err != nil {
        log.Panic(err)
      }

      err = b.Put([]byte("1"), genesis.Hash)
      if err != nil {
        log.Panic(err)
      }
      tip = genesis.Hash
    } else {
      tip = b.Get([]byte("1"))
    }

    return nil
  })

  if err != nil {
    log.Panic(err)
  }

  bc := Blockchain{tip, db}

  return &bc
}
