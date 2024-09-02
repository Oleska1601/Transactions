package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
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

func generateTr(c Client) {
	//выбираем случайного клиента и проверяем, чтобы это не был тот же самый клиент
	var (
		randClientID int
		randAmount   float64
	)
	for {
		randClientID = rand.Intn(len(clients)) + 1 //[1;len(clients)] -> randIntn(len(clients)-1+1)+1
		if randClientID != c.ID {
			break
		}
	}
	//выбираем рандомно тип счета для транзакции
	typesOfAccount := []string{"credit", "debit"}
	typeOfAccount := typesOfAccount[rand.Intn(len(typesOfAccount))]
	if typeOfAccount == "credit" {
		randAmount = rand.Float64()*(c.CreditAccount) + 1 // [1;c.CreditAccount]
	} else {
		randAmount = rand.Float64()*(c.DebitAccount) + 1
	}
	tr := Transaction{c.ID, randClientID, randAmount, typeOfAccount}
	trProcessing(tr)

}

func trProcessing(tr Transaction) {
	m.Lock()
	defer m.Unlock()
	if _, ok := clients[tr.SrcClientId]; !ok {
		clients[tr.SrcClientId] = Client{tr.SrcClientId, 0, 0}
		return //транзакция в любом случае будет невозможна, тк у клиента оба баланса 0
	}
	if _, ok := clients[tr.DstClientId]; !ok {
		clients[tr.DstClientId] = Client{tr.DstClientId, 0, 0}
	}
	scrClient := clients[tr.SrcClientId]
	dstClient := clients[tr.DstClientId]
	switch tr.Type {
	case "credit":
		if scrClient.CreditAccount-tr.Amount < 0 {
			fmt.Printf("client %d CreditAccount is not enough\n", tr.SrcClientId)
		} else {
			scrClient.CreditAccount -= tr.Amount
			dstClient.CreditAccount += tr.Amount
		}
	case "debit":
		if scrClient.DebitAccount-tr.Amount < 0 {
			fmt.Printf("client %d DebitAccount is not enough\n", tr.SrcClientId)

		} else {
			scrClient.DebitAccount -= tr.Amount
			dstClient.DebitAccount += tr.Amount
		}
	default:
		fmt.Printf("incorrect type %s", tr.Type)
		return
	}
	clients[tr.SrcClientId] = scrClient
	clients[tr.DstClientId] = dstClient
}

func Clear() {
	//exec.Command - создание объекта cmd
	//1) cmd - испольняемый файл - интерпретатор командной строки Windows
	//2) /c - выполнить команду и завершиться
	//3) cls - очистка ком строки в Windows
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	//определяем, куда выводятся данные - моя консоль
	cmd.Stdout = os.Stdout
	cmd.Run()

}

func main() {
	ticker := time.NewTicker(5 * time.Second)
	for _, tr := range transactions {
		wg.Add(1)
		go func(tr Transaction) {
			defer wg.Done()
			trProcessing(tr)
		}(tr)
	}
	wg.Wait()

	//возможность добавления новых транзакций в процессе выполнения программы
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
		Clear()
		//каждые 5 сек генерируем транзакции для каждого клиента
		for _, cl := range clients {
			generateTr(cl)
		}
		for id, client := range clients {
			fmt.Printf("ClientId: %d [Credit: %.2f, Debit: %.2f]\n", id, client.CreditAccount, client.DebitAccount)
		}
	}
}
