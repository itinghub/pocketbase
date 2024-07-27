package book

import (
	"github.com/labstack/echo/v5"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type CustomContext struct {
	echo.Context
}

func (c *CustomContext) JSON(code int, i interface{}) error {
	if c.Request().Header.Get("Accept") == "application/x-protobuf" {
		return c.ProtobufResponse(code, i)
	}
	return c.Context.JSON(code, i)
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
