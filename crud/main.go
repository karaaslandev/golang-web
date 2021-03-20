package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"encoding/json"
	"time"
)

var db *pgxpool.Pool

func main() {
	r := chi.NewRouter()
	username := "postgres"
	password := "12345"
	dbname := "syz"
	connectionString := fmt.Sprintf("postgresql://%v:%v@localhost:5432/%v", username, password, dbname)
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Fatal(err)
	}

	db, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal(err)
	}

	r.Route("/api", apiRouter)

	err = http.ListenAndServe(":8080", r)
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func apiRouter(r chi.Router) {
	r.Route("/categories", categoriesRoute)
}

func categoriesRoute(r chi.Router) {
	r.Get("/", listCategories)
}

type Category struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	VisibleName   string `json:"visible_name"`
	SuperCategory int    `json:"super_category"`
}

func listCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(context.Background(), "SELECT * FROM home.categories")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer rows.Close()
	categories := []Category{}
	for rows.Next() {
		if err = rows.Err(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		values, err := rows.Values()
		json.NewEncoder(w).Encode(values)
		return
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		c := Category{
			int(values[0].(int32)),
			values[1].(string),
			values[2].(string),
			int(values[3].(int32)),
		}
		categories = append(categories, c)
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(categories)
	categories = nil
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
