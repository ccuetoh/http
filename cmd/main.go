package main

import (
	"encoding/json"

	"http"
)

func main() {
	r := &http.Router{}
	r.HandleFunc("", index)
	r.HandleFunc("/ping", ping)
	r.HandleFunc("/login", login)

	s, err := http.New(http.WithPort(8080), http.WithHandler(r))
	if err != nil {
		panic(err)
	}

	err = s.ListenAndServe()
	if err != nil {
		panic(err)
	}

	print("ok")
}

func index(_ *http.Request, w *http.ResponseWriter) {
	_ = w.WriteStatus(200)
}

func ping(_ *http.Request, w *http.ResponseWriter) {
	_ = w.WriteText(200, "pong")
}

func login(r *http.Request, w *http.ResponseWriter) {
	if r.Method != "POST" {
		_ = w.WriteStatus(405)
		return
	}

	req := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := json.Unmarshal(r.Body, &req)
	if err != nil {
		_ = w.WriteText(400, err.Error())
		return
	}

	if req.Username != "admin" || req.Password != "123" {
		_ = w.WriteStatus(401)
		return
	}

	_ = w.WriteJSON(200, struct {
		Token string `json:"token"`
	}{"example"})
}
