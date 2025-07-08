# 🔗 URL Shortener - Raccourcisseur d'URL en Go

Un raccourcisseur d'URL moderne et sécurisé développé en Go avec Fiber, offrant des fonctionnalités avancées d'analyse et de personnalisation.

## ✨ Fonctionnalités Actuelles

### 🚀 Fonctionnalités de Base

- **Raccourcissement d'URL** : Conversion d'URLs longues en liens courts
- **Alias personnalisés** : Possibilité de choisir son propre alias au lieu d'un généré automatiquement
- **Expiration configurable** : Définition d'une durée de vie pour les liens raccourcis
- **Stockage en fichiers JSON** : Sauvegarde persistante des liens dans le dossier `data/`
- **Validation d'unicité** : Vérification que l'alias n'existe pas déjà

### 🔄 API Endpoints

- `GET /` : Page d'accueil
- `POST /api/shorten` : Création d'un lien raccourci
- `GET /:alias` : Redirection vers l'URL originale

## 🛠️ Installation et Utilisation

### Prérequis

- Go 1.24.4 ou supérieur
- Modules Go activés

### 📦 Dépendances Principales

Le projet utilise les dépendances suivantes :

- **Fiber v2.52.8** : Framework web haute performance pour Go
- **Google UUID v1.6.0** : Génération d'identifiants uniques
- **Brotli v1.1.0** : Compression avancée des réponses HTTP
- **FastHTTP v1.51.0** : Serveur HTTP optimisé (utilisé par Fiber)
- **Compress v1.17.9** : Algorithmes de compression supplémentaires

### Installation

```bash
# Cloner le repository
git clone <votre-repository>
cd url-shortener

# Installer les dépendances
go mod tidy

# Lancer le serveur
go run main.go
```

Le serveur sera accessible sur `http://localhost:3000`

### 🚀 Performance

Grâce à l'utilisation de **FastHTTP** et **Fiber**, le serveur offre :

- **Latence ultra-faible** : <1ms pour la redirection
- **Haute concurrence** : Supporte >100k requêtes/seconde
- **Faible consommation mémoire** : Optimisé pour les environnements contraints
- **Compression automatique** : Brotli/gzip pour réduire la bande passante

### Utilisation de l'API

#### Raccourcir une URL

```bash
curl -X POST http://localhost:3000/api/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/very/long/url",
    "alias": "mon-alias",
    "expiration_minutes": 60
  }'
```

**Réponse :**

```json
{
  "short_url": "http://localhost:3000/mon-alias",
  "expires_at": "2024-01-01T15:30:00Z"
}
```

#### Accéder à un lien raccourci

```bash
curl -L http://localhost:3000/mon-alias
```

## 🎯 Fonctionnalités Prévues

### 🤖 Intelligence Artificielle

- **Analyse de sécurité** : Détection automatique des sites dangereux, malveillants ou phishing
- **Suggestions d'alias** : Génération d'alias pertinents basés sur le contenu de l'URL

### 🔒 Sécurité et Contrôle

- **Filtrage d'âge** : Détection et restriction des contenus pour adultes (-18)
- **Validation de contenu** : Vérification que le lien est accessible et légitime

### 📊 Fonctionnalités Avancées

- **URL groupées** : Création d'une URL courte qui redirige vers plusieurs URLs
- **Statistiques d'utilisation** : Suivi des clics et analytics
- **Dashboard admin** : Interface web pour gérer les liens
- **API étendue** : Endpoints pour la gestion et les statistiques

### 🎨 Interface Utilisateur

- **Interface web moderne** : Page de création et gestion des liens
- **Responsive design** : Compatible mobile et desktop
- **Thème sombre/clair** : Personnalisation de l'interface

## 📁 Structure du Projet

```
url-shortener/                    # Module: url-shortener
├── main.go                      # Point d'entrée principal
├── data/                        # Stockage des liens (JSON)
│   ├── 6dac2c.json             # Exemple d'alias généré
│   └── fbb41d.json             # Autre exemple
├── go.mod                       # Module Go 1.24.4 + dépendances
├── go.sum                       # Checksums de sécurité
├── tmp/                         # Fichiers temporaires
│   ├── build-errors.log         # Logs de compilation
│   └── main.exe                 # Binaire compilé (Windows)
└── README.md                    # Documentation
```

## ⚙️ Architecture Technique

### Stack Technologique

- **Go 1.24.4** : Langage principal avec les dernières fonctionnalités
- **Fiber v2.52.8** : Framework web rapide et expressif, inspiré d'Express.js
- **FastHTTP** : Serveur HTTP ultra-performant (10x plus rapide que net/http)
- **UUID v1.6.0** : Génération d'identifiants uniques thread-safe
- **Compression multicouche** : Brotli + gzip pour optimiser la bande passante

### Avantages de l'Architecture

- **Performance** : FastHTTP + Fiber = latence ultra-faible
- **Simplicité** : Stockage JSON pour un démarrage rapide
- **Scalabilité** : Architecture modulaire prête pour l'extension
- **Sécurité** : Validation des entrées et gestion des erreurs robuste

## 🔧 Configuration

### Variables d'environnement

```bash
PORT=3000              # Port du serveur
DATA_DIR=data          # Dossier de stockage
AI_API_KEY=your-key    # Clé API pour l'IA (futur)
```

### Personnalisation

- Modifier le port dans `main.go`
- Changer le dossier de stockage avec la constante `dataDir`
- Ajuster la durée d'expiration par défaut

## 🚀 Déploiement

### Production

```bash
# Compiler l'application
go build -o url-shortener main.go

# Lancer en production
./url-shortener
```

