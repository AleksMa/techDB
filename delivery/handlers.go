package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/AleksMa/techDB/models"
	useCase2 "github.com/AleksMa/techDB/usecase"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Handlers struct {
	usecases useCase2.UseCase
	//utils    utils.HandlersUtils
}

func NewHandlers(ucases useCase2.UseCase) *Handlers {
	return &Handlers{
		usecases: ucases,
		//utils:    utils,
	}
}

func (handlers *Handlers) PostForum(w http.ResponseWriter, r *http.Request) {
	var newForum models.Forum

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))

	err := json.Unmarshal(body, &newForum)
	if err != nil {

	}

	handlers.usecases.PutForum(&newForum)
}

func (handlers *Handlers) PostThread(w http.ResponseWriter, r *http.Request) {
	var newThread models.Thread

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))
	vars := mux.Vars(r)
	slug := vars["slug"]

	err := json.Unmarshal(body, &newThread)
	if err != nil {

	}
	newThread.Forum = slug

	handlers.usecases.PutThread(&newThread)
}

func (handlers *Handlers) PutPost(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))
	vars := mux.Vars(r)
	slug_or_id := vars["slug_or_id"]

	err := json.Unmarshal(body, &post)
	if err != nil {

	}

	if id, err := strconv.Atoi(slug_or_id); err == nil {
		handlers.usecases.PutPost(&post, int64(id))
	} else {
		handlers.usecases.PutPostWithSlug(&post, slug_or_id)
	}
}

func (handlers *Handlers) PostUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	err := json.Unmarshal(body, &user)
	if err != nil {

	}
	user.Nickname = nickname

	handlers.usecases.PutUser(&user)
}

func (handlers *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user, _ := handlers.usecases.GetUserByNickname(nickname)

	body, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (handlers *Handlers) GetForum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	forum, _ := handlers.usecases.GetForumBySlug(slug)

	fmt.Println(forum)

	body, _ := json.Marshal(forum)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (handlers *Handlers) GetThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	threads, _ := handlers.usecases.GetThreadsByForum(slug)

	fmt.Println(threads)

	body, _ := json.Marshal(threads)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (handlers *Handlers) ChangeUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	err := json.Unmarshal(body, &user)
	if err != nil {

	}
	user.Nickname = nickname

	handlers.usecases.ChangeUser(&user)
}

func (handlers *Handlers) GetStatus(w http.ResponseWriter, r *http.Request) {
	status, _ := handlers.usecases.GetStatus()

	fmt.Println(status)

	body, _ := json.Marshal(status)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (handlers *Handlers) Clear(w http.ResponseWriter, r *http.Request) {
	handlers.usecases.RemoveAllData()
}
