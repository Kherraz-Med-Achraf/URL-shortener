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
	"golang.org/x/crypto/bcrypt"
)

type Link struct {
	Alias      string    `json:"alias"`
	TargetURL  string    `json:"target_url"`
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
			URL               string `json:"url"`
			Alias             string `json:"alias"`
			ExpirationMinutes int    `json:"expiration_minutes"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "JSON invalide")
		}
		if body.URL == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Champ url manquant")
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
		}

		expiresAt := time.Time{} // 0 = jamais expiré
		if body.ExpirationMinutes > 0 {
			expiresAt = time.Now().Add(time.Duration(body.ExpirationMinutes) * time.Minute)
		}

		link := Link{
			Alias: alias, TargetURL: body.URL,
			Owner:     owner,
			CreatedAt: time.Now(),
			ExpiresAt: expiresAt,
		}
		if err := saveLink(link); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"short_url":  "http://localhost:3000/" + alias,
			"expires_at": expiresAt,
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

	app.Listen(":3000")
}
