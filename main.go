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
	options := []string{"Start Game", "Exit"}

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
				exitGame()
			}
		}

		time.Sleep(100 * time.Millisecond)
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

func exitGame() {
	utils.ClearScreen()
	fmt.Println("Thank you for playing!")
	os.Exit(0)
}
