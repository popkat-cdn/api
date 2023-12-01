package tests

import (
	types "Popkat/types"
	"github.com/go-chi/chi/v5"
	docs "github.com/infinitybotlist/eureka/doclib"
	"github.com/infinitybotlist/eureka/uapi"
	"net/http"
)

const tagName = "Tests"

type Router struct{}

func (b Router) Tag() (string, string) {
	return tagName, "These API endpoints are related to testing our services, and is primarily used by Staff; however is open to the public."
}

func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/tests/ping",
		OpId:    "ping",
		Method:  uapi.GET,
		Docs:    PingDocs,
		Handler: PingRoute,
	}.Route(r)

	uapi.Route{
		Pattern: "/tests/github",
		OpId:    "github",
		Method:  uapi.GET,
		Docs:    GithubDocs,
		Handler: GithubRoute,
	}.Route(r)
}

// Ping
func PingDocs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Ping",
		Description: "Check if i am working or not",
		Resp:        types.MessageResp{},
	}
}

func PingRoute(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	Response := &types.MessageResp{
		Message: "Hello, world!",
	}

	return uapi.HttpResponse{
		Status: http.StatusOK,
		Json:   Response,
	}
}

// Github
func GithubDocs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Github",
		Description: "Check out our Github Organization",
		Resp:        types.MessageResp{},
	}
}

func GithubRoute(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	return uapi.HttpResponse{
		Status:   http.StatusTemporaryRedirect,
		Redirect: "https://github.com/selectlist",
	}
}
