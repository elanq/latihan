package model

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	name string
	Age  int
}
type Student struct {
	ID    int
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Class string `json:"class"`
}

type ResponseMessage struct {
	Data       []interface{} `json:"data"`
	Errors     []string      `json:"errors"`
	HttpStatus int           `json:"http_status"`
}

func (r *ResponseMessage) ToStudents() []Student {
	students := make([]Student, 0)
	for _, v := range r.Data {
		body, err := json.Marshal(v)
		if err != nil {
			continue
		}
		var student Student
		err = json.Unmarshal(body, &student)
		if err != nil {
			fmt.Println("not ok")
			continue
		}
		students = append(students, student)

	}

	return students
}

func (r *ResponseMessage) Serialize() []byte {
	raw, err := json.Marshal(r)
	if err != nil {
		return []byte("{}")
	}
	return raw
}
