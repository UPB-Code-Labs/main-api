package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	accounts_http "github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/http"
	"github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/requests"
	config_infrastructure "github.com/UPB-Code-Labs/main-api/src/config/infrastructure"
	courses_http "github.com/UPB-Code-Labs/main-api/src/courses/infrastructure/http"
	session_http "github.com/UPB-Code-Labs/main-api/src/session/infrastructure/http"
	shared_infrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

// --- Globals ---
var (
	router *gin.Engine

	registeredStudentEmail string
	registeredStudentPass  string

	registeredAdminEmail = "development.admin@gmail.com"
	registeredAdminPass  = "changeme123*/"

	registeredTeacherEmail string
	registeredTeacherPass  string
)

type GenericTestCase struct {
	Payload            map[string]interface{}
	ExpectedStatusCode int
}

// --- Setup ---
func TestMain(m *testing.M) {
	// Setup
	setupDatabase()
	defer shared_infrastructure.ClosePostgresConnection()

	setupRouter()
	setupControllers()

	registerBaseAccounts()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func setupDatabase() {
	shared_infrastructure.GetPostgresConnection()
	config_infrastructure.RunMigrations()
}

func setupRouter() {
	router = gin.Default()

}

func setupControllers() {
	group := router.Group("")
	group.Use(shared_infrastructure.ErrorHandlerMiddleware())

	session_http.StartSessionRoutes(group)
	accounts_http.StartAccountsRoutes(group)
	courses_http.StartCoursesRoutes(group)
}

func registerBaseAccounts() {
	registerBaseStudent()
	registerBaseTeacher()
}

func registerBaseStudent() {
	studentEmail := "greta.mann.2020@upb.edu.co"
	studentPassword := "greta/password/2023"

	code := RegisterStudent(requests.RegisterUserRequest{
		FullName:        "Greta Mann",
		Email:           studentEmail,
		InstitutionalId: "000123456",
		Password:        studentPassword,
	})
	if code != http.StatusCreated {
		panic("Error registering base student")
	}

	registeredStudentEmail = studentEmail
	registeredStudentPass = studentPassword
}

func registerBaseTeacher() {
	teacherEmail := "judy.arroyo.2020@upb.edu.co"
	teacherPassword := "judy/password/2023"

	code := RegisterTeacherAccount(requests.RegisterTeacherRequest{
		FullName: "Judy Arroyo",
		Email:    teacherEmail,
		Password: teacherPassword,
	})
	if code != http.StatusCreated {
		panic("Error registering base teacher")
	}

	registeredTeacherEmail = teacherEmail
	registeredTeacherPass = teacherPassword
}

// --- Helpers ---
func PrepareRequest(method, endpoint string, payload interface{}) (*httptest.ResponseRecorder, *http.Request) {
	var req *http.Request

	if payload != nil {
		payloadBytes, _ := json.Marshal(payload)
		req, _ = http.NewRequest(method, endpoint, bytes.NewReader(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, endpoint, nil)
	}

	w := httptest.NewRecorder()
	return w, req
}

func ParseJsonResponse(buffer *bytes.Buffer) map[string]interface{} {
	var response map[string]interface{}
	json.Unmarshal(buffer.Bytes(), &response)
	return response
}
