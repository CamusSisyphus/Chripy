package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/CamusSisyphus/Chripy/internal/auth"
	"github.com/CamusSisyphus/Chripy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerpolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't retrive Polka api key", err)
		return
	}
	fmt.Println(apiKey)
	fmt.Println(cfg.polkaKey)
	if apiKey != cfg.polkaKey {
		err = errors.New("Polka api don't match")
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}
	user_id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse User ID", err)
		return
	}
	_, err = cfg.db.UpgradeRedUser(r.Context(), database.UpgradeRedUserParams{
		IsChirpyRed: true,
		ID:          user_id,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find User ID", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)

}
