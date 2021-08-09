package ports

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/message-board/identity/app"
	"github.com/message-board/identity/app/command"
	"github.com/message-board/identity/pkg/requests"
	util "github.com/message-board/identity/util"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{
		app: app,
	}
}

func (h HttpServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		util.WriteResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	request := &requests.CreateUserRequest{}
	if err := render.Decode(r, request); err != nil {
		util.WriteResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}

	command := command.NewCreateUserCommand(request.Id, request.EmailAddress, request.Password)
	err := h.app.Commands.CreateUserCommandHandler.Handle(r.Context(), command)
	if err != nil {
		util.WriteResponse(w, "Failed to create user "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
