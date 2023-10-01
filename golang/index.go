package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type TeamData struct {
	Posicao          string   `json:"Posição"`
	Variacao         string   `json:"Variação"`
	Time             Team     `json:"Time"`
	PTS              string   `json:"PTS"`
	Jogos            string   `json:"Jogos"`
	Vitorias         string   `json:"Vitórias"`
	Empates          string   `json:"Empates"`
	Derrotas         string   `json:"Derrotas"`
	GolsPro          string   `json:"Gols Pró"`
	GolsContra       string   `json:"Gols Contra"`
	SaldoGols        string   `json:"Saldo de Gols"`
	CartoesAmarelos  string   `json:"Cartões Amarelo"`
	CartoesVermelhos string   `json:"Cartões Vermelho"`
	Aproveitamento   string   `json:"Aproveitamento"`
	Recentes         []string `json:"Recentes"`
	ProximoJogo      NextGame `json:"Próximo Jogo"`
}

type Team struct {
	Escudo string `json:"Escudo"`
	Nome   string `json:"Nome"`
}

type NextGame struct {
	Time   string `json:"time"`
	Escudo string `json:"escudo"`
}

func main() {
	url := "https://www.cbf.com.br/futebol-brasileiro/competicoes/campeonato-brasileiro-serie-a/2023"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("Erro na requisição HTTP:", response.Status)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Carregar o documento HTML com goquery
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Encontrar linhas da tabela com a classe "expand-trigger"
	var teamDataList []TeamData

	doc.Find("tr.expand-trigger").Each(func(i int, row *goquery.Selection) {
		teamData := TeamData{}

		tds := row.Find("td")
		ths := row.Find("th")

		if tds.Length() > 0 {
			teamData.Posicao = tds.Eq(0).Find("b").Text()
			teamData.Variacao = tds.Eq(0).Find("span").First().Text() // Pegue apenas o primeiro span

			shieldElement := tds.Eq(0).Find("img")
			if src, exists := shieldElement.Attr("src"); exists {
				teamData.Time.Escudo = src
			}

			teamData.Time.Nome = tds.Eq(0).Find("span").Eq(1).Text()
			teamData.PTS = ths.Eq(0).Text()
			teamData.Jogos = tds.Eq(1).Text()
			teamData.Vitorias = tds.Eq(2).Text()
			teamData.Empates = tds.Eq(3).Text()
			teamData.Derrotas = tds.Eq(4).Text()
			teamData.GolsPro = tds.Eq(5).Text()
			teamData.GolsContra = tds.Eq(6).Text()
			teamData.SaldoGols = tds.Eq(7).Text()
			teamData.CartoesAmarelos = tds.Eq(8).Text()
			teamData.CartoesVermelhos = tds.Eq(9).Text()
			teamData.Aproveitamento = tds.Eq(10).Text()

			var recentes []string
			tds.Eq(11).Find("span").Each(func(j int, recentElement *goquery.Selection) {
				recentes = append(recentes, recentElement.Text())
			})
			teamData.Recentes = recentes

			nextGameElement := tds.Eq(12).Find("img")
			if alt, exists := nextGameElement.Attr("alt"); exists {
				teamData.ProximoJogo.Time = alt
			}
			if src, exists := nextGameElement.Attr("src"); exists {
				teamData.ProximoJogo.Escudo = src
			}

			teamDataList = append(teamDataList, teamData)
		}
	})

	// Exibir o resultado
	for _, teamData := range teamDataList {
		fmt.Printf("Posição: %s\n", teamData.Posicao)
		fmt.Printf("Variação: %s\n", teamData.Variacao)
		fmt.Printf("Time:\n")
		fmt.Printf("  Nome: %s\n", teamData.Time.Nome)
		fmt.Printf("  Escudo: %s\n", teamData.Time.Escudo)
		fmt.Printf("PTS: %s\n", teamData.PTS)
		fmt.Printf("Jogos: %s\n", teamData.Jogos)
		fmt.Printf("Vitórias: %s\n", teamData.Vitorias)
		fmt.Printf("Empates: %s\n", teamData.Empates)
		fmt.Printf("Derrotas: %s\n", teamData.Derrotas)
		fmt.Printf("Gols Pró: %s\n", teamData.GolsPro)
		fmt.Printf("Gols Contra: %s\n", teamData.GolsContra)
		fmt.Printf("Saldo de Gols: %s\n", teamData.SaldoGols)
		fmt.Printf("Cartões Amarelo: %s\n", teamData.CartoesAmarelos)
		fmt.Printf("Cartões Vermelho: %s\n", teamData.CartoesVermelhos)
		fmt.Printf("Aproveitamento: %s\n", teamData.Aproveitamento)
		fmt.Printf("Recentes: %s\n", strings.Join(teamData.Recentes, ", "))
		fmt.Printf("Próximo Jogo:\n")
		fmt.Printf("  Time: %s\n", teamData.ProximoJogo.Time)
		fmt.Printf("  Escudo: %s\n", teamData.ProximoJogo.Escudo)
		fmt.Println("------------------------")
	}

	// Salvar o resultado em um arquivo JSON
	jsonData, err := json.MarshalIndent(teamDataList, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile("dados.json", jsonData, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
