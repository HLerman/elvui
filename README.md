# ElvUI Updater

Un utilitaire simple écrit en Go pour mettre à jour automatiquement l'addon ElvUI pour World of Warcraft.

## Fonctionnalités

- Détection automatique de la version d'ElvUI installée
- Vérification des mises à jour disponibles via l'API TukUI
- Mise à jour automatique d'ElvUI si une nouvelle version est disponible
- Journalisation des opérations dans un fichier log

## Prérequis

- Go 1.23.5 ou supérieur
- World of Warcraft installé

## Installation

1. Clonez le dépôt :
   ```bash
   git clone https://github.com/votre-nom-de-utilisateur/elvui-updater.git
   cd elvui-updater
   ```

2. Compilez le programme :
   ```bash
   just build
   ```

3. Exécutez le programme :
   ```bash
   just run
   ```

## Fonctionnement

1. Le programme vérifie la présence d'ElvUI dans le dossier des addons
2. Il compare la version installée avec la dernière version disponible
3. Si une mise à jour est nécessaire :
   - Suppression des dossiers ElvUI existants
   - Téléchargement de la nouvelle version
   - Extraction automatique dans le dossier des addons

## Logs

Les opérations sont enregistrées dans le fichier `log_elvui.log` pour le suivi et le débogage.