package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"techDB/api/delivery"
	"techDB/api/repository"
	useCase "techDB/api/usecase"
)


func main() {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "techdbuser", "1234", "techdb")

	db, err := sql.Open("postgres", dbinfo)
	users := useCase.NewUserUseCase(repository.NewUserDBStore(db))
	usersApi := delivery.NewUsersHandlers(users)



	r := mux.NewRouter()
	r.HandleFunc("/users", usersApi.SignUp).Methods("POST")
	r.HandleFunc("/users/login", usersApi.Login).Methods("POST")
	log.Println("Users http server started on 5000 port: ")
	err = http.ListenAndServe(":5000", r)
	if err != nil {
		log.Println(err)
		return
	}
}
