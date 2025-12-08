package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func StartSyncLoop() {
	config := LoadConfig()

	os.MkdirAll("logs", 0777)

	for {
		last, _ := LoadLastImport()

		if time.Since(last) > 5*time.Hour {
			log.Println("⏳ Plus de 5h depuis le dernier import → lancement import")
			err := RunSync(config)
			if err != nil {
				log.Printf("❌ Import échoué : %v", err)
			} else {
				SaveLastImport()
			}
		} else {
			next := last.Add(5 * time.Hour)
			wait := time.Until(next)
			log.Printf("🟢 Pas besoin d'import. Prochain import dans %.1f minutes", wait.Minutes())
		}

		time.Sleep(10 * time.Minute)
	}
}

func RunSync(cfg Config) error {
	log.Println("➡️ Export local PostgreSQL…")
	dump := fmt.Sprintf("dump_%d.sql", time.Now().Unix())

	cmd := exec.Command("pg_dump", "-U", cfg.LocalDBUser, "-d", cfg.LocalDBName, "-f", dump)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", cfg.LocalDBPass))
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Printf("❌ pg_dump : %s", string(out))
		return err
	}

	log.Println("➡️ Transfert du dump via SCP…")
	scp := exec.Command("scp", "-i", cfg.SSHKey, dump, fmt.Sprintf("%s@%s:/tmp/%s", cfg.SSHUser, cfg.SSHHost, dump))
	if out, err := scp.CombinedOutput(); err != nil {
		return fmt.Errorf("SCP failed: %s", string(out))
	}

	log.Println("➡️ Import dans Docker distant…")
	sshCmd := fmt.Sprintf(
		"docker exec -e PGPASSWORD=%s -i %s psql -U %s -d %s < /tmp/%s",
		cfg.RemoteDBPass,
		cfg.RemoteDocker,
		cfg.RemoteDBUser,
		cfg.RemoteDBName,
		dump,
	)

	ssh := exec.Command("ssh", "-i", cfg.SSHKey, fmt.Sprintf("%s@%s", cfg.SSHUser, cfg.SSHHost), sshCmd)
	if out, err := ssh.CombinedOutput(); err != nil {
		return fmt.Errorf("Import Docker échoué : %s", string(out))
	}

	log.Println("✔️ Import terminé avec succès")
	return nil
}
