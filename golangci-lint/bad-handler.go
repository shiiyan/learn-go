package golangcilint

import (
	"encoding/json"
	"net/http"
	"strings"
)

type UserHandler struct {
}

type RegisterRequest struct {
	Age      int
	Email    string
	Password string
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.Age == 0 {
		http.Error(w, "age is required", 400)
		return
	}
	if req.Age < 18 {
		http.Error(w, "too young", 400)
		return
	}
	if req.Email == "" {
		http.Error(w, "email is required", 400)
		return
	}
	if !strings.Contains(req.Email, "@") {
		http.Error(w, "invalid email", 400)
		return
	}
	if req.Password == "" {
		http.Error(w, "password is required", 400)
		return
	}
	if len(req.Password) < 8 {
		http.Error(w, "password too short", 400)
		return
	}

	// ...
}

func (h *UserHandler) Reregister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.Age == 0 {
		http.Error(w, "age is required", 400)
		return
	}
	if req.Age < 18 {
		http.Error(w, "too young", 400)
		return
	}
	if req.Email == "" {
		http.Error(w, "email is required", 400)
		return
	}
	if !strings.Contains(req.Email, "@") {
		http.Error(w, "invalid email", 400)
		return
	}

	// ...
}

type SaveRequest struct {
	Age    int
	Email  string
	Active bool
}

func (h *UserHandler) Save(w http.ResponseWriter, r *http.Request) {
	var u SaveRequest
	json.NewDecoder(r.Body).Decode(&u)

	if u.Age > 0 {
		if u.Email != "" {
			if !strings.Contains(u.Email, "@") {
				// business rule
				if u.Active {
					// more business rule
				}
			}
		}
	}

}
