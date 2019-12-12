package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/AleksMa/techDB/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (handlers *Handlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	var posts models.Posts
	var err error
	var tempID int
	id := -1

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))
	vars := mux.Vars(r)
	slug_or_id := vars["slug_or_id"]

	err = json.Unmarshal(body, &posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tempID, err = strconv.Atoi(slug_or_id)
	if err == nil {
		id = tempID
	}

	postsAdded := make(models.Posts, len(posts))
	var e *models.Error
	created := time.Now()

	for i, _ := range posts {
		posts[i].Created = created
		if id != -1 {
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

	if len(posts) == 0 {
		//var thread models.Thread
		if id != -1 {
			_, e = handlers.usecases.GetThreadByID(int64(id))
		} else {
			_, e = handlers.usecases.GetThreadBySlug(slug_or_id)
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

func (handlers *Handlers) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var setPost, post models.Post
	var e *models.Error

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)

	err = json.Unmarshal(body, &setPost)
	if err != nil {
		http.Error(w, "unmarshal error", http.StatusInternalServerError)
		return
	}

	setPost.ID = int64(id)

	post, e = handlers.usecases.ChangePost(&setPost)
	if e != nil {
		body, _ = json.Marshal(e)
		WriteResponse(w, body, e.Code)
		return
	}

	body, _ = json.Marshal(post)
	WriteResponse(w, body, http.StatusOK)
}

func (handlers *Handlers) GetPostFull(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, _ := strconv.Atoi(idStr)

	fields := strings.Split(r.URL.Query().Get("related"), ",")
	fmt.Println("RELATED???:", fields)

	postFull, err := handlers.usecases.GetPostFull(int64(id), fields)
	if err != nil {
		body, _ := json.Marshal(err)
		WriteResponse(w, body, err.Code)
		return
	}

	fmt.Println(postFull)

	body, _ := json.Marshal(postFull)
	WriteResponse(w, body, http.StatusOK)
}

func (handlers *Handlers) GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts models.Posts
	var e *models.Error

	vars := mux.Vars(r)
	slug_or_id := vars["slug_or_id"]

	query := r.URL.Query()
	var params models.PostParams
	var err error

	params.Limit, err = strconv.Atoi(query.Get("limit"))
	if err != nil {
		params.Limit = -1
	}
	params.Since, err = strconv.Atoi(query.Get("since"))
	if err != nil {
		params.Since = -1
	}
	fmt.Println("SINCE: ", params.Since)
	params.Desc = query.Get("desc") == "true"

	switch query.Get("sort") {
	case "flat":
		params.Sort = models.Flat
	case "tree":
		params.Sort = models.Tree
	case "parent_tree":
		params.Sort = models.ParentTree
	}

	if id, err := strconv.Atoi(slug_or_id); err == nil {
		posts, e = handlers.usecases.GetPostsByThreadID(int64(id), params)
	} else {
		posts, e = handlers.usecases.GetPostsByThreadSlug(slug_or_id, params)
	}
	if e != nil {
		body, _ := json.Marshal(e)
		WriteResponse(w, body, e.Code)
		return
	}

	fmt.Println(posts)

	body, _ := json.Marshal(posts)
	WriteResponse(w, body, http.StatusOK)
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

	var thread models.Thread
	var e *models.Error

	if id, err := strconv.Atoi(slug_or_id); err == nil {
		vote.ThreadID = int64(id)
		thread, e = handlers.usecases.PutVote(&vote)
	} else {
		thread, e = handlers.usecases.PutVoteWithSlug(&vote, slug_or_id)
	}
	if e != nil {
		body, _ = json.Marshal(e)
		WriteResponse(w, body, http.StatusNotFound)
		return
	}

	body, _ = json.Marshal(thread)
	WriteResponse(w, body, http.StatusOK)
}
