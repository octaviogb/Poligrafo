package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jdkato/prose/summarize"
)

// Transaction é a trasação no contrato que a gente pediu no teste
type Transaction struct {
	// Description tem que ser uma descrição legível
	Description string `json:"descricao"`
	// Value tem que ser um valor inteiro, cujas duas últimas casas
	// representam os centavos
	Value int32 `json:"valor"`
	// Data tem que ser uma data inteira representando um
	// timestamp (secs since 1970)
	Data int64 `json:"data"`
	// Duplicated é um booleano que representa se essa transação
	// está duplicada... mas se eu tiver duas transações duplicadas
	// uma tem que estar true e a outra false
	Duplicated bool `json:"duplicated"`
}

// Application vai guardar a URL... sei lá pq fiz assim
type Application struct {
	url string
}

func (a Application) tryAccess() error {
	_, err := http.Get(a.url + "/9505/transactions/2020/02")

	return err
}

// Result have original error and friendly message
type Result struct {
	err     error
	message string
}

func createResult(err error, message string, results *[]Result) []Result {
	r := append(*results, Result{
		err:     err,
		message: message,
	})
	return r
}

func (a Application) verifyContract() []Result {
	var results []Result
	url := a.url + "/9505/transacoes/2020/02"
	results = createResult(fmt.Errorf("%s", url), "Vou pegar transacoes aqui", &results)

	result, err := http.Get(url)

	if err != nil {
		results = createResult(err, "Nao consegui acessar /9505/transactions/2020/02", &results)
		return results
	}

	if result.StatusCode >= 400 {
		results = createResult(fmt.Errorf("%s", result.Status), "Deu erro HTTP aqui", &results)
		return results
	}

	data, err := ioutil.ReadAll(result.Body)
	defer result.Body.Close()

	if len(data) == 0 {
		results = createResult(fmt.Errorf("%s", data), "Tem nada nessa chamada meo!?", &results)
		return results
	}

	var transactions []Transaction
	err = json.Unmarshal([]byte(data), &transactions)

	if err != nil {
		results = createResult(err, "Nao consegui parsear o json", &results)
		return results
	}

	if len(transactions) == 0 {
		results = createResult(err, "Nao achei nenhuma transacao", &results)
		return results
	}

	firstTransaction := transactions[0]

	if firstTransaction.Description == "" && firstTransaction.Data == 0 && firstTransaction.Value == 0 {
		results = createResult(fmt.Errorf(""), "Nao parseou nada, o contrato esta errado", &results)
		return results
	}

	description := firstTransaction.Description
	if len(description) < 3 {
		results = createResult(fmt.Errorf("%s", description), "Descricao vazia ou nao consegui ler", &results)
	}

	readabilityScore := 0.0
	for _, transaction := range transactions[:3] {
		doc := summarize.NewDocument(transaction.Description)
		readabilityScore += doc.SMOG()
	}
	if readabilityScore < 30 {
		results = createResult(fmt.Errorf("Score de legibilidade SMOG %f", readabilityScore/3), "Nao parece que da pra humano ler", &results)
	}

	for _, transaction := range transactions {
		ts := transaction.Data / 1000
		tm := time.Unix(ts, 0)
		month := tm.Month()
		if month != 2 {
			results = createResult(fmt.Errorf("%s %d", tm, ts), "Tem uma transacao com mes errado", &results)
		}
	}

	negatives := 0
	positives := 0
	zero := 0
	hundredminus := 0
	for _, transaction := range transactions {
		if transaction.Value < 0 {
			negatives++
		}
		if transaction.Value >= 0 {
			positives++
		}
		if transaction.Value == 0 {
			zero++
		}
		if transaction.Value < 10000 {
			hundredminus++
		}
	}
	if negatives == 0 {
		results = createResult(fmt.Errorf("%d", negatives), "Nenhuma transacao negativa", &results)
	}
	if positives == 0 {
		results = createResult(fmt.Errorf("%d", positives), "Nenhuma transacao positiva", &results)
	}
	if zero != 0 {
		results = createResult(fmt.Errorf("%d", zero), "Legal, ganhei brinde no cartao de credito eh? Transacao de 0,00?", &results)
	}
	if hundredminus == 0 {
		results = createResult(fmt.Errorf("%d", hundredminus), "Bichao rico, soh gasta mais de 100,00", &results)
	} else {
		results = createResult(fmt.Errorf("%d", hundredminus), "Transacoes abaixo de 100,00", &results)
	}

	var duplicated *Transaction
	duplicated = nil
	for _, transaction := range transactions {
		if transaction.Duplicated {
			duplicated = &transaction
			continue
		}
	}
	duplicateds := 0
	if duplicated != nil {
		for _, transaction := range transactions {
			if transaction.Description == duplicated.Description && transaction.Data == duplicated.Data && transaction.Value == duplicated.Value && !transaction.Duplicated {
				duplicateds++
			}
		}
	}
	if duplicated != nil && duplicateds == 0 {
		results = createResult(err, "Duplicou errado!!", &results)
	}

	return results
}

func (a Application) request(id, year, month int) (int, string, error) {
	url := fmt.Sprintf("%s/%d/transacoes/%d/%d", a.url, id, year, month)
	response, err := http.Get(url)
	if err != nil {
		return 0, "", err
	}

	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	return response.StatusCode, string(data), err
}

func (a Application) verifyMultiRequestRules() []Result {
	var results []Result

	// two same requests must return same value
	status1, resp1, err := a.request(9505, 2020, 2)
	if err != nil {
		results = createResult(err, "Impossivel obter 9505 2020 2 1x", &results)
	}
	if status1 != 200 {
		results = createResult(fmt.Errorf("%d", status1), "Status incorreto 1x", &results)
	}
	status2, resp2, err := a.request(9505, 2020, 2)
	if err != nil {
		results = createResult(err, "Impossivel obter 9505 2020 2 2x", &results)
	}
	if status2 != 200 {
		results = createResult(fmt.Errorf("%d", status2), "Status incorreto 2x", &results)
	}
	if resp1 != resp2 {
		results = createResult(fmt.Errorf(""), "Duas requisicoes diferentes para os mesmos parmetros", &results)
	}

	// invalid ids
	status1, _, _ = a.request(100, 2020, 20)
	if status1 != 400 {
		results = createResult(fmt.Errorf("%d", status1), "Id 100 invalido nao retornou 400", &results)
	}
	status1, _, _ = a.request(10, 2020, 20)
	if status1 != 400 {
		results = createResult(fmt.Errorf("%d", status1), "Id 10 invalido nao retornou 400", &results)
	}
	status1, _, _ = a.request(15, 2020, 20)
	if status1 != 400 {
		results = createResult(fmt.Errorf("%d", status1), "Id 15 invalido nao retornou 400", &results)
	}
	status1, _, _ = a.request(1000001, 2020, 20)
	if status1 != 400 {
		results = createResult(fmt.Errorf("%d", status1), "Id 1000001 invalido nao retornou 400", &results)
	}
	// invalid years
	status1, _, _ = a.request(9505, -2, 20)
	if status1 != 400 {
		results = createResult(fmt.Errorf("%d", status1), "Ano -2 invalido nao retornou 400", &results)
	}
	// invalid months
	status1, _, _ = a.request(9505, 2020, 15)
	if status1 != 400 {
		results = createResult(fmt.Errorf("%d", status1), "Mês 15 invalido nao retornou 400", &results)
	}
	status1, _, _ = a.request(9505, 2020, 0)
	if status1 != 400 {
		results = createResult(fmt.Errorf("%d", status1), "Mês 0 invalido nao retornou 400", &results)
	}
	// 30% duplicated
	for year := 2010; year <= 2020; year++ {
		duplicated := 0.0
		for month := 1; month <= 12; month++ {
			status, data, _ := a.request(9505, year, month)
			if status != 200 {
				continue
			}
			if strings.Index(data, "true") != -1 {
				duplicated += 1.0
			}
		}
		duplicated = duplicated / 12
		if duplicated < 0.3 {
			results = createResult(fmt.Errorf("Ano %d com %f", year, duplicated), "Ano com menos de 30% duplicados", &results)
			break
		}
	}

	return results
}

func main() {
	a := Application{
		url: os.Args[2],
	}

	// tenta baixar 1 vez
	err := a.tryAccess()
	if err != nil {
		fmt.Printf("Vish, nao consegui nem acessar a URL pow!? %s\n", err)
		return
	}

	fmt.Println("Validando contrato")
	results := a.verifyContract()
	for _, result := range results {
		fmt.Printf("* %s: %s\n", result.message, result.err)
	}
	fmt.Println("Contrato validado")

	fmt.Println("Validando multiplas requests")
	results = a.verifyMultiRequestRules()
	for _, result := range results {
		fmt.Printf("* %s: %s\n", result.message, result.err)
	}
	fmt.Println("Multiplas requests validadas")
	fmt.Println("Teste validado")
}
