package binding

import (
	"encoding/json"
	"errors"
	"fmt"
	validator "github.com/shisan1379/msgo/Validator"
	"net/http"
	"reflect"
)

type jsonBinding struct {
	DisallowUnknownFields bool
	IsValidate            bool
}

func (j jsonBinding) Bind(request *http.Request, data any) error {
	body := request.Body
	if request == nil || body == nil {
		return errors.New("invalid request")
	}

	decoder := json.NewDecoder(body)

	//json 参数中 存在，但是结构体中不存在
	//如果想要实现json参数中有的属性，但是对应的结构体没有，报错，也就是检查结构体是否有效
	if j.DisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if j.IsValidate {
		err := validateParam(data, decoder)
		if err != nil {
			return err
		}
	} else {
		err := decoder.Decode(data)
		if err != nil {
			return err
		}
		return validate(data)
	}
	return nil
}
func validateParam(data any, decoder *json.Decoder) error {
	//解析为map，并根据map 中的key 做对比
	//判断类型为 结构体  才能解析为 map

	rVal := reflect.ValueOf(data)

	//是否 为指针
	if rVal.Kind() != reflect.Pointer {
		return errors.New("data is not a pointer")
	}
	elem := rVal.Elem().Interface()

	of := reflect.ValueOf(elem)
	switch of.Kind() {
	case reflect.Struct:
		//将 json 解析为 map
		mapVal := map[string]interface{}{}
		decoder.Decode(&mapVal)
		for i := 0; i < of.NumField(); i++ {
			field := of.Type().Field(i)
			tag := field.Tag.Get("json")
			value := mapVal[tag]
			if value == nil {
				return errors.New(fmt.Sprintf("filed [%s] is not exist", tag))
			}
		}
		//对 map 进行序列化
		marshal, _ := json.Marshal(mapVal)
		// 对 map 进行反序列化，赋值给 data
		_ = json.Unmarshal(marshal, data)

	case reflect.Slice, reflect.Array:
		elem := of.Type().Elem()
		elemType := elem.Kind()
		if elemType == reflect.Struct {
			return checkParamSlice(elem, data, decoder)
		}
	default:
		err := decoder.Decode(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkParamSlice(elem reflect.Type, data any, decoder *json.Decoder) error {
	mapData := make([]map[string]interface{}, 0)
	_ = decoder.Decode(&mapData)
	if len(mapData) <= 0 {
		return nil
	}
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		required := field.Tag.Get("msgo")
		tag := field.Tag.Get("json")
		value := mapData[0][tag]
		if value == nil && required == "required" {
			return errors.New(fmt.Sprintf("filed [%s] is required", tag))
		}
	}
	if data != nil {
		marshal, _ := json.Marshal(mapData)
		_ = json.Unmarshal(marshal, data)
	}
	return nil
}

func validate(obj any) error {
	return validator.Validator.ValidateStruct(obj)
}

func (j jsonBinding) Name() string {

	return "jsonBinding"
}
