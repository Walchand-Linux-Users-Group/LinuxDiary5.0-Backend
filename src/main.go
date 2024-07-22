package main

import (
	"backend/src/controllers"
	"backend/src/db"
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	port := os.Getenv("BACKEND_PORT")

	if port == "" {
		port = "5000"
	}

	muxRouter := mux.NewRouter()
	dbServ, err := db.NewDbAdapter(context.Background())
	if err != nil {
		panic(err)
	}

	userService := controllers.NewUserService(dbServ)

	muxRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Welcome to LinuxDiary5.0"}`))
	}).Methods("GET")

	muxRouter.HandleFunc("/user/registration", func(w http.ResponseWriter, r *http.Request) {
		response, _ := userService.CreateUser(context.Background(), r)

		jsonResponse, err := json.Marshal(response)
		log.Println(response)
		if !response.Success {
			http.Error(w, string(jsonResponse), http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}).Methods("POST")

	corsOptions := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"**", "*"},
			AllowedHeaders:   []string{"X-Requested-With", "Content-Type", "Authorization"},
			AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
			AllowCredentials: true,
		},
	)

	httpRouter := corsOptions.Handler(muxRouter)

	log.Println("Server started at port " + port)
	slog.Error(http.ListenAndServe(":"+port, httpRouter).Error())
}
