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

func (handlers *Handlers) CreateThread(w http.ResponseWriter, r *http.Request) {
	var newThread models.Thread

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))

	vars := mux.Vars(r)
	slug := vars["slug"]

	err := json.Unmarshal(body, &newThread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("TIME: ", newThread.Created)

	newThread.Forum = slug

	thread, e := handlers.usecases.PutThread(&newThread)
	if e != nil {
		if e.Code == http.StatusConflict {
			body, _ = json.Marshal(thread)
			WriteResponse(w, body, e.Code)
			return
		}
		body, _ = json.Marshal(e)
		WriteResponse(w, body, e.Code)
		return
	}

	body, _ = json.Marshal(thread)
	WriteResponse(w, body, http.StatusCreated)
}

func (handlers *Handlers) GetThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	query := r.URL.Query()
	var params models.ThreadParams
	var err error

	params.Limit, err = strconv.Atoi(query.Get("limit"))
	if err != nil {
		params.Limit = -1
	}
	fmt.Println("SINCE", query.Get("since"))
	params.Since, err = time.Parse(time.RFC3339Nano, query.Get("since"))
	if err != nil {
		params.Since = time.Time{}
	}
	params.Desc = query.Get("desc") == "true"

	threads, e := handlers.usecases.GetThreadsByForum(slug, params)
	if e != nil {
		body, _ := json.Marshal(err)
		WriteResponse(w, body, e.Code)
		return
	}
	fmt.Println(threads)

	body, _ := json.Marshal(threads)

	WriteResponse(w, body, http.StatusOK)
}

func (handlers *Handlers) GetThread(w http.ResponseWriter, r *http.Request) {
	var thread models.Thread
	var e *models.Error

	vars := mux.Vars(r)
	slug_or_id := vars["slug_or_id"]

	if id, err := strconv.Atoi(slug_or_id); err == nil {
		thread, e = handlers.usecases.GetThreadByID(int64(id))
	} else {
		thread, e = handlers.usecases.GetThreadBySlug(slug_or_id)
	}
	if e != nil {
		body, _ := json.Marshal(e)
		WriteResponse(w, body, e.Code)
		return
	}

	fmt.Println(thread)

	body, _ := json.Marshal(thread)
	WriteResponse(w, body, http.StatusOK)
}

func (handlers *Handlers) UpdateThread(w http.ResponseWriter, r *http.Request) {
	var thread models.Thread
	var e *models.Error

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))
	vars := mux.Vars(r)
	slug_or_id := vars["slug_or_id"]

	err := json.Unmarshal(body, &thread)
	if err != nil {
		http.Error(w, "unmarshal error", http.StatusInternalServerError)
		return
	}

	if id, err := strconv.Atoi(slug_or_id); err == nil {
		thread.ID = int64(id)
		thread, e = handlers.usecases.UpdateThreadWithID(&thread)
	} else {
		thread.Slug = slug_or_id
		thread, e = handlers.usecases.UpdateThreadWithSlug(&thread)
	}
	if e != nil {
		body, _ = json.Marshal(e)
		WriteResponse(w, body, e.Code)
		return
	}

	fmt.Println(thread)

	body, _ = json.Marshal(thread)
	WriteResponse(w, body, http.StatusOK)
}
