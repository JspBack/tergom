package main

import (
	"fmt"
	"os"
	"tergom/game"
	"tergom/utils"
	"time"

	"github.com/eiannone/keyboard"
)

func main() {
	utils.ClearScreen()
	width, height := utils.GetTerminalSize()

	selected := 0
	options := []string{"Start Game", "Scoreboard", "Exit"}

	for {
		renderMenu(options, selected, width, height)
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if key == keyboard.KeyArrowUp {
			if selected > 0 {
				selected--
			}
		} else if key == keyboard.KeyArrowDown {
			if selected < len(options)-1 {
				selected++
			}
		} else if key == keyboard.KeyEnter || char == '\r' {
			switch selected {
			case 0:
				startGame(width, height)
				utils.ClearScreen()
			case 1:
				scoreboardRender()
			case 2:
				exitGame()
			}
		}

		time.Sleep(time.Millisecond * 33)
	}
}

func renderMenu(options []string, selected int, width, height int) {
	utils.ClearScreen()
	menuHeight := len(options) * 2
	startY := (height / 2) - (menuHeight / 2)
	startX := (width / 2) - 10

	for i, option := range options {
		if i == selected {
			fmt.Printf("\033[%d;%dH> %s", startY+i*2, startX, option)
		} else {
			fmt.Printf("\033[%d;%dH  %s", startY+i*2, startX, option)
		}
	}
}

func startGame(width, height int) {
	gameInstance := game.NewGame(width, height)
	gameInstance.Start()
}

func scoreboardRender() {
	utils.ClearScreen()

	sb := struct {
		Entries []struct {
			Time  time.Time
			Score int
			Waves int
		}
	}{}

	fmt.Println("=== Scoreboard ===")
	fmt.Println("-------------------")
	if len(sb.Entries) == 0 {
		fmt.Println("No scores yet.")
	} else {
		for i, entry := range sb.Entries {
			fmt.Printf("%d. Time: %s | Score: %d | Waves: %d\n",
				i+1,
				entry.Time.Format("2006-01-02 15:04:05"),
				entry.Score,
				entry.Waves)
		}
	}
	fmt.Println("-------------------")
	fmt.Println("Press Enter to return to the menu.")

	for {
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			continue
		}
		if key == keyboard.KeyEnter || char == '\r' {
			break
		}
	}
}

func exitGame() {
	utils.ClearScreen()
	fmt.Println("Thank you for playing!")
	os.Exit(0)
}
