package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/latihan/model"
)

func main() {
	buffer := &model.ResponseMessage{}

	url := "http://localhost:12345/students"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error at ", err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("status not OK")
		return

	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error at ", err.Error())
		return
	}

	json.Unmarshal(body, buffer)

	students := buffer.ToStudents()
	fmt.Println(students)
}
