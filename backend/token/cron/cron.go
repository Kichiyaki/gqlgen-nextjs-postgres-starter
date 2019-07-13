package cron

import (
	"context"
	"time"

	"github.com/kichiyaki/graphql-starter/backend/models"
	pgfilter "github.com/kichiyaki/pg-filter"

	"github.com/kichiyaki/graphql-starter/backend/token"
	"github.com/robfig/cron"
)

type tokenCronHandler struct {
	tokenRepo token.Repository
}

func InitTokenCron(c *cron.Cron, repo token.Repository) {
	handler := &tokenCronHandler{
		repo,
	}

	c.AddFunc("@every 1m", handler.checkObsoleteActivationTokens)
}

func (handler *tokenCronHandler) checkObsoleteActivationTokens() {
	filter := &models.TokenFilter{
		Type:      models.AccountActivationTokenType,
		CreatedAt: "lt__" + time.Now().Add(-1*time.Hour).Format("2006-01-02 15:04:05"),
	}
	ctx := context.Background()
	tokens, err := handler.tokenRepo.Fetch(ctx, pgfilter.New(filter.ToMap()))
	if len(tokens) > 0 && err == nil {
		ids := []int{}
		for _, token := range tokens {
			ids = append(ids, token.ID)
		}
		_, err = handler.tokenRepo.Delete(ctx, ids)
	}
}
