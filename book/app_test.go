package book_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/book"
	"github.com/pocketbase/pocketbase/tests"
)

func TestBookApp(t *testing.T) {
	t.Parallel()
	appFactory := func() func(t *testing.T) *tests.TestApp {
		return func(t *testing.T) *tests.TestApp {
			testApp, testAppErr := tests.NewTestAppWithFlag(true)
			if testAppErr != nil {
				t.Fatalf("Failed to initialize the test app instance: %v", testAppErr)
			}
			return testApp
		}
	}
	beforeTestFunc := func() func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
		return func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
			book.InitBook(app, e)
		}
	}

	scenarios := []tests.ApiScenario{
		{
			Name:           "empty data",
			Method:         http.MethodGet,
			Url:            "/book/groups",
			Body:           strings.NewReader(``),
			BeforeTestFunc: beforeTestFunc(),
			TestAppFactory: appFactory(),
			ExpectedStatus: 200,
			// ExpectedContent: []string{},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
