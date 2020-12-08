package controllers

import (
	"net/http"

	"github.com/diptyojha/goLngFirstProject/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To The GoMap2020 API")

}
