package meetup_endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/when-are-you-free/backend/auth"
	"github.com/when-are-you-free/backend/httputil"
	"github.com/when-are-you-free/backend/storage"
)

type MeetupHandler struct {
	s *storage.Storage
}

func New(s *storage.Storage) *MeetupHandler {
	return &MeetupHandler{
		s: s,
	}
}

type postRequestData struct {
	Meetup *storage.Meetup        `json:"meetup"`
	Person *storage.InvitedPerson `json:"person"`
}

type postResponseData struct {
	storage.Meetup
	Uuid string `json:"uuid"`
}

func (mh *MeetupHandler) Post(responseWriter http.ResponseWriter, request *http.Request) {
	userToken := request.Header.Get(auth.UserTokenHeader)

	postRequestData := &postRequestData{}
	if errParse := httputil.ParseBody(request.Body, postRequestData, responseWriter); errParse != nil {
		return
	}

	uuid := uuid.NewString()
	postRequestData.Person.UserToken = userToken
	postRequestData.Meetup.InvitedPeople = []*storage.InvitedPerson{
		postRequestData.Person,
	}

	mh.s.ReadWrite(func(data *storage.Data) {
		data.Meetups[uuid] = postRequestData.Meetup
	})

	responseBytes, errMarshalResponse := json.Marshal(&postResponseData{
		Uuid:   uuid,
		Meetup: *postRequestData.Meetup,
	})
	if errMarshalResponse != nil {
		log.Println("Error marshalling MeetupHandler.Post Response:", errMarshalResponse)
		httputil.WritePlainError(responseWriter, http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	_, errWrite := responseWriter.Write(responseBytes)
	if errWrite != nil {
		log.Println("Error writing MeetupHandler.Post Response:", errWrite)
	}
}

func (mh *MeetupHandler) Get(responseWriter http.ResponseWriter, request *http.Request) {
	userToken := request.Header.Get(auth.UserTokenHeader)
	uuid := chi.URLParam(request, "uuid")
	var (
		meetup   *storage.Meetup
		contains bool
	)
	mh.s.ReadWrite(func(data *storage.Data) {
		meetup, contains = data.Meetups[uuid]
	})

	if !contains {
		httputil.WritePlainError(responseWriter, http.StatusNotFound)
		return
	}

	if !meetup.HasAccess(userToken) {
		httputil.WritePlainError(responseWriter, http.StatusForbidden)
		return
	}

	responseBytes, errMarshalResponse := json.Marshal(meetup)
	if errMarshalResponse != nil {
		log.Println("Error marshalling MeetupHandler.Post Response:", errMarshalResponse)
		httputil.WritePlainError(responseWriter, http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	_, errWrite := responseWriter.Write(responseBytes)
	if errWrite != nil {
		log.Println("Error writing MeetupHandler.Post Response:", errWrite)
	}
}

func (mh *MeetupHandler) Delete(responseWriter http.ResponseWriter, request *http.Request) {
	// userToken := request.Header.Get(UserTokenHeader)
}

func (mh *MeetupHandler) Patch(responseWriter http.ResponseWriter, request *http.Request) {
	// userToken := request.Header.Get(UserTokenHeader)
}
