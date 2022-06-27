package validateAuthor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Book struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Author        Author `json:"author"`
	Publication   string `json:"publication"`
	PublishedDate string `json:"published_date"`
}

type Author struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Dob       string `json:"dob"`
	PenName   string `json:"pen_name"`
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "raramuri"
	dbPass := "Dpyadav@123"
	dbName := "test"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getBook(response http.ResponseWriter, request *http.Request) {
	db := dbConn()
	defer db.Close()
	title := request.URL.Query().Get("title")
	includeAuthor := request.URL.Query().Get("includeAuthor")
	var rows *sql.Rows
	var err error
	if title == "" {
		rows, err = db.Query("select * from Books;")
	} else {
		rows, err = db.Query("select * from Books where title=?;", title)
	}
	if err != nil {
		log.Print(err)
	}
	books := []Book{}
	for rows.Next() {
		book := Book{}
		err = rows.Scan(&book.Id, &book.Title, &book.Publication, &book.PublishedDate, &book.Author.Id)
		if err != nil {
			log.Print(err)
		}
		if includeAuthor == "true" {
			row := db.QueryRow("select * from Authors where id=?", book.Author.Id)
			row.Scan(&book.Author.Id, &book.Author.FirstName, &book.Author.LastName, &book.Author.Dob, &book.Author.PenName)
		}
		books = append(books, book)
	}
	json.NewEncoder(response).Encode(books)
}

func getBookById(response http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	db := dbConn()
	defer db.Close()
	bookrow := db.QueryRow("select * from Books where id=?;", id)
	book := Book{}
	author := Author{}
	author_id := 0
	err := bookrow.Scan(&book.Id, &book.Title, &book.Publication, &book.PublishedDate, &author_id)
	if err != nil {
		log.Print(err)
	}
	authorrow := db.QueryRow("select * from Authors where id=?;", author_id)
	err = authorrow.Scan(&author.Id, &author.FirstName, &author.LastName, &author.Dob, &author.PenName)
	if err != nil {
		log.Print(err)
	}
	book.Author = author
	json.NewEncoder(response).Encode(book)
}

func postBook(response http.ResponseWriter, request *http.Request) {
	db := dbConn()
	defer db.Close()
	decoder := json.NewDecoder(request.Body)
	b := Book{}
	err := decoder.Decode(&b)
	if err != nil {
		panic(err)
	}
	res, err := db.Exec("INSERT INTO Books (title, publication, published_date, author_id)\nVALUES (?,?,?,?);", b.Title, b.Publication, b.PublishedDate, b.Author.Id)
	id, _ := res.LastInsertId()
	if err != nil {
		log.Print(err)
		json.NewEncoder(response).Encode(Book{})
	} else {
		b.Id = int(id)
		json.NewEncoder(response).Encode(b)
	}
}

func postAuthor(response http.ResponseWriter, request *http.Request) {
	db := dbConn()
	defer db.Close()
	decoder := json.NewDecoder(request.Body)
	a := Author{}
	err := decoder.Decode(&a)
	if a.FirstName == "" || a.Dob == "" {
		response.WriteHeader(400)
		json.NewEncoder(response).Encode(Author{})
		return
	}
	if err != nil {
		panic(err)
	}
	res, err := db.Exec("INSERT INTO Authors (first_name, last_name, dob, pen_name)\nVALUES (?,?,?,?);", a.FirstName, a.LastName, a.Dob, a.PenName)
	id, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		response.WriteHeader(400)
		json.NewEncoder(response).Encode(Author{})
	} else {
		a.Id = int(id)
		json.NewEncoder(response).Encode(a)
	}
}

func putAuthor(response http.ResponseWriter, request *http.Request) {

}

func putBook(response http.ResponseWriter, request *http.Request) {

}

func deleteAuthor(response http.ResponseWriter, request *http.Request) {

}

func deleteBook(response http.ResponseWriter, request *http.Request) {

}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/book", getBook).Methods(http.MethodGet)

	r.HandleFunc("/book/{id}", getBookById).Methods(http.MethodGet)

	r.HandleFunc("/book", postBook).Methods(http.MethodPost)

	r.HandleFunc("/author", postAuthor).Methods(http.MethodPost)

	r.HandleFunc("/book/{id}", putBook).Methods(http.MethodPut)

	r.HandleFunc("/author/{id}", putAuthor).Methods(http.MethodPut)

	r.HandleFunc("/book/{id}", deleteBook).Methods(http.MethodDelete)

	r.HandleFunc("/author/{id}", deleteAuthor).Methods(http.MethodDelete)

	Server := http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	fmt.Println("Server started at 8000")
	Server.ListenAndServe()

}
