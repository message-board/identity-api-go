package rest

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/message-board/identity-go/pkg/requests"
)

type UserResource struct {
}

func NewUserResource() UserResource {
	return UserResource{}
}

func (ur UserResource) CreateUser(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		render.Render(w, r, ErrUnsupportedMediaType())
		return
	}

	request := &requests.CreateUserRequest{}
	if err := render.Decode(r, request); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// command := command.NewCreateUserCommand(request.Id, request.EmailAddress, request.Password)
	// err := h.app.Commands.CreateUserCommandHandler.Handle(r.Context(), command)
	// if err != nil {
	// 	util.WriteResponse(w, "Failed to create user "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusNoContent)
}

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrUnsupportedMediaType() render.Renderer {
	return &ErrResponse{
		Err:            nil,
		HTTPStatusCode: 415,
		StatusText:     "Unsupported Media Type",
		ErrorText:      "Unsupported Media Type",
	}
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
