package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAMovieById(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/movie/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("_id", "62416541f69b2bec6950cc8b")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAMovieById)

	handler.ServeHTTP(rr, req)
	if statusCode := rr.Code; statusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", statusCode, http.StatusOK)
	}
}

func TestGetAMovieById_Bad_request(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/movie/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("_id", "62416541f69b2bec6950cc8b")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAMovieById)

	handler.ServeHTTP(rr, req)
	if statusCode := rr.Code; statusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", statusCode, http.StatusBadRequest)
	}
}

func TestGetMyAllMovies(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/movies", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetMyAllMovies)
	handler.ServeHTTP(rr, req)
	if statusCode := rr.Code; statusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", statusCode, http.StatusOK)
	}
	// Check the response body is what we expect.
	// expected := `[{"_id":"62416541f69b2bec6950cc8b","movie":"IronMan","watched":false}]`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}

func TestCreateMovie(t *testing.T) {
	var jsonStr = []byte(`{"movie":"Deadman","watched":"true",}`)

	req, err := http.NewRequest("POST", "/api/movie", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateMovie)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// func TestMarkAsWatched(t *testing.T) {

// }

// func TestDeleteAMovie(t *testing.T) {

// }

func TestDeleteAllMovies(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/api/deleteallmovie", nil)
	if err != nil { 
		t.Fatal(err) 
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteAllMovies)
	handler.ServeHTTP(rr, req)
	if statusCode := rr.Code; statusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", statusCode, http.StatusOK)
	}
}