package book

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const kQueryParamKey = "_c"
const kRespProtoMsgKey = "proto_type"

func BookContextMiddleware(app core.App, router ProtoRouter) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &ProtoContext{c}

			if result, matched := router.Match(cc); matched && result.Err != nil {
				c.JSON(400, result.Err)
				return result.Err
			}
			return next(cc)
		}
	}
}

type ProtoContext struct {
	echo.Context
}

func (c *ProtoContext) createProtoMessage(reqType reflect.Type) proto.Message {
	// 检查 reqType 是否为指针类型
	if reqType.Kind() != reflect.Ptr {
		panic("reqType must be a pointer")
	}

	// 获取元素类型
	elemType := reqType.Elem()

	// 创建新实例
	instance := reflect.New(elemType).Interface()

	// 将实例转换为 proto.Message
	msg, ok := instance.(proto.Message)
	if !ok {
		panic("reqType must be a proto.Message")
	}
	return msg
}

func (c *ProtoContext) PopulateQueryParam(reqType reflect.Type, respType reflect.Type) error {
	base64Param := c.QueryParam(kQueryParamKey)
	if base64Param == "" {
		return errors.New("param " + kQueryParamKey + " not found")
	}

	pbBytes, err := base64.StdEncoding.DecodeString(base64Param)
	if err != nil {
		return fmt.Errorf("base64 invalid %s:\n - %v", kQueryParamKey, err)
	}

	reqMsg := c.createProtoMessage(reqType)

	if err = proto.Unmarshal(pbBytes, reqMsg); err != nil {
		return fmt.Errorf("proto unmarshal failed %s:\n - %v", kQueryParamKey, err)
	}

	jsonObj, err := c.protoToMap(reqMsg)
	if err != nil {
		return fmt.Errorf("proto to map failed %s:\n - %v", kQueryParamKey, err)
	}

	for k, v := range jsonObj {
		c.QueryParams().Set(k, v)
	}
	respMsg := c.createProtoMessage(reqType)

	c.Set(kRespProtoMsgKey, respMsg)
	return nil
}

func (c *ProtoContext) JSON(code int, i interface{}) error {
	// protoHeader := c.Request().Header.Get("Accept") == "application/x-protobuf"
	if protoMsg, ok := c.Get(kRespProtoMsgKey).(proto.Message); ok {
		return c.ProtobufResponse(protoMsg, code, i)
	}
	return c.Context.JSON(code, i)
}

func (c *ProtoContext) ProtobufResponse(protoMsg proto.Message, code int, i interface{}) error {
	protoData, err := c.anyToProto(i, protoMsg)
	if err != nil {
		return err
	}

	// 设置响应
	c.Response().Header().Set(echo.HeaderContentType, "application/x-protobuf")
	c.Response().WriteHeader(code)
	_, err = c.Response().Write(protoData)
	return err
}

// 使用 protojson 将 protobuf 消息转换为 JSON
func (c *ProtoContext) protoToMap(msg proto.Message) (map[string]string, error) {
	marshaler := protojson.MarshalOptions{
		UseProtoNames: true,
	}
	jsonBytes, err := marshaler.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// 解析 JSON 到 map[string]interface{}
	var rawMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &rawMap)
	if err != nil {
		return nil, err
	}

	// 将 map[string]interface{} 转换为 map[string]string
	strMap := make(map[string]string)
	for k, v := range rawMap {
		strMap[k] = fmt.Sprintf("%v", v)
	}

	return strMap, nil
}

func (c *ProtoContext) anyToProto(i interface{}, protoMsg proto.Message) (byte []byte, err error) {
	// 将 i 转换为 JSON 字符串
	jsonData, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	// 将 JSON 数据反序列化为 Protobuf 消息
	options := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	err = options.Unmarshal(jsonData, protoMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to proto: %w", err)
	}
	byte, err = proto.Marshal(protoMsg)
	return
}
