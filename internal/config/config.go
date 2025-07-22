package config

import (
	"os"
	"strconv"
	"test-assignment/internal/sshclient"

	"github.com/joho/godotenv"
)

func LoadSSHConfig() sshclient.SSHConfig {
	_ = godotenv.Load(".env")
	port, err := strconv.Atoi(os.Getenv("SSH_PORT"))
	if err != nil {
		port = 22
	}
	return sshclient.SSHConfig{
		Host:     os.Getenv("SSH_HOST"),
		Port:     port,
		User:     os.Getenv("SSH_USER"),
		Password: os.Getenv("SSH_PASSWORD"),
	}
}
