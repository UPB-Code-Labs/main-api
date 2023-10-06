package errors

import (
	"fmt"
	"net/http"
)

type StudentAlreadyInCourse struct {
	CourseName string
}

func (err StudentAlreadyInCourse) Error() string {
	return fmt.Sprintf("Student is already in the course %s", err.CourseName)
}

func (err StudentAlreadyInCourse) StatusCode() int {
	return http.StatusConflict
}
