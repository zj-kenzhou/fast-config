package fastconfig

import "testing"

func TestCamelCaseToCenterLine(t *testing.T) {
	str := "aaaBaaaCCaaa"
	res := camelCaseToCenterLine(str)
	if "aaa-baaa-c-caaa" != res {
		t.Error("convert err")
	}
	t.Log(res)
}

func TestCenterLineToCamelCase(t *testing.T) {
	str := "aaa-baaa-c-caaa"
	res := centerLineToCamelCase(str)
	if "aaaBaaaCCaaa" != res {
		t.Error("convert err")
	}
	t.Log(res)
}

func TestGetString(t *testing.T) {
	res := GetString("main.testValue1", "de")
	if res != "aa" {
		t.Error("value is not aa")
	}
	notRes := GetString("main.testValue1not", "de")
	if notRes != "de" {
		t.Error("value is not de")
	}
	if GetString("main.test-value1", "de") != "aa" {
		t.Error("value is not aa")
	}
	if GetString("main.testValue2", "de") != "bb" {
		t.Error("value is not aa")
	}
}

func TestGetInt(t *testing.T) {
	if GetInt("main.int-value", -1) != 1 {
		t.Error("value is not 1")
	}
}

func TestGetBool(t *testing.T) {
	if GetBool("main.bool-value", false) != true {
		t.Error("value is not true")
	}
}

func TestGetValue(t *testing.T) {
	res := GetValue("main.testMap", map[string]any{})
	t.Log(res)
}

func TestGetValueList(t *testing.T) {
	res := GetValue("main.testList", []any{})
	t.Log(res)
}

// 在env 中覆盖配置
func TestGetStringEnv(t *testing.T) {
	if GetString("main.test-value1", "de") != "cc" {
		t.Error("value is not cc")
	}
}

// app_env test
func TestGetStringAppEnv(t *testing.T) {
	if GetString("main.test-value1", "de") != "dd" {
		t.Error("value is not dd")
	}
	if GetString("test.value", "de") != "test" {
		t.Error("value is not test")
	}
}
