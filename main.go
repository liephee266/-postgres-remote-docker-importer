package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/joho/godotenv"

	"golang.org/x/crypto/ssh"
)

// Read environment variables
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("❌ Variable d'environnement manquante : %s", key)
	}
	return value
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erreur chargement .env : %v", err)
	}

	// Load environment variables
	localUser := getEnv("LOCAL_DB_USER")
	localDB := getEnv("LOCAL_DB_NAME")
	localPass := getEnv("LOCAL_DB_PASSWORD")

	sshUser := getEnv("SSH_USER")
	sshHost := getEnv("SSH_HOST")
	sshPass := os.Getenv("SSH_PASSWORD") // optional
	sshKey := os.Getenv("SSH_KEY")

	remoteDump := getEnv("REMOTE_DUMP")
	remoteDocker := getEnv("REMOTE_DOCKER")
	remoteDBPass := getEnv("REMOTE_DB_PASSWORD")
	remoteDBUser := getEnv("REMOTE_DB_USER")
	remoteDBName := getEnv("REMOTE_DB_NAME")

	localDumpFile := "dump.sql"

	// 1️⃣ Dump local database
	fmt.Println("➡️  Export local PostgreSQL…")
	if err := dumpLocalDB(localUser, localDB, localPass, localDumpFile); err != nil {
		log.Fatalf("❌ Dump local échoué : %v", err)
	}

	// 2️⃣ SSH connect
	fmt.Println("➡️  Connexion SSH…")
	client, err := sshConnect(sshUser, sshHost, sshPass, sshKey)
	if err != nil {
		log.Fatalf("❌ Connexion SSH échouée : %v", err)
	}
	defer client.Close()

	// 3️⃣ Transfer dump via SCP
	fmt.Println("➡️  Envoi du dump via SSH/SCP…")
	if err := scpFile(client, localDumpFile, remoteDump); err != nil {
		log.Fatalf("❌ Erreur SCP : %v", err)
	}

	// 4️⃣ Import inside Docker
	fmt.Println("➡️  Import dans Docker distant…")
	cmd := fmt.Sprintf(
		"docker exec -e PGPASSWORD=%s -i %s psql -U %s -d %s < %s",
		remoteDBPass, remoteDocker, remoteDBUser, remoteDBName, remoteDump,
	)

	if err := runRemoteSSH(client, cmd); err != nil {
		log.Fatalf("❌ Import Docker échoué : %v", err)
	}

	fmt.Println("🎉 Succès ! Base importée dans le conteneur Docker distant.")
}

// ------------------
// Dump local database
// ------------------
func dumpLocalDB(user, db, pass, output string) error {

	cmd := exec.Command(
		"pg_dump",
		"-U", user,
		"-d", db,
		"-f", output,
	)

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PGPASSWORD=%s", pass),
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_dump erreur : %s (%v)", string(out), err)
	}

	return nil
}

// ------------------
// SSH connection
// ------------------
func sshConnect(user, host, password, keyPath string) (*ssh.Client, error) {
	var auth []ssh.AuthMethod

	if keyPath != "" {
		key, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	if password != "" {
		auth = append(auth, ssh.Password(password))
	}

	if len(auth) == 0 {
		return nil, fmt.Errorf("aucune méthode d'authentification SSH fournie")
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return ssh.Dial("tcp", host, config)
}

// ------------------
// SCP file transfer
// ------------------
func scpFile(client *ssh.Client, localPath, remotePath string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		info, err := os.Stat(localPath)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, "C0644 %d %s\n", info.Size(), info.Name())

		f, _ := os.Open(localPath)
		defer f.Close()
		io.Copy(w, f)

		fmt.Fprint(w, "\x00")
	}()

	return session.Run(fmt.Sprintf("scp -t %s", remotePath))
}

// ------------------
// Run remote SSH command
// ------------------
func runRemoteSSH(client *ssh.Client, cmd string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var errBuf bytes.Buffer
	session.Stderr = &errBuf

	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("%s (%v)", errBuf.String(), err)
	}
	return nil
}
