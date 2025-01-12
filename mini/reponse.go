package main

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Token   string `json:"token,omitempty"`
}
