package json

//
import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

type JsonObject struct {
	value interface{}
}

func JsonObjectFromFile(filename string) (*JsonObject, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var value = make(map[string]interface{})

	err = json.Unmarshal(data, &value)

	return &JsonObject{value: value}, err
}

func JsonObjectFromString(text string) (*JsonObject, error) {
	var value = make(map[string]interface{})

	err := json.Unmarshal(bytes.NewBufferString(text).Bytes(), &value)

	return &JsonObject{value: value}, err
}

func CreateJsonObject(b []byte) (*JsonObject, error) {
	var value = make(map[string]interface{})

	err := json.Unmarshal(b, &value)

	return &JsonObject{value: value}, err
}

func (jobj *JsonObject) Get(key string) *JsonObject {

	if jobj.value != nil {

		var m, ok = (jobj.value).(map[string]interface{})
		if ok {
			return &JsonObject{value: m[key]}
		}
	}
	return &JsonObject{value: nil}
}

func (jobj *JsonObject) Exists(key string) bool {

	if jobj.value != nil {

		var _, ok = (jobj.value).(map[string]interface{})
		return ok
	}
	return false
}

func (jobj *JsonObject) Len() int {
	if jobj.value == nil {
		return 0
	}
	var vv, ok = (jobj.value).([]interface{})
	if ok {
		return len(vv)
	}
	return 0
}

func (jobj *JsonObject) GetByIndex(index int) *JsonObject {

	if jobj.value != nil {

		var s, ok = (jobj.value).([]interface{})
		if ok && index < len(s) {
			return &JsonObject{value: s[index]}
		}
	}
	return &JsonObject{value: nil}
}

func (jobj *JsonObject) GetString(key string) (string, bool) {

	if jobj.value != nil {

		var m, ok = (jobj.value).(map[string]interface{})
		if ok {
			v := m[key]
			if v != nil {
				var s, ok = v.(string)
				if ok {
					return s, ok
				}
			}
		}
	}
	return "", false
}

func (jobj *JsonObject) GetStringByIndex(index int) string {

	if jobj.value != nil {

		var vv, ok = (jobj.value).([]interface{})
		if ok && index < len(vv) {
			v := vv[index]
			if v != nil {
				var s, ok = v.(string)
				if ok {
					return s
				}
			}
		}
	}
	return ""
}

func (jobj *JsonObject) GetInteger(key string) int {
	if jobj.value != nil {

		var mapValue, ok = (jobj.value).(map[string]interface{})
		if ok {
			val := mapValue[key]
			if val != nil {
				var v, ok = val.(int)
				if ok {
					return v
				}
			}
		}
	}
	return 0
}

func (jobj *JsonObject) GetBool(key string) (bool, bool) {
	if jobj.value != nil {

		var mapValue, ok = (jobj.value).(map[string]interface{})
		if ok {
			val := mapValue[key]
			if val != nil {
				var v, ok = val.(bool)
				if ok {
					return v, ok
				}
			}
		}
	}
	return false, false
}
