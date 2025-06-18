package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func ReadCmdLine(consoleMsg string, validate *func(string) (bool, error)) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(consoleMsg)
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		line = strings.TrimSpace(line)
		if validate != nil {
			isLine, _ := (*validate)(line)
			if !isLine {
				continue
			}
		}

		return line, err
	}

}

func serverList() {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}

	// Определяем максимальную длину строк
	maxAddrLen := len("Адрес")
	maxNameLen := len("Имя сервера")
	maxNumLen := len(fmt.Sprintf("%d", len(config))) // Длина номера (зависит от кол-ва серверов)

	for alias, userSSH := range config {
		if len(userSSH.Address) > maxAddrLen {
			maxAddrLen = len(userSSH.Address)
		}
		if len(alias) > maxNameLen {
			maxNameLen = len(alias)
		}
	}

	// Добавляем небольшой отступ (по 2 пробела с каждой стороны)
	padding := 2
	numWidth := maxNumLen + padding
	addrWidth := maxAddrLen + padding
	nameWidth := maxNameLen + padding

	// Формируем строки разделителей
	borderTop := fmt.Sprintf("┌%s┬%s┬%s┐",
		strings.Repeat("─", numWidth),
		strings.Repeat("─", addrWidth),
		strings.Repeat("─", nameWidth),
	)
	borderMiddle := fmt.Sprintf("├%s┼%s┼%s┤",
		strings.Repeat("─", numWidth),
		strings.Repeat("─", addrWidth),
		strings.Repeat("─", nameWidth),
	)
	borderBottom := fmt.Sprintf("└%s┴%s┴%s┘",
		strings.Repeat("─", numWidth),
		strings.Repeat("─", addrWidth),
		strings.Repeat("─", nameWidth),
	)

	// Выводим таблицу
	fmt.Println(borderTop)
	fmt.Printf("│ %-*s │ %-*s │ %-*s │\n", numWidth-padding, "№", addrWidth-padding, "Адрес", nameWidth-padding, "Имя сервера")
	fmt.Println(borderMiddle)

	// Вывод данных с нумерацией
	i := 1
	for alias, userSSH := range config {
		fmt.Printf("│ %-*d │ %-*s │ %-*s │\n", numWidth-padding, i, addrWidth-padding, userSSH.Address, nameWidth-padding, alias)
		i++
	}

	fmt.Println(borderBottom)
}

func portValidate(s string) (bool, error) {
	_, err := strconv.Atoi(s)
	if err != nil {
		return false, errors.New("invalid port")
	}
	return true, nil
}

func serverAdd(serverName string) {
	portValidator := portValidate
	address, _ := ReadCmdLine("[write_host]: ", nil)
	username, _ := ReadCmdLine("[write_username]: ", nil)
	port, _ := ReadCmdLine("[write_port]: ", &portValidator)
	userSsh := UserSSH{address, username, port}

	// Загружаем конфиг и сохраняем данные нового сервера
	config, _ := LoadConfig()
	config[serverName] = userSsh
	err := SaveConfig(config)
	if err != nil {
		fmt.Println("error saving config:", err)
		return
	}

	// Выводим подтверждение
	fmt.Printf("Сервер \"%s\" успешно добавлен\n", serverName)

	// Формируем адрес подключения
	fullAddress := fmt.Sprintf("%s@%s", username, address)

	// Если порт отличный от стандартного (22), добавим его в команду
	var sshCopyCmd *exec.Cmd
	if port != "22" {
		sshCopyCmd = exec.Command("ssh-copy-id", "-p", port, fullAddress)
	} else {
		sshCopyCmd = exec.Command("ssh-copy-id", fullAddress)
	}

	// Подключаем ввод/вывод к системной команде, чтобы пользователь мог ввести пароль
	sshCopyCmd.Stdin = os.Stdin
	sshCopyCmd.Stdout = os.Stdout
	sshCopyCmd.Stderr = os.Stderr
	// Запускаем ssh-copy-id
	fmt.Println("Попытка передачи SSH-ключа на сервер с помощью ssh-copy-id...")
	err = sshCopyCmd.Run()
	if err != nil {
		fmt.Println("Ошибка при выполнении ssh-copy-id:", err)
	} else {
		fmt.Println("SSH-ключ успешно передан на удалённый сервер.")
	}
}

func serverConnect(serverName string) {
	config, _ := LoadConfig()
	if userSSH, ok := config[serverName]; ok {
		cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", userSSH.Username, userSSH.Address), "-p", userSSH.Port)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println("Ошибка при выполнении SSH:", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Профиль: \"%s\" не найден \n")
	}
}
func printLogo() {
	// Генерация ASCII-графики
	figure.NewFigure("gossh", "slant", true).Print()
	fmt.Println()
	fmt.Println()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "gossh",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("list") {
			serverList()
			return
		} else if cmd.Flags().Changed("add") {
			if len(args) < 1 {
				fmt.Println("Нужно указать [server_name]: gossh -a server_name")
				return
			}
			serverAdd(args[0])
			return
		} else if cmd.Flags().Changed("server") {
			if len(args) < 1 {
				fmt.Println("Нужно указать [server_name]: gossh -c server_name")
				return
			}
			serverConnect(args[0])
			return
		}

		if len(args) == 0 {
			printLogo()
			_ = cmd.Help()
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("list", "l", false, "servers list")
	rootCmd.Flags().BoolP("server", "s", false, "-s server_name to connect server")
	rootCmd.Flags().BoolP("add", "a", false, "-a server_name to add server")
}
