package model

type Endpoint string

const (
	EndpointChat      Endpoint = "chat"
	EndpointResponses Endpoint = "responses"
)

type Info struct {
	Name     string
	Provider string
	Endpoint Endpoint
}

type Model[C any] struct {
	Info   Info
	Config C
}
