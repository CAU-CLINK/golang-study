package main

import (
  "encoding/hex"
  "fmt"
  "log"
  "os"

  "github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"


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


//Blockchain.MineBlock 메서드
//블록을 DB에 저장
//MineBlock은 provided transactions들과 함께 새로운 블록을 채굴
func (bc *Blockchain) MineBlock(transactions []*Transaction) {
  var lastHash []byte

  err := bc.db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))
    lastHash = b.Get([]byte("1"))

    return nil
  })

  if err != nil {
    log.Panic(err)
  }

  newBlock := NewBlock(transactions, lastHash)

  err = bc.db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))
    err := b.Put(newBlock.Hash, newblock.Serialize())
    if err != nil {
      log.Panic(err)
    }

    err = b.Put([]byte("1"), newBlock.Hash)
    if err != nil {
      log.Panic(err)
    }

    bc.tip = newBlock.Hash

    return nil
  })
}


//미사용 출력을 포함하는 트랜잭션을 찾는 작업
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
  var unspendTXs []Transaction
  spentTXOs := make(map[string][]int)
  bci := bc.Iterator()

  for {
    block := bci.Next()

    for _, tx := range block.Transactions {
      txID := hex.EncodeToString(tx.ID)

      output:
        for outIdx, out := range tx.Vout {
          //출력 사용 여부 검사
          //어떤 출력이 우리가 찾고 있는 미사용 트랜잭션의 출력과 동일한 주소로 잠궈져있다면, 정답!
          //그러나 출력을 가져오기 전에 먼저 이 출력이 입력에서 이미 참조되었는지 확인할 필요가 있다.
          if spentTXOs[txID] != nil {
            for _, spentOut := range spentTXOs[txID] {
              IF spentOut == outIdx {
                continue Outputs
              }
            }
          }

          //출력 시작
          //입력에서 이미 참조된 출력들은 무시 - 이미 다른 출력으로 이동해서 카운트 불가
          if out.CanBeUnlockedWith(address) {
            unspentTXs = append(unspentTXs, *tx)
          }
        }

        //출력 검사가 끝나면 주어진 주소로 잠긴 출력을 해제할 수 있는 모든 입력들을 가져온다.
        //코인베이스 트랜잭션은 제외
        //밑의 함수는 미사용된 출력을 포함하고 있는 트랜잭션의 리스트를 반환
        if tx.IsCoinbase() == false {
          for _, in := range tx.Vin {
            if in.CanUnlockOutputWith(address) {
              inTxID := hex.EncodeToString(in.Txid)
              spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
            }
          }
        }
    }

    if len(block.PrevBlockHash) == 0 {
      break
    }
  }

  return unspentTXs
}


//잔고 게산을 위해선 트랜잭션 리스트에서 출력들만 반환하는 함수가 하나 더 필요하다.
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
  var UTXOs []TXOutput
  unspentTransactions := bc.FindUnspentTransactions(address)

  for _, tx := range unspentTransactions {
    for _, out := range tx.Vout {
      if out.CanBeUnlockedWith(address) {
        UTXOs = append(UTXOs, out)
      }
    }
  }

  return UTXOs
}


//FindSpendableOutputs 메서드는 이전에 정의한 FindUnspentTransactions을 기반으로 한다.
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
  unspentOutputs := make(map[string][]int)
  unspentTXs := bc.FindUnspentTransactions(address)
  accumulated := 0

Work:
  for _, tx := range unspentTXs {
    txID := hex.EncodeToString(tx.ID)

    for outIdx, out := range tx.Vout {
      if out.CanBeUnlockedWith(address) && accumulated < amount {
        accumulated += out.Value
        unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

        if accumulated >= amount {
          break work
        }
      }
    }
  }

  return accumulated, unspentOutputs
}


//반복자는 블록체인의 블록을 반복할 때마다 만들어지며, 현재 반복의 블록 해시와 DB커넥션을 저장
//블록의 해시와 DB커넥션이 필요하기 때문에 반복자는 논리적으로 블록체인과 연결되어 있으며
//이는 Blockchain의 메서드로 만든다.
func (bc *Blockchain) Iterator() *BlockchainIterator {
  bci := &BlockchainIterator{bc.tip, bc.db}

  return bci
}


//BlockchianIterator
//블록체인으로부터 다음 블록을 반환하는 일 수행
func (i *BlockchainIterator) Next() *Block {
  var block *Block

  err := i.db.View(func(tx.*bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))
    encodedBlock := b.Get(i.currentHash)
    block = DeserializeBlock(encodedBlock)

    return nil
  })

  if err != nil {
    log.Panic(err)
  }

  i.currentHash = block.PrevBlockHash

  return block
}

func dbExists() bool {
  if _, err := os.Stat(dbFile); os.IsNotExist(err) {
    return false
  }

  return true
}


//영속성
//NewBlockchain - 새로운 Blockchain 인스턴스 생성, 여기에 제네시스 블록 추가

//DB 파일을 연다.
//저장된 블록체인이 있는지 확인
//블록체인 존재 O - 새로운 Blockchain 인스턴스 생성, Blockchain 인스턴스의 끝부분을 DB에 저장된 마지막 블록의 해시로 설정
//블록체인 존재 X - 제네시스 블록 생성, DB에 저장, 제네시스 블록의 해시를 마지막 블록의 해시로 저장, 제네시스 블록을 끝부분으로 하는 새로운 Blockchain 인스턴스 생성

//NewBlockchain은 제네시스블록과 함께 새로운 블록체인을 만든다.
func NewBlockchain(address string) *Blockchain {
  if dbExists() == false {
    fmt.Println("No existing blockchain found. Creat one first.")
    os.Exit(1)
  }

  var tip []byte
  db. err := bolt.Open(dbFile, 0600, nil)
  if err != nil {
    log.Panic(err)
  }

  err = db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))
    tip = b.Get([]byte("1"))

    return nil
  })

  if err != nil {
    log.Panic(err)
  }

  bc := Blockchain{tip, db}

  return &bc
}


//CreateBlockchain은 새로운 블록체인 DB를 생성
//새로운 블록체인을 만드는 함수
func CreatBlockchain(address string) *Blockchain {
  if dbExists() {
    fmt.Println("Blockchain already exists.")
    os.Exit(1)
  }

  var tip []byte
  db, err := bolt.Open(dbFile, 0600, nil)
  if err != nil {
    log.Panic(err)
  }

  err = db.Update(func(tx *bolt.Tx) error {
    cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
    genesis := NewGenesisBlock(cbtx)

    b, err := tx.CreatBucket([]byte(blocksBucket))
    if err != nil {
      log.Panic(err)
    }

    err = b.Put(genesis.Hash, genesis.Serialize())
    if err != nil {
      log.Panic(err)
    }

    err = b.Put([]byte("1"), gensis.Hash)
    if err != nil {
      log.Panic(err)
    }
    tip = genesis.Hash

    return nil
  })

  if err != nil {
    log.Panic(err)
  }

  bc := Blockchain{tip, db}

  return &bc
}
