package types

type Ping struct {
	Response string `json:"response" validate:"required" description:"Should return 'Hello, world!'."`
}