package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/diptyojha/goLngFirstProject/api/auth"
	"github.com/diptyojha/goLngFirstProject/api/models"
	"github.com/diptyojha/goLngFirstProject/api/responses"
	"github.com/diptyojha/goLngFirstProject/api/utils/formaterror"
	"github.com/gorilla/mux"
)

func (server *Server) CreateLocation(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	location := models.Location{}
	err = json.Unmarshal(body, &location)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	location.Prepare()
	err = location.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != location.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	locationCreated, err := location.SaveLocation(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, locationCreated.ID))
	responses.JSON(w, http.StatusCreated, locationCreated)
}

func (server *Server) GetLocations(w http.ResponseWriter, r *http.Request) {

	location := models.Location{}

	locations, err := location.FindAllLocations(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, locations)
}

func (server *Server) GetLocation(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	Location := models.Location{}

	LocationReceived, err := Location.FindLocationByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, LocationReceived)
}

func (server *Server) UpdateLocation(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the Location id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the Location exist
	Location := models.Location{}
	err = server.DB.Debug().Model(models.Location{}).Where("id = ?", pid).Take(&Location).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Location not found"))
		return
	}

	// If a user attempt to update a Location not belonging to him
	if uid != Location.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data Locationed
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	LocationUpdate := models.Location{}
	err = json.Unmarshal(body, &LocationUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != LocationUpdate.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	LocationUpdate.Prepare()
	err = LocationUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	LocationUpdate.ID = Location.ID //this is important to tell the model the Location id to update, the other update field are set above

	LocationUpdated, err := LocationUpdate.UpdateALocation(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, LocationUpdated)
}

func (server *Server) DeleteLocation(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid Location id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the Location exist
	Location := models.Location{}
	err = server.DB.Debug().Model(models.Location{}).Where("id = ?", pid).Take(&Location).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this Location?
	if uid != Location.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = Location.DeleteALocation(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
