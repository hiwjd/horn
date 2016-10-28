package dispatcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

func push(url string, bs []byte) error {
	bodyType := "application/json"
	r := bytes.NewReader(bs)

	rsp, err := http.Post("http://"+url+"/push", bodyType, r)
	if err != nil {
		return err
	}

	rbs, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	log.Printf(" -> push return: %s \r\n", rbs)

	var v struct {
		Code int
		Msg  string
	}
	err = json.Unmarshal(rbs, &v)
	if err != nil {
		return err
	}

	if v.Code != 0 {
		return errors.New("push return fail")
	}

	return nil
}
