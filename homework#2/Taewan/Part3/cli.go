package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "strconv"
)


//커맨드라인과 관련된 모든 연산들은 CLI 구조체에 의해 처리
type CLI struct {
  bc *Blockchain
}

func (cli *CLI) printUsage() {
  fmt.Println("Usage:")
  fmt.Println(" addblock -data BLOCK_DATA - add a block to the blockchain")
  fmt.Println(" printchain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
  if len(os.Args) < 2 {
    cli.printUsage()
    os.Exit(1)
  }
}


//어떤 하위커맨드가 파싱되었는지 확인한 뒤 관련 함수를 실행
func (cli *CLI) addBlock(data string) {
  cli.bc.AddBlock(data)
  fmt.Println("Success!")
}

func (cli *CLI) printChain() {
  bci := cli.bc.Iterator()

  for {
    block := bci.Next()

    fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
    fmt.Pritnf("Data: %s\n", block.Data)
    fmt.Printf("Hash: %x\n", block.Hash)
    pow := NewProofOfWork(block)
    fmt.Prinf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
    fmt.Println()

    if len(block.PrevBlockHash) == 0 {
      break
    }
  }
}


//CLI의 엔트리포인트는 Run함수
func (cli *CLI) Run() {
  cli.validateArgs()

  addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
  printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

  addBlockData := addBlockCmd.String("data", "", "Block data")

  switch os.Args[1] {
  case "addblock":
    err := addBlockCmd.Parse(os.Args[2:])
    if err != nil {
      log.Panic(err)
    }
  case "printchain":
    err := printChainCmd.Parge(os.Args[2:])
    if err != nil {
      log.Panic(err)
    }
  default:
    cli.printUsage()
    os.Exit(1)
  }

  if addBlockCmd.Parsed() {
    if *addBlockData == "" {
      addBlockCmd.Usage()
      os.Exit(1)
    }
    cli.addBlock(*addBlockData)
  }

  if printChainCmd.Parsed() {
    cli.printChain()
  }
}
