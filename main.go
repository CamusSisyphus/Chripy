package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/CamusSisyphus/Chripy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	// get env data
	godotenv.Load()

	// get the DB_URL from the environment
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWTSECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	//connection to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	const port = "8080"
	const filepathRoot = "."

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		jwtSecret:      jwtSecret,
		polkaKey:       polkaKey,
	}

	servemux := http.NewServeMux()
	servemux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	servemux.HandleFunc("GET /api/healthz", healthzCheck)
	servemux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirps)
	servemux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	servemux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerGetChirpByID)
	servemux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandlerFunc)
	servemux.HandleFunc("POST /admin/reset", apiCfg.metricsResetHandlerFunc)
	servemux.HandleFunc("POST /api/users", apiCfg.handlerCreateUsers)
	servemux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)
	servemux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	servemux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)
	servemux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerpolkaWebhook)

	servemux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUsers)

	servemux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
	server := http.Server{
		Handler: servemux,
		Addr:    ":" + port}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	jwtSecret      string
	polkaKey       string
}

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(w, req)
		a.fileserverHits.Add(1)
	})
}

func (a *apiConfig) metricsHandlerFunc(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hitcountString := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
<body>
<h1>Welcome, Chirpy Admin</h1>
<p>Chirpy has been visited %d times!</p>
</body>
</html>
	`,
		a.fileserverHits.Load())
	io.WriteString(w, hitcountString)
}

func (a *apiConfig) metricsResetHandlerFunc(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	a.fileserverHits.Store(0)
	a.db.DeleteUsers(context.Background())

}

func healthzCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")

}
