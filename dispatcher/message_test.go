package dispatcher

import (
	"encoding/json"
	"testing"
)

func TestImageUnmarshal(t *testing.T) {
	str := `
    {
        "t":{
            "0":123123,
            "1":123125,
            "2":123126
        },
        "mid":"mid0001",
        "from":{
            "id":"uid1",
            "name":"un1"
        },
        "chat":{
            "id":"chat1"
        },
        "image":{
            "src":"src1",
            "width":20,
            "height":20,
            "size":100
        }
    }`
	t.Log(str)

	var v MessageImage
	err := json.Unmarshal([]byte(str), &v)
	if err != nil {
		t.Error(err)
	}

	t.Log(v)
}

func TestTextUnmarshal(t *testing.T) {
	str := `
    {
        "t":{
            "0":123123,
            "1":123125,
            "2":123126
        },
        "mid":"mid0001",
        "from":{
            "id":"uid1",
            "name":"un1"
        },
        "chat":{
            "id":"chat1"
        },
        "text":"你好"
    }`
	t.Log(str)

	var v MessageText
	err := json.Unmarshal([]byte(str), &v)
	if err != nil {
		t.Error(err)
	}

	t.Log(v)
}
