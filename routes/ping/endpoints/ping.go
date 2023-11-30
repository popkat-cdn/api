package ping_end

import (
	"net/http"
	types "Popkat/types"
	docs "github.com/infinitybotlist/eureka/doclib"
	"github.com/infinitybotlist/eureka/uapi"
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Ping",
		Description: "Check if i am working or not",
		Resp: types.Ping{},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	return uapi.HttpResponse{
		Status: http.StatusOK,
		Json:   types.Ping{
			Response: "Hello, world!",
		},
	}
}