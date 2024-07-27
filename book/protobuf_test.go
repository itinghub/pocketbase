package book_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"google.golang.org/protobuf/proto"
)

// 定义函数 serializer： 将 i interface{} 转成 JSON 或 Protobuf 格式
func serializer(i interface{}) ([]byte, error) {
	switch v := i.(type) {
	case proto.Message:
		// 如果是 Protobuf message，序列化为 Protobuf 格式
		return proto.Marshal(v)
	default:
		// 如果是其他类型，序列化为 JSON
		return json.Marshal(i)
	}
}

// 示例 Protobuf 消息结构体（假设已经生成并导入）
type MyMessage struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Age  int32  `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty"`
}

func TestProtoBuf(t *testing.T) {
	// 示例 JSON 数据
	data := map[string]interface{}{
		"name":    "Alice",
		"age":     30,
		"married": false,
		"hobbies": []string{"reading", "gaming"},
	}

	// 调用 serializer 函数（JSON）
	jsonData, err := serializer(data)
	if err != nil {
		fmt.Println("JSON Error:", err)
	} else {
		fmt.Println("JSON:", string(jsonData))
	}

	// 示例 Protobuf 数据
	message := &MyMessage{
		Name: "Bob",
		Age:  25,
	}

	// 调用 serializer 函数（Protobuf）
	protoData, err := serializer(message)
	if err != nil {
		fmt.Println("Protobuf Error:", err)
	} else {
		fmt.Println("Protobuf:", protoData)
	}
}
