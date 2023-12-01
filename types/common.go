package types

type ApiError struct {
	Context map[string]string `json:"context,omitempty" description:"Context of the error. Usually used for validation error contexts"`
	Message string            `json:"message" description:"Message of the error"`
}

type BasicAPIResp struct {
	Message string `json:"message" description:"Response Message for Endpoint"`
}
