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
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "techdbuser", "1234", "techdb")

	db, err := sql.Open("postgres", dbinfo)
	usecases := useCase.NewUseCase(repository.NewDBStore(db))
	api := delivery.NewHandlers(usecases)

	r := mux.NewRouter()
	r.HandleFunc("/forum/create", api.PostForum).Methods("POST")
	//r.HandleFunc("/users", api.SignUp).Methods("POST")
	//r.HandleFunc("/users/login", api.Login).Methods("POST")
	log.Println("http server started on 5000 port: ")
	err = http.ListenAndServe(":5000", r)
	if err != nil {
		log.Println(err)
		return
	}
}
