package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
	Timer   int
}

type GameState struct {
	Name      string
	Points    int
	Questions []Question
}

func (g *GameState) Init() {
	fmt.Println("Seja bem-vindo(a) ao quiz!")
	fmt.Println("Escreva seu nome:")
	reader := bufio.NewReader(os.Stdin)

	name, err := reader.ReadString('\n')

	if err != nil {
		panic("Erro ao ler a string")
	}

	g.Name = name

	fmt.Printf("Vamos ao jogo, %s\n", g.Name)
}

func (g *GameState) ProcessCSV() {
	f, err := os.Open("quizgo.csv")
	if err != nil {
		panic("erro ao ler arquivo")
	}

	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		panic("Erro ao ler CSV")
	}

	for index, record := range records {
		if index > 0 {
			correctAnswer, _ := toInt(record[5])
			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer:  correctAnswer,
				Timer: 15,
			}

			g.Questions = append(g.Questions, question)
		}
	}
}

func (g *GameState) Run() {
	//Exibir a pergunta pro usuário
	for index, question := range g.Questions {
		fmt.Printf("\033[33m %d. %s \033[0m\n", index+1, question.Text)

		//Iterar sobre as opções do game state e exibir no terminal

		for j, option := range question.Options {
			fmt.Printf("[%d] %s\n", j+1, option)
		}

		fmt.Printf("Informe a alternativa correta (Você tem até %d segundos):\n", question.Timer)

		answerChan := make(chan string)

		go func () {
			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\n')
			answerChan <- strings.TrimSpace(read)
		}()

		var answer int
		var err error
		var timeout bool

		select {
		case userAnswer := <-answerChan:
			answer, err = toInt(userAnswer)
			if err != nil || answer < 1 || answer > 4 {
				fmt.Println("Resposta inválida. O jogo será encerrado.")
				os.Exit(1)
			}
		case <-time.After(time.Duration(question.Timer) * time.Second):
			fmt.Println("\nTempo esgotado! O jogo será encerrado.")
			timeout = true
		}

		if timeout {
			return
		}

		if answer == question.Answer {
			fmt.Println("Parabéns, você acertou!")
			g.Points += 10
		} else {
			fmt.Println("Ops! Resposta errada!")
		}
		fmt.Println("------------------------------------------")
	}
}

func main() {
	game := &GameState{Points: 0}
	go game.ProcessCSV()
	game.Init()
	game.Run()

	fmt.Printf("Fim de jogo. Você fez %d pontos\n", game.Points)
}

func toInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("não é permitido caractere diferente de número. Por favor, insira um número")
	}

	return i, nil
}
