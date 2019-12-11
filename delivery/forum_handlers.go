package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/AleksMa/techDB/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

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

	forum, err := handlers.usecases.GetForumBySlug(slug)
	if err != nil {
		body, _ := json.Marshal(err)
		WriteResponse(w, body, err.Code)
		return
	}

	body, _ := json.Marshal(forum)

	WriteResponse(w, body, http.StatusOK)
}
