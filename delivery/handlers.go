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
}

func NewHandlers(ucases useCase2.UseCase) *Handlers {
	return &Handlers{
		usecases: ucases,
	}
}

func WriteResponse(w http.ResponseWriter, body []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
}

func (handlers *Handlers) CreateForum(w http.ResponseWriter, r *http.Request) {
	var newForum models.Forum

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))

	err := json.Unmarshal(body, &newForum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	forum, e := handlers.usecases.PutForum(&newForum)
	if e != nil {
		if e.Code == http.StatusConflict {
			body, _ = json.Marshal(forum)
			WriteResponse(w, body, e.Code)
			return
		}
		body, _ = json.Marshal(e)
		WriteResponse(w, body, e.Code)
		return
	}

	body, _ = json.Marshal(forum)
	WriteResponse(w, body, http.StatusCreated)
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

func (handlers *Handlers) GetThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	threads, _ := handlers.usecases.GetThreadsByForum(slug)

	fmt.Println(threads)

	body, _ := json.Marshal(threads)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (handlers *Handlers) ChangeThread(w http.ResponseWriter, r *http.Request) {
	var thread models.Thread

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))
	vars := mux.Vars(r)
	slug_or_id := vars["slug_or_id"]

	err := json.Unmarshal(body, &thread)
	if err != nil {

	}

	if id, err := strconv.Atoi(slug_or_id); err == nil {
		thread.ID = int64(id)
		thread, _ = handlers.usecases.UpdateThreadWithID(&thread)
	} else {
		thread.Slug = slug_or_id
		thread, _ = handlers.usecases.UpdateThreadWithSlug(&thread)
	}

	fmt.Println(thread)

	body, _ = json.Marshal(thread)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (handlers *Handlers) GetThread(w http.ResponseWriter, r *http.Request) {
	var thread models.Thread
	vars := mux.Vars(r)
	slug_or_id := vars["slug_or_id"]

	if id, err := strconv.Atoi(slug_or_id); err == nil {
		thread, _ = handlers.usecases.GetThreadByID(int64(id))
	} else {
		thread, _ = handlers.usecases.GetThreadBySlug(slug_or_id)
	}

	fmt.Println(thread)

	body, _ := json.Marshal(thread)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
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

func (handlers *Handlers) ChangePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)

	err = json.Unmarshal(body, &post)
	if err != nil {

	}
	post.ID = int64(id)

	handlers.usecases.ChangePost(&post)
}

func (handlers *Handlers) GetPostFull(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, _ := strconv.Atoi(idStr)

	postFull, _ := handlers.usecases.GetPostFull(int64(id))

	fmt.Println(postFull)

	body, _ := json.Marshal(postFull)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (handlers *Handlers) GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts models.Posts
	vars := mux.Vars(r)
	slug_or_id := vars["slug_or_id"]

	if id, err := strconv.Atoi(slug_or_id); err == nil {
		posts, _ = handlers.usecases.GetPostsByThreadID(int64(id))
	} else {
		posts, _ = handlers.usecases.GetPostsByThreadSlug(slug_or_id)
	}

	fmt.Println(posts)

	body, _ := json.Marshal(posts)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (handlers *Handlers) Vote(w http.ResponseWriter, r *http.Request) {
	var vote models.Vote

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))
	vars := mux.Vars(r)
	slug_or_id := vars["slug_or_id"]

	err := json.Unmarshal(body, &vote)
	if err != nil {

	}

	if id, err := strconv.Atoi(slug_or_id); err == nil {
		vote.ThreadID = int64(id)
		handlers.usecases.PutVote(&vote)
	} else {
		handlers.usecases.PutVoteWithSlug(&vote, slug_or_id)
	}
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
