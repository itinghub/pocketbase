package book

import (
	"encoding/base64" // Add this import statement
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tests"
	"google.golang.org/protobuf/proto"
)

func TestProtobuf(t *testing.T) {
	util := TestUtil{
		urlFmt: "/api/collections/%s/records?_c=%s",
		testItems: []TestItem{
			{
				Model: "vbooks_basic",
				Req:   &ListReq{Page: 1, PerPage: 10},
				buildExpectResp: func(t *testing.T, app *tests.TestApp) proto.Message {
					records, err := app.Dao().FindRecordsByFilter("vbooks_basic", "1=1", "", 10, 0)
					if err != nil {
						t.Fatalf("error: %v", err.Error())
					}
					protoBooks := make([]*Basic, len(records))

					for i, record := range records {
						protoBooks[i] = pb_toBookBasic(record)
					}

					return &BookListResp{Page: 1, PerPage: 10, Items: protoBooks, TotalItems: uint32(len(records)), TotalPages: 1}
				},
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
	Model           string
	Req             proto.Message
	buildExpectResp func(t *testing.T, app *tests.TestApp) proto.Message
}

func (b *TestUtil) protoToBase64StrRespJson(t *testing.T, msg proto.Message) string {
	protoData, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("error: %v", err.Error())
	}
	encoded := base64.StdEncoding.EncodeToString(protoData)
	hash := map[string]interface{}{
		"_c": encoded,
	}
	jsonData, err := json.Marshal(hash)
	if err != nil {
		t.Fatalf("error: %v", err.Error())
	}
	return string(jsonData)
}
func (b *TestUtil) compareRespJsonStr(t *testing.T, expect []string, actual []string, reqtype reflect.Type) bool {
	msg1 := b.respJsonStrToProto(t, expect, reqtype)
	msg2 := b.respJsonStrToProto(t, actual, reqtype)
	return b.compareProtoMessag(t, msg1, msg2)
}

func (b *TestUtil) respJsonStrToProto(t *testing.T, jsonStr_ []string, reqtype reflect.Type) proto.Message {
	jsonStr := strings.Join(jsonStr_, "")
	var rawMap map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &rawMap)
	base64Str, ok := rawMap["_c"].(string)
	if !ok {
		t.Fatalf("error: _c not found, in: %v ", jsonStr)
	}
	converter := ProtoConverter{}
	msg, err := converter.DecodeToProto(base64Str, reqtype)
	if err != nil {
		t.Fatalf("error: decode failed in: %v", jsonStr)
	}
	return msg
}

func (b *TestUtil) compareProtoMessag(t *testing.T, expect, actual proto.Message) bool {
	converter := ProtoConverter{}
	map1, _ := converter.ProtoToMap(expect)
	map2, _ := converter.ProtoToMap(actual)
	return b.compareMaps(t, map1, map2)
}
func prettyPrintJson(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func (b *TestUtil) compareMaps(t *testing.T, map1, map2 map[string]interface{}) bool {
	t.Helper() // This marks the function as a helper, which improves test output

	if len(map1) != len(map2) {
		t.Errorf("Map lengths differ: %d vs %d \n expect:\n %v vs \n actual:\n %v", len(map1), len(map2), prettyPrintJson(map1), prettyPrintJson(map2))
		return false
	}

	for key, value1 := range map1 {
		value2, exists := map2[key]
		if !exists {
			t.Errorf("Key %s not found in second map", key)
			return false
		}
		if !reflect.DeepEqual(value1, value2) {
			t.Errorf("Values for key %s are different: %v vs %v, \n expect:\n %v \n vs \n actual:\n %v", key, value1, value2, prettyPrintJson(map1), prettyPrintJson(map2))
			return false
		}
	}

	return true
}

func (b *TestUtil) runTests(t *testing.T) {
	beforeTestFunc := func() func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
		return func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
			InitBook(app, e)
		}
	}
	for _, item := range b.testItems {
		reqUrl := b.reqUrl(item.Req)
		fmt.Println(reqUrl)

		testcase := tests.ApiScenario{
			Name:            "Test " + item.Model,
			Method:          http.MethodGet,
			Url:             reqUrl,
			BeforeTestFunc:  beforeTestFunc(),
			TestAppFactory:  appFactoryFunc(),
			ExpectedStatus:  200,
			ExpectedContent: []string{""},
			ExpectedEvents:  map[string]int{"OnRecordsListRequest": 1},
			CheckResponseContentFunc: func(t *testing.T, app *tests.TestApp, res *http.Response, jsonStr string, api tests.ApiScenario) {
				expectProto := item.buildExpectResp(t, app)
				expect := b.protoToBase64StrRespJson(t, expectProto)
				if !b.compareRespJsonStr(t, []string{expect}, []string{jsonStr}, reflect.TypeOf(expectProto)) {
					t.Errorf("response content mismatch")
				}
			},
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
		InitBook(app, e)
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
