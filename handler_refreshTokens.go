package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/CamusSisyphus/Chripy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	req_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	refresh_token, err := cfg.db.GetToken(context.Background(), req_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	if refresh_token.ExpiresAt.Before(time.Now()) || refresh_token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired/revoked!", errors.New("Refresh token expired!/revoked"))
		return
	}

	accessToken, err := auth.MakeJWT(refresh_token.UserID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {

	req_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	_, err = cfg.db.GetToken(context.Background(), req_token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Refresh token doesn't exist!", err)
		return
	}
	err = cfg.db.RevokeToken(context.Background(), req_token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke token!", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
