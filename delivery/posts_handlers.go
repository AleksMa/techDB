package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/AleksMa/techDB/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func (handlers *Handlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	var posts models.Posts

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))
	vars := mux.Vars(r)
	slug_or_id := vars["slug_or_id"]

	err := json.Unmarshal(body, &posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postsAdded := make(models.Posts, len(posts))
	var e *models.Error
	created := time.Now()

	for i, _ := range posts {
		posts[i].Created = created
		if id, err := strconv.Atoi(slug_or_id); err == nil {
			posts[i].Thread = int64(id)
			postsAdded[i], e = handlers.usecases.PutPost(posts[i])
		} else {
			postsAdded[i], e = handlers.usecases.PutPostWithSlug(posts[i], slug_or_id)
		}
		if e != nil {
			body, _ = json.Marshal(e)
			WriteResponse(w, body, e.Code)
			return
		}
	}

	body, _ = json.Marshal(postsAdded)
	WriteResponse(w, body, http.StatusCreated)
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
