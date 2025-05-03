package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"txparser/internal/service"
)

type Subscription struct {
	userService service.Subscription
}

func NewSubscription(userService service.Subscription) *Subscription {
	return &Subscription{
		userService: userService,
	}
}

func (u *Subscription) Subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	var body struct {
		Address string `json:"address"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("subscribed to %s", body.Address)
	err = u.userService.Subscribe(body.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
