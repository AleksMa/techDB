package main

import (
	"database/sql"
	"fmt"
	"github.com/AleksMa/techDB/delivery"
	"github.com/AleksMa/techDB/repository"
	useCase "github.com/AleksMa/techDB/usecase"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {

	//t, _ := time.Parse(time.RFC3339Nano, "2017-01-01 00:00:00 +0000 UTC")
	//fmt.Println(t)

	//layout := "2006-01-02 00:00:00 +0000 UTC"

	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "techdbuser", "1234", "techdb")

	db, err := sql.Open("postgres", dbinfo)
	usecases := useCase.NewUseCase(repository.NewDBStore(db))
	api := delivery.NewHandlers(usecases)

	r := mux.NewRouter()
	r.HandleFunc("/forum/create", api.PostForum).Methods("POST")

	r.HandleFunc("/forum/{slug}/create", api.PostThread).Methods("POST")
	r.HandleFunc("/forum/{slug}/details", api.GetForum).Methods("GET")
	r.HandleFunc("/forum/{slug}/threads", api.GetThreads).Methods("GET")

	r.HandleFunc("/thread/{slug_or_id}/create", api.PutPost).Methods("POST")

	r.HandleFunc("/user/{nickname}/create", api.PostUser).Methods("POST")
	r.HandleFunc("/user/{nickname}/profile", api.GetUser).Methods("GET")
	r.HandleFunc("/user/{nickname}/profile", api.ChangeUser).Methods("POST")

	r.HandleFunc("/service/status", api.GetStatus).Methods("GET")
	r.HandleFunc("/service/clear", api.Clear).Methods("POST")

	log.Println("http server started on 5000 port: ")
	err = http.ListenAndServe(":5000", r)
	if err != nil {
		log.Println(err)
		return
	}
}
