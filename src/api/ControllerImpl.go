package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"chatServer/src/chatserver/data"
)

// ControllerImpl struct for api controller
type ControllerImpl struct {
	service Service
}


// NewControllerImpl returns ControllerImpl
func NewControllerImpl(service Service) *ControllerImpl {
	return &ControllerImpl{
		service:	service,
	}
}


// Register the endpoints this controller handles
func (controller *ControllerImpl) Register() {
	http.HandleFunc("/rest/v1/messages", controller.APIHandler)
	http.ListenAndServe(":3000", nil)
}


// APIHandler handles the endpoints
func (controller *ControllerImpl) APIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		controller.PostMessage(w, r)
	} else if r.Method == http.MethodGet {
		controller.GetMessages(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}


// PostMessage controller is for posting a message
func (controller *ControllerImpl) PostMessage(w http.ResponseWriter, r *http.Request) {

	type BadResponse struct {
		StatusCode int       `json:"statusCode"`
		Message    string    `json:"message"`
	}

	w.Header().Set("Content-Type", "application/json")

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(BadResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
		return
	}

	// Unmarshal the body
	var message data.Message
	err = json.Unmarshal(bodyBytes, &message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(BadResponse{
			StatusCode: http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	//validate request body
	if message.UserID == 0 || message.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadResponse{
			StatusCode: http.StatusBadRequest,
			Message: "UserId or Text is empty",
		})
		return
	}

	resp, err := controller.service.PostMessage(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(BadResponse{
			StatusCode: http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}


// GetMessages controller is for retrieving the messages
func (controller *ControllerImpl) GetMessages(w http.ResponseWriter, r *http.Request) {

	userIDKey, _ := r.URL.Query()["userId"]
	roomIDKey, _ := r.URL.Query()["roomId"]

	userID := 9223372036854775807 // int64 max value
	roomID := 9223372036854775807 // int64 max value
	if userIDKey != nil {
		userID, _ = strconv.Atoi(userIDKey[0])
	}
	if roomIDKey != nil {
		roomID, _ = strconv.Atoi(roomIDKey[0])
	}

	messages := controller.service.GetMessages(userID, roomID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

