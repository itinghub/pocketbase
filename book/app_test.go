package book_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/book"
	"github.com/pocketbase/pocketbase/tests"
	"google.golang.org/protobuf/proto"
)

func TestProtobuf(t *testing.T) {
	util := TestUtil{
		urlFmt: "/api/collections/%s/records?_c=%s",
		testItems: []TestItem{
			{
				Model: "vbooks_basic",
				Req:   &book.ListReq{Page: 1, PerPage: 10},
				// ExpectResp: &book.ListResp{},
				ExpectResp: &book.BookListResp{},
			},
		},
	}
	util.runTests(t)
}

type TestUtil struct {
	urlFmt    string
	testItems []TestItem
}

type TestItem struct {
	Model      string
	Req        proto.Message
	ExpectResp proto.Message
}

func (b *TestUtil) runTests(t *testing.T) {
	beforeTestFunc := func() func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
		return func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
			book.InitBook(app, e)
		}
	}
	for _, item := range b.testItems {
		reqUrl := b.reqUrl(item.Req)
		fmt.Println(reqUrl)

		response, err := proto.Marshal(item.ExpectResp)
		if err != nil {
			panic("error: " + err.Error())
		}

		testcase := tests.ApiScenario{
			Name:            "Test " + item.Model,
			Method:          http.MethodGet,
			Url:             reqUrl,
			BeforeTestFunc:  beforeTestFunc(),
			TestAppFactory:  appFactoryFunc(),
			ExpectedStatus:  200,
			ExpectedContent: []string{string(response) + "ssss"},
			ExpectedEvents:  map[string]int{"OnRecordsListRequest": 1},
		}
		testcase.Test(t)
	}
}

func (b *TestUtil) reqUrl(msg proto.Message) string {
	data, err := proto.Marshal(msg)
	if err != nil {
		panic("error: " + err.Error())
	}

	token := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf(b.urlFmt, "vbooks_basic", token)
}

func TestBookApp(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:           "empty data",
			Method:         http.MethodGet,
			Url:            "/book/groups",
			Body:           strings.NewReader(``),
			BeforeTestFunc: beforeTestFunc(),
			TestAppFactory: appFactoryFunc(),
			ExpectedStatus: 200,
			// ExpectedContent: []string{},
		},
		{
			Name:            "test2",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo1/records",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func beforeTestFunc() func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
	return func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
		book.InitBook(app, e)
	}
}
func appFactoryFunc() func(t *testing.T) *tests.TestApp {
	return func(t *testing.T) *tests.TestApp {
		testApp, testAppErr := tests.NewTestAppWithFlag(true)
		if testAppErr != nil {
			t.Fatalf("Failed to initialize the test app instance: %v", testAppErr)
		}
		return testApp
	}
}
