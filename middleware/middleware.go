package middleware

import (
	"context"
	"corona/helpers"
	"corona/models"
	"database/sql"
	"encoding/base64"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		tokenKey := r.Header.Get("Token")

		token, err := tokenValidation(ctx, tokenKey)

		if err != nil {
			helpers.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, "user_id", token.UserId.String())
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		tokenKey := r.Header.Get("Token")

		token, err := rateLimitValidation(ctx, tokenKey)

		if err != nil {
			helpers.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, "user_id", token.UserId.String())
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}

func tokenValidation(ctx context.Context, token string) (models.TokenModel, error) {

	encToken := base64.StdEncoding.EncodeToString([]byte(token))

	tokenData, err := models.GetOneTokenByTokenKey(ctx, dbPool, encToken)
	if err != nil {

		if err == sql.ErrNoRows {
			return models.TokenModel{}, errors.Wrap(err, "Unauthorized")
		}
		return models.TokenModel{}, err
	}

	if !tokenData.IsActive {
		return models.TokenModel{}, errors.Wrap(err, "Token Inactive")
	}

	if tokenData.ExpiredAt.Sub(time.Now()).Seconds() < 0 {
		return models.TokenModel{}, errors.Wrap(err, "Token Expired")
	}

	return tokenData, nil
}

func rateLimitValidation(ctx context.Context, token string) (models.TokenModel, error) {

	encToken := base64.StdEncoding.EncodeToString([]byte(token))

	tokenData, err := models.GetOneTokenByTokenKey(ctx, dbPool, encToken)
	if err != nil {

		if err == sql.ErrNoRows {
			return models.TokenModel{}, errors.Wrap(err, "Unauthorized")
		}
		return models.TokenModel{}, err
	}

	user, err := models.GetOneUser(ctx, dbPool, tokenData.UserId)

	if err != nil {
		return models.TokenModel{}, errors.Wrap(err, "No user")
	}

	subscription, err := models.GetOneSubscription(ctx, dbPool, user.SubscriptionId)

	if err != nil {
		return models.TokenModel{}, errors.Wrap(err, "No subscription")
	}

	rateLimit, err := models.GetOneRateLimitByUserId(ctx, dbPool, tokenData.UserId)

	if err != nil {
		return models.TokenModel{}, errors.Wrap(err, "No rate limit")
	}

	if rateLimit.TotalRequest > subscription.RequestPerDay {
		return models.TokenModel{}, errors.New("Rate Limit Exceeded")
	}

	rateLimit.TotalRequest++

	rateLimit.UpdatedBy = uuid.NullUUID{
		UUID:  uuid.NewV4(),
		Valid: true,
	}

	err = rateLimit.Update(ctx, dbPool)

	if err != nil {
		return models.TokenModel{}, errors.Wrap(err, "Update Failed")
	}

	return tokenData, nil

}
