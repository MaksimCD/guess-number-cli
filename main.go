package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	NoHint TempHint = iota
	Hot
	Warm
	Cold
)

type TempHint int

type Result struct {
	Date     string `json:"date"`
	Win      bool   `json:"win"`
	Attempts int    `json:"attempts"`
	Max      int    `json:"max_number"`
}

func saveResult(r Result) {
	var results []Result

	data, err := os.ReadFile("results.json")
	if err == nil {
		json.Unmarshal(data, &results)
	}
	results = append(results, r)

	updatedData, _ := json.MarshalIndent(results, "", "  ")
	os.WriteFile("results.json", updatedData, 0644)
}

func choseDifficulty(reader *bufio.Reader) (int, int) {
	for {
		fmt.Println("Выберите сложность:")
		fmt.Println("1 - Easy (1-50, 15 попыток)")
		fmt.Println("2 - Medium (1-100, 10 попыток)")
		fmt.Println("3 - Hard (1-200, 5 попыток)")
		fmt.Print("Ваш выбор: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			return 50, 15
		case "2":
			return 100, 10
		case "3":
			return 200, 5
		default:
			fmt.Println("Введите 1,2 или 3")
		}

	}
}

func directionHint(secret, guess int) string {
	if guess < secret {
		return "Секретное число больше 👆"
	}
	if guess > secret {
		return "Секретное число меньше👇"
	}
	return ""
}

func distanceHint(secret, guess int) TempHint {
	diff := secret - guess
	if diff < 0 {
		diff = -diff
	}

	switch {
	case diff == 0:
		return NoHint
	case diff <= 5:
		return Hot
	case diff <= 15:
		return Warm
	default:
		return Cold
	}
}

func printTemperatureHint(h TempHint) {
	switch h {
	case Hot:
		color.Red("🔥 Горячо")
	case Warm:
		color.RGB(255, 140, 0).Println("🙂 Тепло")
	case Cold:
		color.Blue("❄️ Холодно")
	}
}

func readGuess(reader *bufio.Reader, min, max int) int {
	badNumberTries := 0
	outOfRangeTries := 0
	for {
		fmt.Printf("Введите число (%d-%d): ", min, max)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка ввода, повторите попытку!")
			continue
		}

		input = strings.TrimSpace(input)

		number, err := strconv.Atoi(input)
		if err != nil {
			badNumberTries++

			switch badNumberTries {
			case 1:
				fmt.Println("Я же сказал число!")
			case 2:
				fmt.Println("Ты серьезно? 😅")
			case 3:
				fmt.Println("Последнее предупреждение ⚠️")
			default:
				fmt.Println("Это слишком сложно для тебя 😔")
				os.Exit(0)

			}

			continue
		}
		if number < min || number > max {
			outOfRangeTries++

			switch outOfRangeTries {
			case 1:
				fmt.Printf("Я напоминаю: выбери число в диапазоне %d–%d 🙂\n", min, max)
			case 2:
				fmt.Println("Почитай, пожалуйста, правила 😄")
			case 3:
				fmt.Println("Смешно, да? 🙂")
			default:
				fmt.Println("Это слишком сложно для тебя 😔")
				os.Exit(0)
			}
			continue
		}

		outOfRangeTries = 0

		return number
	}
}

func main() {
	fmt.Println("Привет 'Угадай число' - от 1 до 100 за 10 попыток")

	reader := bufio.NewReader(os.Stdin)

	for {
		maxNumber, maxAttempts := choseDifficulty(reader)
		secretNumber := rand.Intn(maxNumber) + 1

		success := false
		var attempts []int

		for guesses := 0; guesses < maxAttempts; guesses++ {
			color.Yellow("Попыток осталось %d", maxAttempts-guesses)

			guess := readGuess(reader, 1, maxNumber)
			attempts = append(attempts, guess)
			color.Blue("Попытки %v", attempts)

			dir := directionHint(secretNumber, guess)
			if dir != "" {
				fmt.Println(dir)
			}

			if guess == secretNumber {
				success = true
				color.Green("Есть пробитие! С победой! 👑")
				break
			}

			h := distanceHint(secretNumber, guess)
			printTemperatureHint(h)
		}

		if !success {
			color.Red("Не сегодня 😔. Секретное число было: %d", secretNumber)
		}

		result := Result{
			Date:     time.Now().Format("2006-01-02 15:04:05"),
			Win:      success,
			Attempts: len(attempts),
			Max:      maxNumber,
		}
		saveResult(result)

		fmt.Print("Сыграть еще раз? (да/нет): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "да" {
			break
		}
	}
}
