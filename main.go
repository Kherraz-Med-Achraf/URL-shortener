package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Link struct {
	Alias      string    `json:"alias"`
	TargetURL  string    `json:"target_url"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	ClickCount int       `json:"click_count"`
}

const dataDir = "data"

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

	// Route POST /api/shorten
	app.Post("/api/shorten", func(c *fiber.Ctx) error {
		var body struct {
			URL               string `json:"url"`
			Alias             string `json:"alias"` // Nouveau champ pour personnalisation
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
			Alias:     alias,
			TargetURL: body.URL,
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

	app.Get("/api/links", func(c *fiber.Ctx) error {
		files, _ := filepath.Glob(filepath.Join(dataDir, "*.json"))
		links := make([]Link, 0, len(files))
		for _, f := range files {
			fjson, err := os.Open(f)
			if err != nil {
				continue
			}
			var l Link
			if json.NewDecoder(fjson).Decode(&l) == nil {
				links = append(links, l)
			}
			fjson.Close()
		}
		return c.JSON(links)
	})

	app.Listen(":3000")
}
