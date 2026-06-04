package analyzer

import (
	"context"
	"fmt"
)

var knownRivals = map[string][]string{
	"gin-gonic/gin":      {"labstack/echo", "gofiber/fiber"},
	"labstack/echo":      {"gin-gonic/gin", "gofiber/fiber"},
	"gofiber/fiber":      {"gin-gonic/gin", "labstack/echo"},
	"facebook/react":     {"vuejs/vue", "sveltejs/svelte"},
	"vuejs/vue":          {"facebook/react", "sveltejs/svelte"},
	"sveltejs/svelte":    {"facebook/react", "vuejs/vue"},
	"expressjs/express":  {"fastify/fastify", "koajs/koa"},
	"fastify/fastify":    {"expressjs/express", "koajs/koa"},
	"koajs/koa":          {"expressjs/express", "fastify/fastify"},
	"django/django":      {"pallets/flask", "tiangolo/fastapi"},
	"pallets/flask":      {"django/django", "tiangolo/fastapi"},
	"tiangolo/fastapi":   {"django/django", "pallets/flask"},
	"rails/rails":        {"sinatra/sinatra", "hanami/hanami"},
	"nestjs/nest":        {"expressjs/express", "fastify/fastify"},
	"laravel/laravel":    {"symfony/symfony", "slimphp/slim"},
	"spring-projects/spring-boot": {"micronaut-projects/micronaut-core", "quarkusio/quarkus"},
}

func (a *Analyzer) Rivals(ctx context.Context, owner, repo string) ([]RivalSuggestion, error) {
	key := fmt.Sprintf("%s/%s", owner, repo)

	if rivals, ok := knownRivals[key]; ok {
		out := make([]RivalSuggestion, 0, len(rivals))
		for _, r := range rivals {
			out = append(out, RivalSuggestion{
				FullName:   r,
				Reason:     "Known alternative in the same ecosystem",
				Similarity: "direct-rival",
			})
		}
		return out, nil
	}

	// Dynamic fallback: suggest based on primary language
	data, err := a.client.FetchAll(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	var primaryLang string
	var maxBytes int
	for lang, bytes := range data.Languages {
		if bytes > maxBytes {
			maxBytes = bytes
			primaryLang = lang
		}
	}

	if primaryLang == "" {
		return []RivalSuggestion{}, nil
	}

	query := fmt.Sprintf("language:%s", primaryLang)
	if len(data.Repo.Topics) > 0 {
		query += fmt.Sprintf(" topic:%s", data.Repo.Topics[0])
	}

	return []RivalSuggestion{{
		FullName:   "",
		Reason:     fmt.Sprintf("Search for similar repos: %s", query),
		Similarity: "same-language",
	}}, nil
}
