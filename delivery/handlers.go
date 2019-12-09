package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/AleksMa/techDB/models"
	useCase2 "github.com/AleksMa/techDB/usecase"
	"io/ioutil"
	"net/http"
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

	handlers.usecases.PostForum(&newForum)
}
