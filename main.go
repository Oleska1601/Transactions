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

type Transaction struct {
	ClientID string
	Amount   float64
	Type     string // "credit" or "debit"
}

type Client struct {
	ID string
}

var (
	m            sync.RWMutex
	wg           = new(sync.WaitGroup)
	balances     = make(map[string]float64)
	clients      = []string{} //содержит ID клиентов
	transactions = []Transaction{
		{"client1", 100.0, "credit"},
		{"client2", 50.0, "debit"},
		{"client1", 75.0, "debit"},
		{"client2", 125.0, "credit"},
		{"client3", 200.0, "credit"},
		{"client1", 200.0, "debit"},
		{"client1", 200.0, "debit"},
	}
)

// ниже пока не до конца
/*
func generateTr(c Client) {
	m.RLock()
	defer m.RUnlock()
	curAmount := balances[c.ID]
	randAmount := rand.Float64()*(2*curAmount+1) - curAmount //[-curAmount; curAmount] rand * (finish - start + 1) + start
	if randAmount < 0 {
		transactions = append(transactions, Transaction{c.ID, randAmount, "credit"})
	} else {
		//генерируем др клиента, у которого есть нужный баланс для совершения нашей транакции debit и пополнения своего баланса
		var randClientIdx int
		var anotherClientID string
		//мб никогда не найдет нужного клиента -ДОДЕЛАТЬ
		for {
			randClientIdx = rand.Intn(len(clients))
			anotherClientID = clients[randClientIdx]
			if anotherClientID != c.ID && balances[anotherClientID] >= randAmount {
				break
			}
		}
		balances[anotherClientID] -= randAmount
		transactions = append(transactions, Transaction{c.ID, randAmount, "debit"})
	}
}
*/

// обработка транзакции
func trProcessing(tr Transaction) {
	m.Lock()
	defer m.Unlock()
	if _, ok := balances[tr.ClientID]; !ok {
		clients = append(clients, tr.ClientID)
		balances[tr.ClientID] = 0
	}
	if tr.Type == "debit" {
		balances[tr.ClientID] += tr.Amount
	} else {
		if balances[tr.ClientID]-tr.Amount < 0 {
			fmt.Printf("transaction for %s with amount %f %s is impossible\n", tr.ClientID, tr.Amount, tr.Type)
			return
		} else {
			balances[tr.ClientID] -= tr.Amount
		}
	}

}

func main() {
	ticker := time.NewTicker(5 * time.Second) //пока 5 секунд для более удобной проверки
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
			amount, _ := strconv.ParseFloat(transaction[1], 64)
			tr := Transaction{transaction[0], amount, transaction[2]}
			trProcessing(tr)
			wg.Done()
		}
	}()
	for range ticker.C {
		for key, val := range balances {
			fmt.Println(key, val)
		}
	}

}
