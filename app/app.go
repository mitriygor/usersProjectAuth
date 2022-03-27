package app

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/mitriygor/usersProjectAuth/domain"
	"github.com/mitriygor/usersProjectAuth/service"
	"github.com/mitriygor/usersProjectLib/logger"
	"log"
	"net/http"
	"os"
)

func Start() {
	sanityCheck()
	router := mux.NewRouter()
	authRepository := domain.NewAuthRepository(getDbClient())
	ah := AuthHandler{service.NewLoginService(authRepository)}

	router.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)
	router.HandleFunc("/auth/register", ah.NotImplementedHandler).Methods(http.MethodPost)
	router.HandleFunc("/auth/refresh", ah.Refresh).Methods(http.MethodPost)
	router.HandleFunc("/auth/verify", ah.Verify).Methods(http.MethodGet)

	address := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	logger.Info(fmt.Sprintf("Starting Auth server on %s:%s", address, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))
}

func getDbClient() *sql.DB {
	client, err := sql.Open("postgres", "user=postgres password=postgres dbname=dbusers sslmode=disable")

	if err != nil {
		panic(err)
	}

	return client
}

func sanityCheck() {
	envProps := []string{
		"SERVER_HOST",
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			logger.Error(fmt.Sprintf("Environment variable %s is not defined", k))
		}
	}
}
