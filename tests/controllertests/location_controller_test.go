package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diptyojha/goLngFirstProject/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateLocation(t *testing.T) {

	err := refreshUserAndLocationTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		Loc_Name     string
		Address      string
		Pincode      string
		TimeZone     string
		creator_id   uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"Loc_Name":"Locations1", "Address": "Raffles place","Pincode":" 111111","TimeZone":"EST", "creator_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			Loc_Name:     "Locations1",
			Address:      "Raffles place",
			Pincode:      "111111",
			TimeZone:     "EST",
			creator_id:   user.ID,
			errorMessage: "",
		},
		{
			inputJSON:    `{"Loc_Name":"Locations1", "Address": "Raffles place","Pincode":" 111111","TimeZone":"EST", "creator_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Incorrect Details",
		},
		{
			// When no token is passed
			inputJSON:    `{"Loc_Name":"When no token is passed", "Address": "Raffles place","Pincode":" 111111","TimeZone":"EST", "creator_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"Loc_Name":"When incorrect token is passed", "Address": "Raffles place","Pincode":" 111111","TimeZone":"EST", "creator_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"Loc_Name": "", "Address": "Raffles place","Pincode":" 111111","TimeZone":"EST", "creator_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Location Name",
		},
		{
			inputJSON:    `{"Loc_Name": "Locations11", "Address": "","Pincode":" 111111","TimeZone":"EST", "creator_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Address",
		},
		{
			inputJSON:    `{"Loc_Name": "Locations17", "Address": "Some where on earth","Pincode":" 111111","TimeZone":"EST"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required CreatorID",
		},
		{
			// When user 2 uses User_1 token
			inputJSON:    `{"Loc_Name": "Locations18", "Address": "Also somewhere on earth","Pincode":" 222222","TimeZone":"EST","creator_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/locations", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateLocation)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["Loc_Name"], v.Loc_Name)
			assert.Equal(t, responseMap["Address"], v.Address)
			assert.Equal(t, responseMap["Pincode"], v.Pincode)
			assert.Equal(t, responseMap["TimeZone"], v.TimeZone)
			assert.Equal(t, responseMap["creator_id"], float64(v.creator_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetLocations(t *testing.T) {

	err := refreshUserAndLocationTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndLocations()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetLocations)
	handler.ServeHTTP(rr, req)

	var posts []models.Location
	err = json.Unmarshal([]byte(rr.Body.String()), &posts)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(posts), 2)
}

// func TestGetLocationByID(t *testing.T) {

// 	err := refreshUserAndLocationTable()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	location, err := seedOneUserAndOneLocation()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	locationSample := []struct {
// 		id           string
// 		statusCode   int
// 		Loc_Name     string
// 		Address      string
// 		Pincode      string
// 		TimeZone     string
// 		creator_id   uint32
// 		errorMessage string
// 	}{
// 		{
// 			id:         strconv.Itoa(int(location.ID)),
// 			statusCode: 200,
// 			Loc_Name:   location.Loc_Name,
// 			Address:    location.Address,
// 			Pincode:    location.Pincode,
// 			TimeZone:   location.TimeZone,
// 			creator_id: location.CreatorID,
// 		},

// 		{
// 			id:         "unknwon",
// 			statusCode: 400,
// 		},
// 	}
// 	for _, v := range locationSample {

// 		req, err := http.NewRequest("GET", "/posts", nil)
// 		if err != nil {
// 			t.Errorf("this is the error: %v\n", err)
// 		}
// 		req = mux.SetURLVars(req, map[string]string{"id": v.id})

// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(server.GetLocation)
// 		handler.ServeHTTP(rr, req)

// 		responseMap := make(map[string]interface{})
// 		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
// 		if err != nil {
// 			log.Fatalf("Cannot convert to json: %v", err)
// 		}
// 		assert.Equal(t, rr.Code, v.statusCode)

// 		if v.statusCode == 200 {
// 			assert.Equal(t, location.Loc_Name, responseMap["Loc_Name"])
// 			assert.Equal(t, location.Address, responseMap["Address"])
// 			assert.Equal(t, location.Pincode, responseMap["Pincode"])
// 			assert.Equal(t, location.TimeZone, responseMap["TimeZone"])
// 			assert.Equal(t, float64(location.CreatorID), responseMap["creator_id"]) //the response author id is float64
// 		}
// 	}
// }
