package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// retrieve the "id" URL paramter from the current request context, then convert it to an int
// if the operation isn't successful, return 0 and an error
func (app *application) readIDParam(r *http.Request) (int64, error) {
	// when httprouter parses a request, interpolated URL params will be stored in the request context
	// ParamsFromContext() func will retrieve a slice containing these names and values
	// []params, params = key, value
	params := httprouter.ParamsFromContext(r.Context())

	// use ByName() method to get value of "id" parameter from the slice
	// all movies will have unique positive integer ID, but val returned by ByName is always a str
	// conver it to a base 10 int (with bit size of 64)
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)

	if err != nil {
		return err
	}
	js = append(js, '\n')

	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
