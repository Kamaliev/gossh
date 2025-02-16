package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type alias = string

type UserSSH struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Port     int    `json:"port"`
}

var configFilePath = filepath.Join(os.Getenv("HOME"), ".gossh.json")

func LoadConfig() (map[alias]UserSSH, error) {
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[alias]UserSSH), nil // Если файла нет, возвращаем пустую карту
		}
		return nil, err
	}

	var storage map[alias]UserSSH
	err = json.Unmarshal(file, &storage)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func SaveConfig(storage map[alias]UserSSH) error {
	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFilePath, data, 0644)
}

func GetServer(name string) (UserSSH, error) {
	storage, err := LoadConfig()
	if err != nil {
		return UserSSH{}, err
	}

	server, exists := storage[name]
	if !exists {
		return UserSSH{}, fmt.Errorf("сервер '%s' не найден", name)
	}

	return server, nil
}

func AddServer(name, address, username string, port int) error {
	storage, err := LoadConfig()
	if err != nil {
		return err
	}

	storage[name] = UserSSH{Address: address, Username: username, Port: port}

	command := fmt.Sprintf("ssh-copy-id -i %s %s@%s", filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub"), username, address)
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("не удалось добавить ключ на сервер: %w", err)
	}

	return SaveConfig(storage)
}
