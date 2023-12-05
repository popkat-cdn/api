package tests

import (
	"github.com/go-chi/chi/v5"
	docs "github.com/infinitybotlist/eureka/doclib"
	"github.com/infinitybotlist/eureka/uapi"
	"net/http"
)

// Tags
const tagName = "Tests"

type Router struct{}

func (b Router) Tag() (string, string) {
	return tagName, "The API endpoints listed below are related to checking our system's performance and stability."
}

// Types
type Ping struct {
	Message string `json:"message" description:"This message should mention whether our services are functioning or not."`
}

// Routers
func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/tests/ping",
		OpId:    "ping",
		Method:  uapi.GET,
		Docs:    PingDocs,
		Handler: PingRoute,
	}.Route(r)
}

// Ping
func PingDocs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Ping Status",
		Description: "Check my current functionality and ping!",
		Resp:        Ping{},
	}
}

func PingRoute(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	Response := &Ping{
		Message: "Hello, there. I'm fully functioning with no issues.",
	}

	return uapi.HttpResponse{
		Status: http.StatusOK,
		Json:   Response,
	}
}
