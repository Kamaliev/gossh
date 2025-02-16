package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func ReadCmdLine(consoleMsg string) (string, error) {
	fmt.Print(consoleMsg)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	return line, err

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование:")
		fmt.Println("  gossh server_name       - Подключиться к серверу")
		fmt.Println("  gossh -a server_name     - Добавить сервер интерактивно")
		os.Exit(1)
	}

	if os.Args[1] == "-a" {
		if len(os.Args) < 3 {
			fmt.Println("Использование: gossh -a server_name")
			os.Exit(1)
		}

		serverName := os.Args[2]

		address, _ := ReadCmdLine("[write_host]: ")
		username, _ := ReadCmdLine("[write_username]: ")
		portStr, _ := ReadCmdLine("[write_port]: ")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Println("Ошибка: порт должен быть числом")
			os.Exit(1)
		}

		err = AddServer(serverName, address, username, port)
		if err != nil {
			fmt.Println("Ошибка при сохранении сервера:", err)
			os.Exit(1)
		}

		fmt.Println("✅ Сервер успешно добавлен!")
		os.Exit(0)
	}

	serverName := os.Args[1]
	server, err := GetServer(serverName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", server.Username, server.Address), "-p", fmt.Sprintf("%d", server.Port))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Ошибка при выполнении SSH:", err)
		os.Exit(1)
	}
}
