package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type J map[string]interface{}

func PathString(req *http.Request, key string) string {
	return mux.Vars(req)[key]
}
func PathInt(req *http.Request, key string) (int, error) {
	str := PathString(req, key)
	return strconv.Atoi(str)
}

// Unmarshals request's body into a given object
func Bind(req *http.Request, entity interface{}) error {
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(entity)
}

func SetUnauthorizedResponse(res http.ResponseWriter) {
	JsonResponse(res)
	res.WriteHeader(http.StatusUnauthorized)
}

func SetConflictResponse(res http.ResponseWriter) {
	JsonResponse(res)
	res.WriteHeader(http.StatusConflict)
}

func SetBadRequestResponse(res http.ResponseWriter) {
	JsonResponse(res)
	res.WriteHeader(http.StatusBadRequest)
}

func SetNotFoundResponse(res http.ResponseWriter) {
	JsonResponse(res)
	res.WriteHeader(http.StatusNotFound)
}

func SetInternalServerErrorResponse(res http.ResponseWriter, err interface{}) {
	log.Println(err)
	JsonResponse(res)
	res.WriteHeader(http.StatusInternalServerError)
}

func SetCreatedResponse(res http.ResponseWriter, entity interface{}, location string) error {
	b, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	body := string(b[:])

	JsonResponse(res)
	res.Header().Set("Location", location)
	res.WriteHeader(http.StatusCreated)
	fmt.Fprint(res, body)
	return nil
}

func SetOKResponse(res http.ResponseWriter, entity interface{}) error {
	JsonResponse(res)
	res.WriteHeader(http.StatusOK)
	return WriteEntity(res, entity)
}

func WriteEntity(res http.ResponseWriter, entity interface{}) error {
	b, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	body := string(b[:])
	fmt.Fprint(res, body)
	return nil
}

func SetNoContentResponse(res http.ResponseWriter) {
	JsonResponse(res)
	res.WriteHeader(http.StatusNoContent)
}

func SetMethodNotAllowedResponse(res http.ResponseWriter) {
	JsonResponse(res)
	res.WriteHeader(http.StatusMethodNotAllowed)
}

func JsonResponse(res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
}
