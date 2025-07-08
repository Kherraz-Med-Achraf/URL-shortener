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
	Alias     string    `json:"alias"`
	TargetURL string    `json:"target_url"`
	CreatedAt time.Time `json:"created_at"`
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

	// Route pour vérifier que ça marche encore
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World en Go avec Fiber ! 🚀")
	})

	// Route POST /api/shorten
	app.Post("/api/shorten", func(c *fiber.Ctx) error {
		var body struct {
			URL string `json:"url"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "JSON invalide")
		}
		if body.URL == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Champ url manquant")
		}

		// Génère un code court (6 caractères aléatoires)
		alias := uuid.NewString()[:6]
		link := Link{
			Alias:     alias,
			TargetURL: body.URL,
			CreatedAt: time.Now(),
		}
		if err := saveLink(link); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"short_url": "http://localhost:3000/" + alias,
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
		// Redirige vers l'URL d'origine
		return c.Redirect(link.TargetURL, fiber.StatusSeeOther)
	})

	app.Listen(":3000")
}
