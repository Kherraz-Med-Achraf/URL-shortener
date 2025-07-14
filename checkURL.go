package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"log"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func init() {
	_ = godotenv.Load()
}

func IsURLSafe(u string) bool {
	if !isFormatValid(u) {
		return false
	}
	if !openAICheck(u) {
		return false
	}
	return true
}

func isFormatValid(u string) bool {
	p, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}
	return p.Scheme == "http" || p.Scheme == "https"
}


func openAICheck(u string) bool {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY not set")
		return true // sans clé on ignore
	}

	client := openai.NewClient(apiKey)
	ctx := context.Background()

	prompt := "Analyse cette URL et réponds UNIQUEMENT par SAFE ou UNSAFE. Marque UNSAFE si l'URL contient des références à : alcool (bière, vin, spiritueux), drogues (cannabis, cocaine, etc.), contenu pornographique/sexuel, ou tout contenu interdit aux mineurs. Sois très strict dans ta détection. URL: " + u
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       "gpt-4.1-nano-2025-04-14",
		MaxTokens:   10,
		Temperature: 0,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "Tu es un filtre de contenu très strict pour mineurs. Tu détectes tout contenu lié à l'alcool, aux drogues, au sexe/pornographie. Même les marques d'alcool ou références subtiles doivent être marquées UNSAFE. Réponds UNIQUEMENT par SAFE ou UNSAFE, rien d'autre."},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil || len(resp.Choices) == 0 {
		log.Println("OpenAI request error:", err)
		return false
	}

	answer := strings.ToUpper(strings.TrimSpace(resp.Choices[0].Message.Content))
	log.Println("OpenAI answer:", answer)
	return strings.HasPrefix(answer, "SAFE")
}

func IsAliasSafe(alias string) bool {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY not set for alias check")
		return true // sans clé on ignore
	}

	client := openai.NewClient(apiKey)
	ctx := context.Background()

	prompt := "Analyse cet alias et réponds UNIQUEMENT par SAFE ou UNSAFE. Marque UNSAFE si l'alias fait référence à : alcool, drogues, contenu pornographique/sexuel, ou tout contenu interdit aux mineurs. Alias: " + alias
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       "gpt-4.1-nano-2025-04-14",
		MaxTokens:   10,
		Temperature: 0,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "Tu es un filtre de contenu très strict pour mineurs. Tu détectes tout contenu lié à l'alcool, aux drogues, au sexe/pornographie dans les alias. Réponds UNIQUEMENT par SAFE ou UNSAFE."},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil || len(resp.Choices) == 0 {
		log.Println("OpenAI alias check error:", err)
		return false
	}

	answer := strings.ToUpper(strings.TrimSpace(resp.Choices[0].Message.Content))
	log.Println("OpenAI alias answer:", answer)
	return strings.HasPrefix(answer, "SAFE")
}

func SuggestAlias(rawURL string) string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY not set for alias suggestion")
		return "link-" + uuid.NewString()[:4] // fallback
	}

	client := openai.NewClient(apiKey)
	ctx := context.Background()

	prompt := "Suggère un alias court (3-8 caractères) et approprié pour cette URL. L'alias doit être simple, mémorable et ne pas contenir de contenu inapproprié. Réponds UNIQUEMENT par l'alias suggéré, sans explication. URL: " + rawURL
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       "gpt-4.1-nano-2025-04-14",
		MaxTokens:   15,
		Temperature: 0.3,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "Tu es un générateur d'alias pour URL. Tu crées des alias courts, simples et appropriés pour tous les âges. Évite tout contenu lié à l'alcool, drogues, sexe."},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil || len(resp.Choices) == 0 {
		log.Println("OpenAI alias suggestion error:", err)
		return "link-" + uuid.NewString()[:4]
	}

	suggestion := strings.TrimSpace(resp.Choices[0].Message.Content)
	log.Println("OpenAI suggested alias:", suggestion)

	if IsAliasSafe(suggestion) {
		return suggestion
	}

	return "link-" + uuid.NewString()[:4]
}
