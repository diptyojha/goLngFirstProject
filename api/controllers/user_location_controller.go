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

func (server *Server) CreateUserLocation(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userLocation := models.UserLocation{}
	err = json.Unmarshal(body, &userLocation)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userLocation.Prepare()
	err = userLocation.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != userLocation.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	userLocationCreated, err := userLocation.SaveUserLocation(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, userLocationCreated.ID))
	responses.JSON(w, http.StatusCreated, userLocationCreated)
}

func (server *Server) GetUserLocation(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	UserLocation := models.UserLocation{}
	Location := models.Location{}

	UserLocationReceived, err := UserLocation.FindUserLocationByID(server.DB, pid)
	LocationReceived, err := Location.FindLocationByID(server.DB, UserLocationReceived.UserLocationID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, LocationReceived)
}

func (server *Server) UpdateUserLocation(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the UserLocation id is valid
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

	// Check if the UserLocation exist
	UserLocation := models.UserLocation{}
	err = server.DB.Debug().Model(models.UserLocation{}).Where("id = ?", pid).Take(&UserLocation).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("UserLocation not found"))
		return
	}

	// If a user attempt to update a UserLocation not belonging to him
	if uid != UserLocation.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data UserLocationed
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	UserLocationUpdate := models.UserLocation{}
	err = json.Unmarshal(body, &UserLocationUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != UserLocationUpdate.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	UserLocationUpdate.Prepare()
	err = UserLocationUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	UserLocationUpdate.ID = UserLocation.ID //this is important to tell the model the UserLocation id to update, the other update field are set above

	UserLocationUpdated, err := UserLocationUpdate.UpdateAUserLocation(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, UserLocationUpdated)
}

func (server *Server) DeleteUserLocation(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid UserLocation id given to us?
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

	// Check if the UserLocation exist
	UserLocation := models.UserLocation{}
	err = server.DB.Debug().Model(models.UserLocation{}).Where("id = ?", pid).Take(&UserLocation).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this UserLocation?
	if uid != UserLocation.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = UserLocation.DeleteAUserLocation(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
