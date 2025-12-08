package main

import (
	"log"
	"os"
)

type Config struct {
	LocalDBUser     string
	LocalDBName     string
	LocalDBPass     string
	SSHUser         string
	SSHHost         string
	SSHPassword     string
	SSHKey          string
	RemoteDump      string
	RemoteDocker    string
	RemoteDBUser    string
	RemoteDBPass    string
	RemoteDBName    string
}

func LoadConfig() Config {
	return Config{
		LocalDBUser:  mustEnv("LOCAL_DB_USER"),
		LocalDBName:  mustEnv("LOCAL_DB_NAME"),
		LocalDBPass:  mustEnv("LOCAL_DB_PASSWORD"),
		SSHUser:      mustEnv("SSH_USER"),
		SSHHost:      mustEnv("SSH_HOST"),
		SSHPassword:  os.Getenv("SSH_PASSWORD"),
		SSHKey:       mustEnv("SSH_KEY"),
		RemoteDump:   mustEnv("REMOTE_DUMP"),
		RemoteDocker: mustEnv("REMOTE_DOCKER"),
		RemoteDBUser: mustEnv("REMOTE_DB_USER"),
		RemoteDBPass: mustEnv("REMOTE_DB_PASSWORD"),
		RemoteDBName: mustEnv("REMOTE_DB_NAME"),
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("❌ Variable d'environnement manquante : %s", key)
	}
	return v
}
