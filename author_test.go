package validateAuthor

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetBook(t *testing.T) {
	testcases := []struct {
		desc      string
		req       string
		expRes    []Book
		expStatus int
	}{
		{"get all books", "/book", []Book{
			{1, "Arvind", nil, "Arvind", "11/03/2002"},
			{2, "Arvind", nil, "Arvind", "11/03/2002"},
		}, http.StatusOK},
		{"get all books with query param", "/book?title=Arvind", []Book{
			{1, "Arvind", nil, "Arvind", "11/03/2002"}}, http.StatusOK},
		{"get all books with query param", "/book?includeAuthor=true", []Book{
			{1, "Arvind", &Author{}, "Arvind", "11/03/2002"}}, http.StatusOK},
	}
	for j, tc := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "localhost:8000"+tc.req, nil)

		getBook(w, req)

		defer w.Result().Body.Close()

		res, err := io.ReadAll(w.Result().Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		resBooks := []Book{}

		err = json.Unmarshal(res, &resBooks)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", j, tc.desc)
		}

		if !reflect.DeepEqual(resBooks, tc.expRes) {
			t.Errorf("%v test failed %v", j, tc.desc)
		}
	}
}

func TestGetBookById(t *testing.T) {
	testcases := []struct {
		desc      string
		req       string
		expRes    Book
		expStatus int
	}{
		{"get book", "1", Book{1, "Arvind", &Author{1, "Arvind", "Yadav", "11/03/2002", ""}, "", "11/03/2002"}, http.StatusOK},
		{"Id doesn't exist", "1000", Book{}, http.StatusBadRequest},
		{"Invalid Id", "abc", Book{}, http.StatusNotFound},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "localhost:8000/book/"+tc.req, nil)
		getBookById(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}

		res, err := io.ReadAll(w.Result().Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		resBook := Book{}

		err = json.Unmarshal(res, &resBook)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		if resBook != tc.expRes {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestPostBook(t *testing.T) {
	testcases := []struct {
		desc      string
		reqBody   Book
		expRes    Book
		expStatus int
	}{
		{"Valid Details", Book{Title: "Arvind", Author: &Author{Id: 1}, Publication: "Pengiun", PublishedDate: "11/03/2002"}, Book{1, "Arvind", &Author{Id: 1}, "Pengiun", "11/03/2002"}, 200},
		{"Publication should be Scholastic/Pengiun/Arihanth", Book{Title: "Arvind", Author: &Author{Id: 1}, Publication: "Arvind", PublishedDate: "11/03/2002"}, Book{}, http.StatusBadRequest},
		{"Published date should be between 1880 and 2022", Book{Title: "", Author: &Author{Id: 1}, Publication: "", PublishedDate: "1/1/1870"}, Book{}, http.StatusBadRequest},
		{"Published date should be between 1880 and 2022", Book{Title: "", Author: &Author{Id: 1}, Publication: "", PublishedDate: "1/1/2222"}, Book{}, http.StatusBadRequest},
		{"Author should exist", Book{Title: "Arvind", Author: &Author{Id: 2}, Publication: "Pengiun", PublishedDate: "11/03/2002"}, Book{}, http.StatusBadRequest},
		{"Title can't be empty", Book{Title: "", Author: nil, Publication: "", PublishedDate: ""}, Book{}, http.StatusBadRequest},
		{"Book already exists", Book{Title: "Arvind", Author: &Author{Id: 1}, Publication: "Pengiun", PublishedDate: "11/03/2002"}, Book{}, http.StatusConflict},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPost, "localhost:8000/book/", bytes.NewReader(body))
		postBook(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
		res, _ := io.ReadAll(w.Result().Body)
		resBook := Book{}
		json.Unmarshal(res, &resBook)
		if resBook != tc.expRes {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestPostAuthor(t *testing.T) {
	testcases := []struct {
		desc      string
		reqBody   Author
		expRes    Author
		expStatus int
	}{
		{"Valid details", Author{FirstName: "RD", LastName: "Sharma", Dob: "2/11/1989", PenName: "Sharma"}, Author{1, "RD", "Sharma", "2/11/1989", "Sharma"}, http.StatusOK},
		{"InValid details", Author{FirstName: "", LastName: "Sharma", Dob: "2/11/1989", PenName: "Sharma"}, Author{}, http.StatusBadRequest},
		{"Author already exists", Author{FirstName: "RD", LastName: "Sharma", Dob: "2/11/1989", PenName: "Sharma"}, Author{}, http.StatusBadRequest},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPost, "localhost:8000/author", bytes.NewReader(body))
		postAuthor(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
		res, _ := io.ReadAll(w.Result().Body)
		resAuthor := Author{}
		json.Unmarshal(res, &resAuthor)
		if resAuthor != tc.expRes {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestPutBook(t *testing.T) {
	testcases := []struct {
		desc      string
		reqId     string
		reqBody   Book
		expRes    Book
		expStatus int
	}{
		{"Valid Details", "1", Book{Title: "Arvind", Author: nil, Publication: "Pengiun", PublishedDate: "11/03/2002"}, Book{}, 200},
		{"Publication should be Scholastic/Pengiun/Arihanth", "1", Book{Title: "Arvind", Author: nil, Publication: "Arvind", PublishedDate: "11/03/2002"}, Book{}, http.StatusBadRequest},
		{"Published date should be between 1880 and 2022", "1", Book{Title: "", Author: nil, Publication: "", PublishedDate: "1/1/1870"}, Book{}, http.StatusBadRequest},
		{"Published date should be between 1880 and 2022", "1", Book{Title: "", Author: nil, Publication: "", PublishedDate: "1/1/2222"}, Book{}, http.StatusBadRequest},
		{"Author should exist", "1", Book{}, Book{}, http.StatusBadRequest},
		{"Title can't be empty", "1", Book{Title: "", Author: nil, Publication: "", PublishedDate: ""}, Book{}, http.StatusBadRequest},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPost, "localhost:8000/book/"+tc.reqId, bytes.NewReader(body))
		putBook(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
		res, _ := io.ReadAll(w.Result().Body)
		resBook := Book{}
		json.Unmarshal(res, &resBook)
		if resBook != tc.expRes {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestPutAuthor(t *testing.T) {
	testcases := []struct {
		desc      string
		reqBody   Author
		expRes    Author
		expStatus int
	}{
		{"Valid details", Author{FirstName: "RD", LastName: "Sharma", Dob: "2/11/1989", PenName: "Sharma"}, Author{1, "RD", "Sharma", "2/11/1989", "Sharma"}, http.StatusOK},
		{"InValid details", Author{FirstName: "", LastName: "Sharma", Dob: "2/11/1989", PenName: "Sharma"}, Author{}, http.StatusBadRequest},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPost, "localhost:8000/author", bytes.NewReader(body))
		putAuthor(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
		res, _ := io.ReadAll(w.Result().Body)
		resAuthor := Author{}
		json.Unmarshal(res, &resAuthor)
		if resAuthor != tc.expRes {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestDeleteBook(t *testing.T) {
	testcases := []struct {
		desc      string
		reqId     string
		expStatus int
	}{
		{"Valid Details", "1", http.StatusOK},
		{"Book does not exists", "100", http.StatusNotFound},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "localhost:8000/book/"+tc.reqId, nil)
		deleteBook(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestDeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc      string
		reqId     string
		expStatus int
	}{
		{"Valid Details", "1", http.StatusOK},
		{"Author does not exists", "100", http.StatusBadRequest},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "localhost:8000/author/"+tc.reqId, nil)
		deleteAuthor(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}
