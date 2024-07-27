package book

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
)

func InitBook(app core.App, e *echo.Echo) {
	api := e.Group("/books", BookContextMiddleware(app))
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
	return c.JSON(200, []string{"group1", "group2"})
}

func BookContextMiddleware(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{Context: c}
			return next(cc)
		}
	}
}
