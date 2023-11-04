package main

import (
	"bufio"
	"fmt"
	"os"

	"instagram-lite/cli"
)

func promptInput(scanner *bufio.Scanner, text string) string {
	fmt.Print(text)
	scanner.Scan()
	return scanner.Text()
}

func outputHandler(err error, a ...interface{}) {
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println()
		return
	}
	fmt.Println(a...)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	exit := false
	menu := "Activity Reporter\n\n" +
		"1. Setup\n" +
		"2. Action\n" +
		"3. Display\n" +
		"4. Trending\n" +
		"5. Exit"

	for !exit {
		fmt.Println(menu)
		input := promptInput(scanner, "Enter menu: ")

		switch input {
		case "1":
			relation := promptInput(scanner, "Setup social graph: ")
			_, _, _, err := cli.HandleSetup(relation)
			outputHandler(err)
		case "2":
			action := promptInput(scanner, "Enter user Actions: ")
			_, _, err := cli.HandleAction(action)
			outputHandler(err)
		case "3":
			display := promptInput(scanner, "Display activity for: ")
			res, err := cli.HandleDisplay(display)
			outputHandler(err, res)
		case "4":
			res := cli.HandleTrending()
			outputHandler(nil, res)
		case "5":
			exit = true
			fmt.Println("")
			fmt.Println("Good bye!")
		default:
			fmt.Println("")
			fmt.Println("invalid menu")
		}
	}
}
