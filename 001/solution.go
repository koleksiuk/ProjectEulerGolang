package main

import (
  "fmt"
)

func isDivisible(num int) bool {
  return num % 3 == 0 || num % 5 == 0
}

func compute(start int, end int, computeChan chan int, haltChan chan bool) {
  fmt.Println("Starting channel", start, end)
  for i := start; i <= end; i++ {
    if isDivisible(i) {
      computeChan<- i
    }
  }

  haltChan<- true
}

func listener(computeChan chan int, haltChan chan bool, resultChan chan int) {
  closedChannels := 0
  sum := 0

  for {
    select {
      case num := <-computeChan:
        sum += num
      case <-haltChan:
        closedChannels += 1
      default:
        // no-op
    }

    if closedChannels == 10 {
      break;
    }
  }

  close(computeChan)
  close(haltChan)

  for num := range computeChan {
    sum += num
  }

  resultChan<- sum
}

func main() {
  haltChan    := make(chan bool, 10)
  computeChan := make(chan int,  200)
  resultChan  := make(chan int,  1)

  for i := 0; i <= 9; i++ {
    go compute(i * 100, (i + 1) * 100 - 1, computeChan, haltChan)
  }

  go listener(computeChan, haltChan, resultChan)

  fmt.Println(<-resultChan)
}
