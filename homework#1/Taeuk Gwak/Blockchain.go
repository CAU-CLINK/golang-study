type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockhash []byte
	Blockhash     []byte
  Version       []byte
  Merklehash    []byte
  Nonce         int
}
 func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockhash, b.Data, timestamp}, []byte{})
	blockhash := sha256.Sum256(headers)
  
	b.Blockhash = blockhash[:]
}
 func NewBlock(data string, prevBlockhash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockhash, []byte{}}
	block.SetHash()
	return block
}
 type Blockchain struct {
	blocks []*Block
}
 func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Blockhash)
	bc.blocks = append(bc.blocks, newBlock)
}
 func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
 func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
 const targetBits = 24
 type ProofOfWork struct {
	block  *Block
	target *big.Int
}
 func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
 	pow := &ProofOfWork{b, target}
 	return pow
}
 func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockhash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
 	return data
}
 func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
 	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
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
 func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockhash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
 	block.Blockhash = hash[:]
	block.Nonce = nonce
 	return block
}
 func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
 	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
 	isValid := hashInt.Cmp(pow.target) == -1
 	return isValid
}
 func main() {
	bc := NewBlockchain()
 	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")
 	for _, block := range bc.blocks {
    fmt.Printf("Blockhash: %x\n", block.Blockhash)
    fmt.Printf("Version: %s\n", block.Version)
		fmt.Printf("Previous Blockhash: %x\n", block.PrevBlockhash)
    fmt.Printf("Merklehash: %x\n", block.Merklehash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Println()
    
    pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
