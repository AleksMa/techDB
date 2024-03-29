package main

import (
	//"database/sql"
	"fmt"
	"github.com/AleksMa/techDB/delivery"
	"github.com/AleksMa/techDB/models"
	"github.com/AleksMa/techDB/repository"
	useCase "github.com/AleksMa/techDB/usecase"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {

	//t, _ := time.Parse(time.RFC3339Nano, "2017-01-01 00:00:00 +0000 UTC")
	//fmt.Println(t)

	//layout := "2006-01-02 00:00:00 +0000 UTC"

	//dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", "docker", "docker", "docker")

	db, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			User:     "docker",
			Password: "docker",
			Port:     5432,
			Database: "docker",
		},
		MaxConnections: 50,
	})

	// db, err := sql.Open("postgres", dbinfo)
	usecases := useCase.NewUseCase(repository.NewDBStore(db))
	api := delivery.NewHandlers(usecases)

	_, err = db.Exec(models.InitScript)

	if err != nil {
		fmt.Println(err)
	}

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.HandleFunc("/forum/create", api.CreateForum).Methods("POST")

	r.HandleFunc("/forum/{slug}/create", api.CreateThread).Methods("POST")
	r.HandleFunc("/forum/{slug}/details", api.GetForum).Methods("GET")
	r.HandleFunc("/forum/{slug}/threads", api.GetThreads).Methods("GET")
	r.HandleFunc("/forum/{slug}/users", api.GetUsers).Methods("GET")

	r.HandleFunc("/thread/{slug_or_id}/create", api.CreatePost).Methods("POST")
	r.HandleFunc("/thread/{slug_or_id}/details", api.GetThread).Methods("GET")
	r.HandleFunc("/thread/{slug_or_id}/details", api.UpdateThread).Methods("POST")
	r.HandleFunc("/thread/{slug_or_id}/posts", api.GetPosts).Methods("GET")
	r.HandleFunc("/thread/{slug_or_id}/vote", api.Vote).Methods("POST")

	r.HandleFunc("/user/{nickname}/create", api.CreateUser).Methods("POST")
	r.HandleFunc("/user/{nickname}/profile", api.GetUser).Methods("GET")
	r.HandleFunc("/user/{nickname}/profile", api.UpdateUser).Methods("POST")

	r.HandleFunc("/post/{id}/details", api.GetPostFull).Methods("GET")
	r.HandleFunc("/post/{id}/details", api.UpdatePost).Methods("POST")

	r.HandleFunc("/service/status", api.GetStatus).Methods("GET")
	r.HandleFunc("/service/clear", api.Clear).Methods("POST")

	log.Println("http server started on 5000 port: ")
	err = http.ListenAndServe(":5000", r)
	if err != nil {
		log.Println(err)
		return
	}
}
