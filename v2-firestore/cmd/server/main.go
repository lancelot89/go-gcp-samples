package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"

	fs "go-gcp-samples/v2-firestore/internal/firestore"
	"go-gcp-samples/v2-firestore/internal/handler"
	"go-gcp-samples/v2-firestore/internal/service"
	"go-gcp-samples/v2-firestore/pkg/config"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig()
	if cfg.ProjectID == "" {
		log.Fatalf("PROJECT_ID environment variable not set")
	}

	// Initialize Firestore client
	var client *firestore.Client
	var err error
	if cfg.FirestoreEmulatorHost != "" {
		// Use emulator if FIRESTORE_EMULATOR_HOST is set
		client, err = firestore.NewClient(ctx, cfg.ProjectID, option.WithEndpoint(cfg.FirestoreEmulatorHost))
	} else {
		// Otherwise, connect to production Firestore
		client, err = firestore.NewClient(ctx, cfg.ProjectID)
	}
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Failed to close Firestore client: %v", err)
		}
	}()

	repo := fs.NewTodoRepository(client)
	svc := service.NewTodoService(repo)
	handler := handler.NewTodoHandler(svc)

	r := mux.NewRouter()
	r.HandleFunc("/todos", handler.CreateTodo).Methods("POST")
	r.HandleFunc("/todos", handler.GetTodo).Methods("GET")
	
	// ヘルスチェックエンドポイント
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods("GET")
	
	r.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		// Firestore接続確認
		_, err := client.Collection("health_check").Doc("test").Get(ctx)
		if err != nil && err.Error() != "rpc error: code = NotFound" {
			// Firestoreへの接続に問題がある場合
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("NOT READY"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("READY"))
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
