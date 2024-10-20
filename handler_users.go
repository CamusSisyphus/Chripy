package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/CamusSisyphus/Chripy/internal/auth"

	"github.com/CamusSisyphus/Chripy/internal/database"
)

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash provided password", err)
		return
	}
	dbuser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPW,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user in Database", err)
		return
	}

	respondWithJSON(w, 201, response{
		User{
			ID:        dbuser.ID,
			CreatedAt: dbuser.CreatedAt,
			UpdatedAt: dbuser.UpdatedAt,
			Email:     dbuser.Email,
		},
	})
}

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	dbuser, err := cfg.db.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get User by email", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, dbuser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Provided password doesn't match", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		dbuser.ID,
		cfg.jwtSecret,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate refresh Token string", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(
		context.Background(),
		database.CreateRefreshTokenParams{
			Token:  refresh_token,
			UserID: dbuser.ID,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh Token in Database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          dbuser.ID,
			CreatedAt:   dbuser.CreatedAt,
			UpdatedAt:   dbuser.UpdatedAt,
			Email:       dbuser.Email,
			IsChirpyRed: dbuser.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refresh_token,
	})
}

func (cfg *apiConfig) handlerUpdateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can not get token", err)
		return
	}

	userIDbyToken, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can't retrive user per token", err)
		return
	}

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash provided password", err)
		return
	}
	dbUser, err := cfg.db.UpdateUser(r.Context(),
		database.UpdateUserParams{
			HashedPassword: hashedPW,
			Email:          params.Email,
			ID:             userIDbyToken,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user in password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt,
			UpdatedAt:   dbUser.UpdatedAt,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
	})
}
