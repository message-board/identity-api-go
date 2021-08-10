package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/message-board/identity-go/internal/interfaces/rest"
)

func main() {
	// ctx := context.Background()

	// application := service.NewApplication(ctx)

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		rest.RunServer(func(router chi.Router) http.Handler {
			return rest.HandlerFromMux(
				rest.NewUserResource(),
				router,
			)
		})
	case "grpc":
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
