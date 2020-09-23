package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/latihan/connection"
	"github.com/latihan/model"
)

var errChan chan error
var ErrOperationTimeout = errors.New("Timeout while performing task")

var (
	db     connection.SimpleDatabase
	router *httprouter.Router
)

func toInterface(students []model.Student) []interface{} {
	interfaces := make([]interface{}, len(students))
	for i, v := range students {
		interfaces[i] = v
	}
	return interfaces
}

func studentsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	defer req.Body.Close()
	var response *model.ResponseMessage
	students, err := selectAllStudent(req.Context())
	if err != nil {
		response = &model.ResponseMessage{
			Data:       nil,
			Errors:     []string{err.Error()},
			HttpStatus: http.StatusInternalServerError,
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response.Serialize())
		return
	}

	response = &model.ResponseMessage{
		Data:       toInterface(students),
		Errors:     []string{""},
		HttpStatus: 200,
	}

	w.WriteHeader(200)
	w.Write(response.Serialize())

}

func newStudentHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	student := &model.Student{}

	err = json.Unmarshal(body, student)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(400)
		return
	}

	timeoutContext, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()
	if student != nil {
		go insertStudent(timeoutContext, *student)
		select {
		case <-timeoutContext.Done():
			response := &model.ResponseMessage{
				Data:       nil,
				Errors:     []string{ErrOperationTimeout.Error()},
				HttpStatus: http.StatusRequestTimeout,
			}
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write(response.Serialize())
			return
		case err := <-errChan:
			if err != nil {
				response := &model.ResponseMessage{
					Data:       nil,
					Errors:     []string{err.Error()},
					HttpStatus: 500,
				}
				w.WriteHeader(500)
				w.Write(response.Serialize())
				return
			}
		}
	}
	//response code 200
	response := &model.ResponseMessage{
		Data:       toInterface([]model.Student{*student}),
		Errors:     []string{""},
		HttpStatus: 200,
	}

	w.WriteHeader(200)
	w.Write(response.Serialize())
}

func insertStudent(ctx context.Context, student model.Student) {
	insertQuery := fmt.Sprintf("INSERT INTO student (name, age, class) VALUES ('%v', %v, '%v')", student.Name, student.Age, student.Class)
	if err := db.Insert(ctx, insertQuery); err != nil {
		errChan <- err
	}

	errChan <- nil
}

func selectAllStudent(ctx context.Context) ([]model.Student, error) {
	selectQuery := "SELECT * FROM student"
	return db.Select(ctx, selectQuery)
}

func setupDB() {
	var err error
	url := os.Getenv("DB_URL")
	if url == "" {
		url = "root:rahasia@tcp(192.168.99.100:3306)/latihan"
	}

	db, err = connection.NewMySQLConnection(url)
	if err != nil {
		panic(err)
	}
}

func setupRouter() {
	router = httprouter.New()
	router.POST("/student/submit", newStudentHandler)
	router.GET("/students", studentsHandler)
}

func main() {
	errChan = make(chan error)
	setupDB()
	setupRouter()

	http.ListenAndServe(":12345", router)
}
