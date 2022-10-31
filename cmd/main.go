package main

import (
	"context"
	"github.com/giusepperoro/avitotest/internal/database"
	handlers "github.com/giusepperoro/avitotest/internal/handlers"
	_ "github.com/jackc/pgx/v4"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	db, err := database.New(ctx)
	if err != nil {
		log.Fatalf("unable to create db connection: %v", err)
	}
	log.Println("...Server started...")
	http.HandleFunc("/create", handlers.HandleCreateUser(db))
	http.HandleFunc("/balance", handlers.HandleGetBalance(db))
	http.HandleFunc("/refill", handlers.HandleRefill(db))
	http.HandleFunc("/withdrawal", handlers.HandleWithdrawal(db))
	http.HandleFunc("/processWithdrawal", handlers.HandleProcessWithdrawal(db))
	http.HandleFunc("/cancelWithdrawal", handlers.HandleCancelWithdrawal(db))
	err = http.ListenAndServe("0.0.0.0:80", nil)
	log.Fatal("err here`:", err)
}
