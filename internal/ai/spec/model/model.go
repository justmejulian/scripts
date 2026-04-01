package model

type Endpoint string

type ProviderName string

type Info struct {
	Name     string
	Provider ProviderName
	Endpoint Endpoint
}

type Model[C any] struct {
	Info   Info
	Config C
}
