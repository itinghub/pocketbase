package book

import (
	"reflect"

	"github.com/labstack/echo/v5"

	// "github.com/pocketbase/pocketbase/book"
	"github.com/pocketbase/pocketbase/core"
)

func InitBook(app core.App, e *echo.Echo) {
	conf := []struct {
		method string
		path   string
		req    interface{}
		resp   interface{}
	}{
		{"GET", "/api/collections/repos/records", &ListReq{}, &RepoListResp{}},
		{"GET", "/api/collections/book_groups/records", &ListReq{}, &GroupListResp{}},
		{"GET", "/api/collections/vbooks_basic/records", &ListReq{}, &BookListResp{}},
	}

	configList := []RouteItem{}

	for _, item := range conf {
		configList = append(configList, RouteItem{
			item.method,
			item.path,
			reflect.TypeOf(item.req),
			reflect.TypeOf(item.resp),
		})
	}

	router := protoRouter(configList)
	e.Use(BookContextMiddleware(app, router))
	// api := e.Group("/book", BookContextMiddleware(app, router), middleware.Gzip())
	api := e.Group("/book")
	bindBooksApi(app, api)
}

func bindBooksApi(app core.App, rg *echo.Group) {
	api := booksApi{app: app}
	rg.GET("/groups", api.GetBookGroups)
}

type booksApi struct {
	app core.App
}

func (api *booksApi) GetBookGroups(c echo.Context) error {
	dao := api.app.Dao()
	bookGroups, err := dao.FindRecordsByFilter("book_groups", "isPublished=true", "+position", 10, 0)
	if err != nil {
		return err
	}

	dao.ExpandRecords(bookGroups, []string{"books"}, nil)
	result := convertToGroupResult(bookGroups)
	return c.JSON(200, &result)
}
