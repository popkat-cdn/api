package ping

import (
	ping_end "Popkat/routes/ping/endpoints"

	"github.com/go-chi/chi/v5"
	"github.com/infinitybotlist/eureka/uapi"
)

const tagName = "Ping"

type Router struct{}

func (b Router) Tag() (string, string) {
	return tagName, "These API endpoints are related to pinging our services."
}

func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern:      "/ping",
		OpId:         "ping",
		Method:       uapi.GET,
		Docs:         ping_end.Docs,
		Handler:      ping_end.Route,
	}.Route(r)
}