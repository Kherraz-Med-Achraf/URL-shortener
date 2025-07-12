# 🔗 URL Shortener - Raccourcisseur d'URL en Go

Un raccourcisseur d'URL moderne et sécurisé développé en Go avec Fiber, offrant des fonctionnalités avancées d'analyse et de personnalisation.

## ✨ Fonctionnalités Actuelles

### 🚀 Fonctionnalités de Base

- **Raccourcissement d'URL** : Conversion d'URLs longues en liens courts (avec QR code automatique)
- **Alias IA** : Suggère un alias pertinent grâce à OpenAI (bouton « Suggérer avec IA »)
- **Filtrage IA** : Détecte les URLs/alias contenant contenu adulte, alcool, drogues, etc. (-18)
- **Multi-URLs** : Un seul lien court peut rediriger vers plusieurs destinations
- **Expiration configurable** : Définition d'une durée de vie pour les liens raccourcis
- **Authentification JWT** : Inscription / connexion, routes privées sécurisées
- **Dashboard Admin** : Gestion des liens et utilisateurs (routes protégées)
- **Stockage JSON** : Sauvegarde persistante dans `data/` (liens) et `data/users/` (utilisateurs)
- **QR code** : Génération à la volée via `/qr/:alias`

### 🔄 API Endpoints

**Public**

- `GET /` : Page d'accueil (statique)
- `GET /:alias` : Redirection ou page multi-liens
- `GET /qr/:alias` : QR code PNG pour le lien court
- `POST /api/register` : Inscription (JSON `username` / `password`)
- `POST /api/login` : Connexion → retourne un JWT

**Protégées (Header `Authorization: Bearer <token>`)**

- `POST /api/private/shorten` : Créer un lien court (simple ou multi-URL)
- `POST /api/private/suggest-alias` : Obtenir une proposition d'alias IA
- `GET  /api/private/links` : Lister ses liens (ou tous si admin)
- `DELETE /api/private/links/:alias` : Supprimer un lien
- `GET  /api/private/admin/users` : Liste des utilisateurs (admin)
- `DELETE /api/private/admin/users/:username` : Supprimer un utilisateur (admin)

## 🛠️ Installation et Utilisation

### Prérequis

- Go 1.24.4 ou supérieur
- Modules Go activés

### 📦 Dépendances Principales

Le projet utilise les dépendances suivantes :

- **Fiber v2.52.8** : Framework web haute performance pour Go
- **Google UUID v1.6.0** : Génération d'identifiants uniques
- **JWT v4.5.0** : Authentification sécurisée par JSON Web Tokens
- **Go-OpenAI v1.40.5** : Appels à l'API OpenAI pour filtrage et suggestion
- **godotenv v1.5.0** : Chargement automatique des variables d'environnement depuis `.env`
- **Go-QRCode** : Génération de QR codes pour chaque lien
- **Brotli / FastHTTP / Compress** : Performance et compression

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

## 📁 Structure du Projet

```
url-shortener/                    # Module: url-shortener
├── main.go                      # Point d'entrée principal
├── data/                        # Stockage des liens (JSON)
│   ├── links/                   # Stockage des liens raccourcis
│   └── users/                   # Données utilisateurs
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

### Dépendances Principales

- **Fiber v2.52.8** : Framework web haute performance
  - `gofiber/jwt/v3` : Middleware JWT pour l'authentification
- **OpenAI v1.40.5** : Intégration IA pour l'analyse de contenu et suggestions
- **UUID v1.6.0** : Génération d'identifiants uniques cryptographiquement sûrs
- **QR Code v0.0.0** : Génération de codes QR pour les liens raccourcis
- **JWT v4.5.0** : Tokens d'authentification sécurisés
- **GoDotEnv v1.5.1** : Gestion des variables d'environnement
- **Crypto** : Chiffrement et hachage sécurisés

### Dépendances Système

- **Brotli v1.1.0** : Compression avancée (meilleure que gzip)
- **FastHTTP v1.51.0** : Serveur HTTP ultra-rapide
- **Colorable/IsATTY** : Support couleurs terminal multiplateforme

### Avantages de l'Architecture

- **Performance** : FastHTTP + Fiber = latence ultra-faible
- **Simplicité** : Stockage JSON pour un démarrage rapide
- **Scalabilité** : Architecture modulaire prête pour l'extension
- **Sécurité** : Validation des entrées et gestion des erreurs robuste

## 🔧 Configuration

### Variables d'environnement

```bash
OPENAI_API_KEY=sk-...  # Clé API OpenAI (requise pour l'IA)
```
