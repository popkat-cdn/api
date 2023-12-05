package auth

import (
	"Popkat/api"
	"Popkat/state"
	"Popkat/types"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	docs "github.com/infinitybotlist/eureka/doclib"
	"github.com/infinitybotlist/eureka/uapi"
	"github.com/infinitybotlist/eureka/uapi/ratelimit"
	"go.uber.org/zap"
)

// Tags
const tagName = "Authentication"

type Router struct{}

func (b Router) Tag() (string, string) {
	return tagName, "The API endpoints listed below are related to allowing our users to authenticate on the website."
}

// Types
type Message struct {
	Message string `json:"message" description:"We haven't actually done things with this endpoint yet, so you are going to get a random message instead."`
}

// Routers
func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/auth/login",
		OpId:    "login",
		Method:  uapi.GET,
		Docs:    LoginDocs,
		Handler: LoginRoute,
		Auth: []uapi.AuthType{
			{
				URLVar: "id",
				Type:   api.TargetTypeUser,
			},
		},
	}.Route(r)
}

// Login
func LoginDocs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Login",
		Description: "Recieve a temporary token to access our panel on the website.",
		Resp:        Message{},
	}
}

func LoginRoute(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	limit, err := ratelimit.Ratelimit{
		Expiry:      5 * time.Minute,
		MaxRequests: 10,
		Bucket:      "auth/login",
		Identifier: func(r *http.Request) string {
			return d.Auth.ID
		},
	}.Limit(d.Context, r)

	if err != nil {
		state.Logger.Error("Error while ratelimiting", zap.Error(err), zap.String("bucket", "auth/login"))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}

	if limit.Exceeded {
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "You are being ratelimited. Please try again in " + limit.TimeToReset.String(),
			},
			Headers: limit.Headers(),
			Status:  http.StatusTooManyRequests,
		}
	}

	Response := &Message{
		Message: "reee",
	}

	return uapi.HttpResponse{
		Status: http.StatusOK,
		Json:   Response,
	}
}
