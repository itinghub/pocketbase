package book

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func BookContextMiddleware(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{Context: c}
			return next(cc)
		}
	}
}

type CustomContext struct {
	echo.Context
}

func (c *CustomContext) JSON(code int, i interface{}) error {
	if true || c.Request().Header.Get("Accept") == "application/x-protobuf" {
		return c.ProtobufResponse(code, i)
	}
	return c.Context.JSON(code, i)
}

func (c *CustomContext) PopulateProtoParams() bool {

	return true
}

func (c *CustomContext) ProtobufResponse(code int, i interface{}) error {
	// 获取 Protobuf 消息类型

	protoMsg, ok := i.(protoreflect.ProtoMessage)
	if !ok {
		return c.Context.JSON(code, i)
	}

	// 序列化 Protobuf
	protoData, err := proto.Marshal(protoMsg)
	if err != nil {
		return err
	}

	// 设置响应
	c.Response().Header().Set(echo.HeaderContentType, "application/x-protobuf")
	c.Response().WriteHeader(code)
	_, err = c.Response().Write(protoData)
	return err
}

type ReqRespType struct {
	Req  proto.Message
	Resp proto.Message
}

type ReqRespRouter struct {
	TypeRoute echo.Router
}

func (r *ReqRespRouter) AddPath(method, path string, reqResp ReqRespType) {
	// r.TypeRoute.Add(method, path, handler)
	// r.TypeRoute.(*echo.RouteInfo).Meta = reqResp
}
