package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

//subsidy는 보상의 양
const subsidy = 10

//비트코인 트랜잭션
//트랜잭션은 입력과 출력의 조합
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

//거래가 코인베이스인지 아닌지 확인 하는 것
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

//SetID sets ID of a transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

//트랜잭션 입력
//입력 구조체 - 이전 출력을 참조
//Txid는 해당 트랜잭션 ID 저장. Vout은 트랜잭션의 출력 인덱스 저장.
//Scrptsig는 출력의 ScriptPubkey에 사용되는 데이터 제공하는 스크립트
type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

//출력 구조체
type TXOutput struct {
	Value        int
	ScriptPubkey string
}

//미사용 트랜잭션 출력
//입력과 출력에 대한 잠금-해제 메서드 정의
//CanUnlockOutputWith는 이 주소가 거래를 시작했는지 안했는지 확인
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

//CanBeUnlockedWith는 output을 provided data와 함께 잠굴 수 있는 지 확인
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubkey == unlockingData
}

//코인베이스 트랜잭션
//단 하나의 입력만 가진다.
//Txid는 비어있으며 Vout은 -1이다.
//ScriptSig에 아무 스크립트도 저장하지 않으며 대신 임의의 데이터가 저장된다.
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

//코인 전송
//새로운 트랜잭션을 생성하여 블록에 넣고 블록을 채굴
//이제 일반 트랜잭션 구현

//FindSpendableOutputs
//새로운 출력을 생성하기 전에, 모든 미사용 출력을 찾아 충분한 잔고가 있는지 확인
//확인이 끝나면 찾아낸 각각의 출력에 대해 이를 참조하는 입력들이 생성된다.
//1. 수신자 주소로 잠근 출력 (실제로 다른 주소로 코인을 전송하는 출력)
//2. 발신자 주소로 잠근 출력 (이는 잔액, 미사용 출력들의 보유량이 새로운 트랜잭션에서 필요한 값보다 큰 경우에만 만들어진다.
//출력은 나눠질 수 없다.)
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	//입력 리스트 생성
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	//출력 리스트 생성
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
