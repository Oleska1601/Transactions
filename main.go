package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	ID            int
	CreditAccount float64
	DebitAccount  float64
}

type Transaction struct {
	SrcClientId int
	DstClientId int
	Amount      float64
	Type        string // "credit" or "debit"
}

var (
	m  sync.RWMutex
	wg = new(sync.WaitGroup)

	clients = map[int]Client{
		1: {1, 2000.0, 5000.0},
		2: {2, 1000.0, 3000.0},
		3: {3, 3000.0, 2000.0},
	}

	transactions = []Transaction{
		{1, 2, 100.0, "credit"},
		{2, 3, 50.0, "debit"},
		{1, 2, 75.0, "debit"},
		{2, 1, 125.0, "credit"},
		{3, 2, 200.0, "credit"},
		{1, 3, 200.0, "debit"},
		{1, 2, 200.0, "debit"},
	}
)

func trProcessing(tr Transaction) {
	m.Lock()
	defer m.Unlock()
	scrClient := clients[tr.SrcClientId]
	dstClient := clients[tr.DstClientId]
	if tr.Type == "credit" {
		if scrClient.CreditAccount-tr.Amount < 0 {
			fmt.Printf("client %d CreditAccount is not enough\n", tr.SrcClientId)
		} else {
			scrClient.CreditAccount -= tr.Amount
			dstClient.CreditAccount += tr.Amount
		}

	} else {
		if scrClient.DebitAccount-tr.Amount < 0 {
			fmt.Printf("client %d DebitAccount is not enough\n", tr.SrcClientId)

		} else {
			scrClient.DebitAccount -= tr.Amount
			dstClient.DebitAccount += tr.Amount
		}
	}
	clients[tr.SrcClientId] = scrClient
	clients[tr.DstClientId] = dstClient
}

func main() {
	ticker := time.NewTicker(time.Second)
	for _, tr := range transactions {
		wg.Add(1)
		go func(tr Transaction) {
			defer wg.Done()
			trProcessing(tr)
		}(tr)
	}
	wg.Wait()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			str, _ := reader.ReadString('\n')
			transaction := strings.Split(str, " ")
			srcClientID, _ := strconv.Atoi(transaction[0])
			dstClientID, _ := strconv.Atoi(transaction[1])
			amount, _ := strconv.ParseFloat(transaction[2], 64)
			trProcessing(Transaction{srcClientID, dstClientID, amount, transaction[3]})
		}
	}()

	for range ticker.C {
		for _, client := range clients {
			fmt.Printf("client %d has CreditAccount %f and DebitAccount %f\n", client.ID, client.CreditAccount, client.DebitAccount)
		}
	}
}
