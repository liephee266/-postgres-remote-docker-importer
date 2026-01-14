
# PostgreSQL Remote Docker Importer

Un outil en **Go** pour automatiser l’export d’une base PostgreSQL locale et son import dans un conteneur Docker distant via SSH/SCP.

---

## Fonctionnalités

* Export de votre base PostgreSQL locale (`pg_dump`).
* Connexion SSH sécurisée au serveur distant.
* Transfert du dump SQL via SCP.
* Import automatique dans un conteneur Docker distant.
* Supporte l’authentification par mot de passe ou clé SSH.

---

## Prérequis

* Go 1.20+
* PostgreSQL installé localement (`pg_dump`)
* Accès SSH au serveur distant
* Docker installé sur le serveur distant
* Fichier `.env` pour stocker les variables d’environnement

---

## Installation

1. Clonez le dépôt :

```bash
git clone https://github.com/votre-utilisateur/postgres-remote-docker-importer.git
cd postgres-remote-docker-importer
```

2. Installez les dépendances Go :

```bash
go mod tidy
```

3. Créez un fichier `.env` à la racine du projet avec les variables suivantes :

```dotenv
# Local PostgreSQL
LOCAL_DB_USER=utilisateur_local
LOCAL_DB_NAME=nom_base_locale
LOCAL_DB_PASSWORD=motdepasse_local

# SSH
SSH_USER=utilisateur_ssh
SSH_HOST=serveur.exemple.com
SSH_PASSWORD=motdepasse_ssh  # optionnel si clé SSH
SSH_KEY=/chemin/vers/cle_privee # optionnel si mot de passe

# Remote Docker & PostgreSQL
REMOTE_DOCKER=nom_du_conteneur_docker
REMOTE_DB_USER=utilisateur_remote
REMOTE_DB_NAME=nom_base_remote
REMOTE_DB_PASSWORD=motdepasse_remote
REMOTE_DUMP=/chemin/remote/dump.sql
```

---

## Utilisation

1. Compilez et lancez le script :

```bash
go run main.go
```

2. Le script effectue automatiquement :

   * Dump de la base locale (`dump.sql`).
   * Connexion SSH au serveur distant.
   * Transfert du fichier dump via SCP.
   * Import du dump dans le conteneur Docker distant.

3. À la fin, vous devriez voir :

```
Succès ! Base importée dans le conteneur Docker distant.
```

---

## Sécurité

* La clé SSH ou le mot de passe est nécessaire pour se connecter au serveur distant.
* Les mots de passe de base de données sont transmis via des variables d’environnement et ne sont **pas** logués.

---

## Limitations

* Fonctionne uniquement avec PostgreSQL et Docker.
* Le conteneur Docker doit déjà exister sur le serveur distant.
* L’authentification SSH doit être configurée correctement (mot de passe ou clé privée).

---

## Contribution

Les contributions sont les bienvenues ! N’hésitez pas à proposer des améliorations ou signaler des bugs.

---

## Licence

MIT © [Orphée Lié]

---

Veux‑tu que je fasse ça ?
