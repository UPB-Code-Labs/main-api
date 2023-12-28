package integration

import "net/http"

func GetSupportedLanguages(cookie *http.Cookie) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("GET", "/api/v1/languages", nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)
	return ParseJsonResponse(w.Body), w.Code
}
