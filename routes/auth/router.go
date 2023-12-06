package auth

import (
	dbtypes "Popkat/database/types"
	"Popkat/state"
	"Popkat/types"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	docs "github.com/infinitybotlist/eureka/doclib"
	"github.com/infinitybotlist/eureka/uapi"
	"github.com/infinitybotlist/eureka/uapi/ratelimit"
	"go.uber.org/zap"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"

	"github.com/goombaio/namegenerator"
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
	}.Route(r)
}

// Login
func LoginDocs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Login",
		Description: "Recieve a temporary token to access our panel on the website.",
		Resp:        dbtypes.User{},
		Params: []docs.Parameter{
			{
				Name:        "token",
				In:          "query",
				Description: "Token provided by Firebase Auth",
				Required:    true,
				Schema:      docs.IdSchema,
			},
		},
	}
}

func LoginRoute(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	limit, err := ratelimit.Ratelimit{
		Expiry:      5 * time.Minute,
		MaxRequests: 10,
		Bucket:      "auth/login",
		Identifier: func(r *http.Request) string {
			return r.RemoteAddr
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

	Token := r.URL.Query().Get("token")

	if Token == "" {
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "You are missing a query, or path. Please check our documentation site to check what you are missing.",
			},
			Status: http.StatusBadRequest,
		}
	}

	opt := option.WithCredentialsFile("firebase_service.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Auth(state.Context)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.VerifyIDToken(state.Context, Token)
	if err != nil {
		// this shit be fake asf bro
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "The token provided in the 'token' query is invalid.",
			},
			Status: http.StatusBadRequest,
		}
	}

	/*
		{
			uid: token.uid,
			email: token.firebase.identities.email[0]
			auth_time: token.auth_time,
			expire_at: token.exp
		}
	*/

	user := &dbtypes.User{
		Username: namegenerator.NewNameGenerator(time.Now().UTC().UnixNano()).Generate(),
		UserID:   token.UID,
		Banned:   false,
		Avatar:   "https://popkat.select-list.xyz/logo.png",
		Token:    uuid.NewString(),
	}

	return uapi.HttpResponse{
		Status: http.StatusOK,
		Json:   user,
	}
}
