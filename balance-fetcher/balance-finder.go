package main

import (
  "fmt"
  "sync"
  "os"
  "bufio"
  "github.com/onrik/ethrpc"
  "log"
  "math/big"
  "time"
)

type Holder struct {
  Address string
  Balance big.Int
}

func getBalanceRequest(
  address string,
  block_number_hex string,
  client *ethrpc.EthRPC,
  ch chan<-bool,
) {
  if len(address) < 40 {
    ch <- true
    return
  }
  // fetch balance at block height
  balance, err := client.EthGetBalance(address, block_number_hex)
  if err != nil {
      log.Fatal(err)
  }

  fmt.Printf("%s,%s\n", address, balance.String())

  ch <- true
}

func main() {
  start := time.Now()
  // File handling
  file_size := 28049940
  // file_size := 10
  block_hex := "0x545a65"

  file, err := os.Open("./filtered_addresses_short.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)

  if err := scanner.Err(); err != nil {
      log.Fatal(err)
  }

  client := ethrpc.New("http://127.0.0.1:8545")
  // client := ethrpc.New("https://mainnet.infura.io")

  _, err = client.Web3ClientVersion()
  if err != nil {
      log.Fatal(err)
  }

  // response channel
  done_chan := make(chan bool)

  // this will allow our program to stay alive until all requests are completed
  var wg sync.WaitGroup
  wg.Add(file_size)

  // limit number of go routines running at once so we don't go over open file limit
  batch_size := 300
  if file_size < batch_size {
    batch_size = file_size
  }

  lines_read := 0
  // do each rpc call as a concurrent request
  for i := 0; i < batch_size; i++ {
    scanner.Scan()
    lines_read++
    address := scanner.Text()
    go getBalanceRequest(address, block_hex, client, done_chan)
  }

  go func() {
    for done := range done_chan {
      if done {
        if lines_read < file_size {
          scanner.Scan()
          lines_read++
          address := scanner.Text()
          go getBalanceRequest(address, block_hex, client, done_chan)
        }

        wg.Done()
      }
    }
  }()

  // wait until we have received all addresses
  wg.Wait()
  // print all addresses NOTE should stdout to a file or you are going to get spammed
  fmt.Println()
  fmt.Println("All done")
  // fmt.Println(len(addresses))
  fmt.Println(time.Since(start))
}
