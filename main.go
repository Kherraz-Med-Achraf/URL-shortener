package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"
)

type Link struct {
	Alias      string    `json:"alias"`
	TargetURL  string    `json:"target_url,omitempty"`
	TargetURLs []string  `json:"target_urls,omitempty"`
	Multi      bool      `json:"multi"`
	Owner      string    `json:"owner"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	ClickCount int       `json:"click_count"`
}

const dataDir = "data/links"

func saveLink(l Link) error {
	path := filepath.Join(dataDir, l.Alias+".json")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(l)
}

func loadLink(alias string) (Link, error) {
	var l Link
	path := filepath.Join(dataDir, alias+".json")
	f, err := os.Open(path)
	if err != nil {
		return l, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&l)
	return l, err
}

func main() {
	// On s'assure que le dossier "data" existe
	_ = os.MkdirAll(dataDir, 0755)

	app := fiber.New()

	app.Static("/", "./static")

	// Route pour vérifier que ça marche encore
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World en Go avec Fiber !!! 🚀")
	})

	// Route GET /:alias pour rediriger
	app.Get("/:alias", func(c *fiber.Ctx) error {
		alias := c.Params("alias")
		link, err := loadLink(alias)
		if err != nil {
			// Lien inconnu
			return fiber.NewError(fiber.StatusNotFound, "Lien non trouvé")
		}
		// Vérifie l'expiration
		if !link.ExpiresAt.IsZero() && time.Now().After(link.ExpiresAt) {
			return fiber.NewError(fiber.StatusGone, "Lien expiré")
		}

		if link.Multi {
			// Incrémente clics
			link.ClickCount++
			_ = saveLink(link)
			// Génère page HTML simple
			html := "<!DOCTYPE html><html><head><meta charset='utf-8'><title>Liens</title><link rel='stylesheet' href='/styles.css'></head><body><div class='container'><h1>Choisissez un lien</h1><ul>"
			for _, u := range link.TargetURLs {
				html += "<li><a href='" + u + "' target='_blank'>" + u + "</a></li>"
			}
			html += "</ul></div></body></html>"
			return c.Type("html").SendString(html)
		}

		// Incrémente le nombre de clics
		link.ClickCount++
		_ = saveLink(link)

		// Redirige vers l'URL d'origine
		return c.Redirect(link.TargetURL, fiber.StatusSeeOther)
	})

	// PROTECTED ROUTES

	const jwtSecret = "vraiment-secret"

	jwtMW := jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
		ContextKey: "userToken",
	})

	// Tout ce qui commence par /api/private/* exige le JWT
	app.Use("/api/private", jwtMW)

	// Route POST /api/shorten
	app.Post("/api/private/shorten", func(c *fiber.Ctx) error {

		tok := c.Locals("userToken").(*jwt.Token)
		owner := tok.Claims.(jwt.MapClaims)["sub"].(string)

		var body struct {
			URL               string   `json:"url"`
			URLs              []string `json:"urls"`
			Alias             string   `json:"alias"`
			ExpirationMinutes int      `json:"expiration_minutes"`
			Multi             bool     `json:"multi"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "JSON invalide")
		}
		if body.Multi {
			if len(body.URLs) == 0 {
				return fiber.NewError(fiber.StatusBadRequest, "Aucune URL fournie")
			}
			// Vérification de sécurité pour chaque URL
			for _, u := range body.URLs {
				if !IsURLSafe(u) {
					return fiber.NewError(fiber.StatusBadRequest, "URL refusée car elle contient du contenu inapproprié (alcool, drogues, contenu adulte): "+u)
				}
			}
		} else {
			if body.URL == "" {
				return fiber.NewError(fiber.StatusBadRequest, "Champ url manquant")
			}
			// Vérification de sécurité de l'URL simple
			if !IsURLSafe(body.URL) {
				return fiber.NewError(fiber.StatusBadRequest, "URL refusée car elle contient du contenu inapproprié (alcool, drogues, contenu adulte)")
			}
		}

		// Choix de l'alias : personnalisé ou généré
		alias := body.Alias
		if alias == "" {
			alias = uuid.NewString()[:6]
		} else {
			// Vérifie que l'alias n'existe pas déjà
			if _, err := loadLink(alias); err == nil {
				return fiber.NewError(fiber.StatusConflict, "Alias déjà utilisé")
			}
			// Vérifie que l'alias ne contient pas de contenu inapproprié
			if !IsAliasSafe(alias) {
				return fiber.NewError(fiber.StatusBadRequest, "Alias refusé car il contient du contenu inapproprié")
			}
		}

		expiresAt := time.Time{} // 0 = jamais expiré
		if body.ExpirationMinutes > 0 {
			expiresAt = time.Now().Add(time.Duration(body.ExpirationMinutes) * time.Minute)
		}

		link := Link{
			Alias:     alias,
			Owner:     owner,
			CreatedAt: time.Now(),
			ExpiresAt: expiresAt,
			Multi:     body.Multi,
		}
		if body.Multi {
			link.TargetURLs = body.URLs
		} else {
			link.TargetURL = body.URL
		}
		if err := saveLink(link); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"short_url":  "http://localhost:3000/" + alias,
			"expires_at": expiresAt,
		})
	})

	// Route pour suggérer un alias basé sur une URL en utilisant l'IA
	app.Post("/api/private/suggest-alias", func(c *fiber.Ctx) error {
		var body struct {
			URL string `json:"url"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "JSON invalide")
		}
		if body.URL == "" {
			return fiber.NewError(fiber.StatusBadRequest, "URL manquante")
		}

		suggestion := SuggestAlias(body.URL)
		return c.JSON(fiber.Map{
			"suggested_alias": suggestion,
		})
	})

	app.Get("/api/private/links", func(c *fiber.Ctx) error {
		tok := c.Locals("userToken").(*jwt.Token)
		username := tok.Claims.(jwt.MapClaims)["sub"].(string)
		isAdmin := tok.Claims.(jwt.MapClaims)["adm"].(bool)

		files, _ := filepath.Glob(filepath.Join(dataDir, "*.json"))
		var list []Link
		for _, f := range files {
			var l Link
			fd, _ := os.Open(f)
			if json.NewDecoder(fd).Decode(&l) == nil && (isAdmin || l.Owner == username) {
				list = append(list, l)
			}
			fd.Close()
		}
		return c.JSON(list)
	})

	app.Delete("/api/private/links/:alias", func(c *fiber.Ctx) error {
		tok := c.Locals("userToken").(*jwt.Token)
		username := tok.Claims.(jwt.MapClaims)["sub"].(string)
		isAdmin := tok.Claims.(jwt.MapClaims)["adm"].(bool)

		alias := c.Params("alias")

		// Charger le lien pour vérifier le propriétaire
		link, err := loadLink(alias)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Lien non trouvé")
		}

		// Vérifier les permissions
		if !isAdmin && link.Owner != username {
			return fiber.NewError(fiber.StatusForbidden, "Accès refusé")
		}

		// Supprimer le fichier
		if err := os.Remove(filepath.Join(dataDir, alias+".json")); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Erreur lors de la suppression")
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	// ADMIN ROUTES

	app.Get("/api/private/admin/users", func(c *fiber.Ctx) error {
		tok := c.Locals("userToken").(*jwt.Token)
		claims := tok.Claims.(jwt.MapClaims)
		if adm, ok := claims["adm"].(bool); !ok || !adm {
			return fiber.NewError(fiber.StatusForbidden, "Accès refusé")
		}

		files, _ := filepath.Glob(filepath.Join(userDir, "*.json"))
		var users []User
		for _, f := range files {
			fd, err := os.Open(f)
			if err != nil {
				continue
			}
			var u User
			if json.NewDecoder(fd).Decode(&u) == nil {
				users = append(users, u)
			}
			fd.Close()
		}
		return c.JSON(users)
	})

	app.Delete("/api/private/admin/users/:username", func(c *fiber.Ctx) error {
		tok := c.Locals("userToken").(*jwt.Token)
		claims := tok.Claims.(jwt.MapClaims)
		if adm, ok := claims["adm"].(bool); !ok || !adm {
			return fiber.NewError(fiber.StatusForbidden, "Accès refusé")
		}

		username := c.Params("username")

		// Supprime le fichier utilisateur
		_ = os.Remove(filepath.Join(userDir, username+".json"))

		// Supprime tous les liens appartenant à cet utilisateur
		linkFiles, _ := filepath.Glob(filepath.Join(dataDir, "*.json"))
		for _, lf := range linkFiles {
			fd, err := os.Open(lf)
			if err != nil {
				continue
			}
			var l Link
			if json.NewDecoder(fd).Decode(&l) == nil && l.Owner == username {
				fd.Close()
				_ = os.Remove(lf)
			} else {
				fd.Close()
			}
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	// Inscription : POST /api/register
	app.Post("/api/register", func(c *fiber.Ctx) error {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Données invalides")
		}
		if body.Username == "" || body.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Champs obligatoires")
		}
		if _, err := loadUser(body.Username); err == nil {
			return fiber.NewError(fiber.StatusConflict, "Utilisateur déjà existant")
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		user := User{Username: body.Username, PasswordHash: string(hash)}
		if err := saveUser(user); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.SendStatus(fiber.StatusCreated)
	})

	// Connexion : POST /api/login
	app.Post("/api/login", func(c *fiber.Ctx) error {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Données invalides")
		}
		u, err := loadUser(body.Username)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(body.Password)) != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Identifiants invalides")
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": u.Username,
			"adm": u.IsAdmin,
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		})
		t, _ := token.SignedString([]byte(jwtSecret))
		return c.JSON(fiber.Map{"token": t})
	})

	app.Get("/qr/:alias", func(c *fiber.Ctx) error {
		alias := c.Params("alias")
		if alias == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Alias manquant")
		}
		if _, err := os.Stat(filepath.Join(dataDir, alias+".json")); err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Lien non trouvé")
		}
		shortURL := "http://" + c.Hostname() + "/" + alias
		png, err := qrcode.Encode(shortURL, qrcode.Medium, 512)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		c.Set("Content-Type", "image/png")
		return c.Send(png)
	})

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
