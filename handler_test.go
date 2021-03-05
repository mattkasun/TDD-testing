package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	name        string
	method      string
	url         string
	formData    url.Values
	contentType string
	code        int
	testBody    bool
	body        string
}

func TestDisplayLanding(t *testing.T) {
	cases := []TestCase{
		TestCase{
			name:        "Landing Page",
			method:      "GET",
			url:         "/",
			formData:    url.Values{},
			contentType: "text/html; charset=utf-8",
			code:        http.StatusOK,
			testBody:    false,
			body:        "<Title>Landing",
		},
	}
	runTest(t, cases)
}

func TestDisplayRegister(t *testing.T) {
	cases := []TestCase{
		TestCase{
			name:        "Registation Page",
			method:      "GET",
			url:         "/register",
			formData:    url.Values{},
			contentType: "text/html; charset=utf-8",
			code:        http.StatusOK,
			testBody:    true,
			body:        "Create Your Account",
		},
	}
	runTest(t, cases)
}

func TestProcessRegistration(t *testing.T) {
	cases := []TestCase{
		TestCase{
			name:   "Successful Registration",
			method: "POST",
			url:    "/register",
			formData: url.Values{
				"user": {"user"},
				"pass": {"user"},
			},
			contentType: "text/html; charset=utf-8",
			code:        http.StatusOK,
			testBody:    true,
			body:        "<Title>Login",
		},

		TestCase{
			name:   "User Taken",
			method: "POST",
			url:    "/register",
			formData: url.Values{
				"user": {"user"},
				"pass": {"user"},
			},
			contentType: "text/html; charset=utf-8",
			code:        http.StatusTemporaryRedirect,
			testBody:    false,
			body:        "<Title>Landing",
		},
		TestCase{
			name:   "Nil Password",
			method: "POST",
			url:    "/register",
			formData: url.Values{
				"user": {"user"},
				"pass": {""},
			},
			contentType: "text/html; charset=utf-8",
			code:        http.StatusTemporaryRedirect,
			testBody:    false,
			body:        "<Title>Registration",
		},
	}
	//Clear the user filestore
	Users.Delete()
	runTest(t, cases)
}

func TestDisplayLogin(t *testing.T) {
	cases := []TestCase{
		TestCase{
			name:        "Display Login Page",
			method:      "GET",
			url:         "/login",
			formData:    url.Values{},
			contentType: "text/html; charset=utf-8",
			code:        http.StatusOK,
			testBody:    false,
			body:        "<Title>Login",
		},
	}
	runTest(t, cases)
}

func TestProcessLogin(t *testing.T) {
	cases := []TestCase{
		TestCase{
			name:   "Invalid User",
			method: "POST",
			url:    "/login",
			formData: url.Values{
				"user": {"junk"},
				"pass": {"user"},
			},
			contentType: "text/html; charset=utf-8",
			code:        http.StatusUnauthorized,
			testBody:    false,
			body:        "<Title>Login",
		},
		TestCase{
			name:   "Invalid Password",
			method: "POST",
			url:    "/login",
			formData: url.Values{
				"user": {"user"},
				"pass": {"demo"},
			},
			contentType: "text/html; charset=utf-8",
			code:        http.StatusUnauthorized,
			testBody:    false,
			body:        "<Title>Login",
		},
		TestCase{
			name:   "Successful Login",
			method: "POST",
			url:    "/login",
			formData: url.Values{
				"user": {"user"},
				"pass": {"user"},
			},
			contentType: "text/html; charset=utf-8",
			code:        http.StatusOK,
			testBody:    false,
			body:        "<Title>Main",
		},
	}
	runTest(t, cases)
}

func TestDisplayMain(t *testing.T) {
	cases := []TestCase{
		TestCase{
			name:        "Display Main",
			method:      "GET",
			url:         "/main",
			formData:    url.Values{},
			contentType: "text/html; charset=utf-8",
			code:        http.StatusUnauthorized,
			testBody:    false,
			body:        "",
		},
	}
	runTest(t, cases)
}

func runTest(t *testing.T, cases []TestCase) {
	router := SetupRouter()
	gin.SetMode(gin.TestMode)
	Users.Dir = "testing"

	for _, tc := range cases {
		var req *http.Request
		assert := assert.New(t)

		resp := httptest.NewRecorder()

		if tc.method == "POST" {
			payload := tc.formData.Encode()
			req, _ = http.NewRequest(tc.method, tc.url, strings.NewReader(payload))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Add("Content-Length", strconv.Itoa(len(payload)))
		} else {
			req, _ = http.NewRequest(tc.method, tc.url, nil)
		}

		router.ServeHTTP(resp, req)
		contentType := resp.Header().Get("Content-Type")

		assert.Equal(tc.code, resp.Code, "Wrong Status Code for test %s, got %d, wanted %d", tc.name, resp.Code, tc.code)
		assert.Equal(tc.contentType, contentType, "Wrong Content-Type for %s: got %s, wanted %s", tc.name, contentType, tc.contentType)
		if tc.testBody {
			assert.Contains(resp.Body.String(), tc.body, "Wrong Body in test %s: wanted %s", tc.name, tc.body)
		}
	}
}
