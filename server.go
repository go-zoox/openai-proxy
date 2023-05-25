package main

import (
	"net/http"
	"strings"

	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/proxy/utils/rewriter"

	"github.com/go-zoox/proxy"
	"github.com/go-zoox/zoox/defaults"
	"github.com/go-zoox/zoox/middleware"
)

const OpenAIServer = "https://api.openai.com"

type Config struct {
	Port      int64
	BasePath  string
	AuthToken string
	APIKey    string
}

type Models struct {
	ChatCompletions Model
	Embeddings      Model
}

type Model struct {
	Resource   string
	Deployment string
}

func Server(cfg *Config) error {
	app := defaults.Application()

	fmt.PrintJSON(cfg)

	if cfg.AuthToken != "" {
		app.Use(middleware.BearerToken(strings.Split(cfg.AuthToken, ",")))
	}

	app.Proxy(cfg.BasePath, OpenAIServer, func(c *proxy.SingleTargetConfig) {
		c.RequestHeaders = http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", cfg.APIKey)},
			"User-Agent":    []string{fmt.Sprintf("GoZoox/OpenAI-Proxy@%s", Version)},
		}

		c.Rewrites = rewriter.Rewriters{
			{
				From: fmt.Sprintf("^%s/(.*)$", cfg.BasePath),
				To:   "/$1",
			},
		}
	})

	return app.Run(fmt.Sprintf(":%d", cfg.Port))
}
