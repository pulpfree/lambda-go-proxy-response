package pres

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/pulpfree/pkgerrors"
)

func extractBody(body string) (ret map[string]interface{}, err error) {
	byt := []byte(body)
	err = json.Unmarshal(byt, &ret)
	return ret, err
}

func TestResponsePass(t *testing.T) {

	tm := time.Now()
	dataRet := "pong"
	statusRet := "success"

	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"

	res := ProxyRes(Response{
		Code:      200,
		Data:      dataRet,
		Status:    statusRet,
		Timestamp: tm.Unix(),
	}, hdrs, nil)

	body, err := extractBody(res.Body)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("got %d, want %d", res.StatusCode, 200)
	}
	if body["data"] != dataRet {
		t.Errorf("got %s, want %s", body["data"], dataRet)
	}
	if body["status"] != statusRet {
		t.Errorf("got %s, want %s", body["status"], statusRet)
	}

	// test timestamp
	var btm int64 = int64(body["timestamp"].(float64))
	if btm != tm.Unix() {
		t.Errorf("got %d, want %d", btm, tm.Unix())
	}
}

func TestResponseFail(t *testing.T) {

	tm := time.Now()
	errMsg := "Invalid and bad return"

	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"

	err := errors.New(errMsg)
	res := ProxyRes(Response{
		Timestamp: tm.Unix(),
	}, hdrs, err)

	body, err := extractBody(res.Body)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 500 {
		t.Errorf("got %d, want %d", res.StatusCode, 500)
	}
	if body["message"] != errMsg {
		t.Errorf("got %s, want %s", body["status"], errMsg)
	}

	// test timestamp
	var btm int64 = int64(body["timestamp"].(float64))
	if btm != tm.Unix() {
		t.Errorf("got %d, want %d", btm, tm.Unix())
	}
}

func TestResponseFailWithpkgErrors(t *testing.T) {

	tm := time.Now()
	errMsg := "pkg error message"

	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"

	err := &pkgerrors.StdError{Err: "", Caller: "db.GetDay", Msg: errMsg}
	res := ProxyRes(Response{
		Timestamp: tm.Unix(),
	}, hdrs, err)

	body, bErr := extractBody(res.Body)
	if bErr != nil {
		panic(bErr)
	}

	if res.StatusCode != 500 {
		t.Errorf("got %d, want %d", res.StatusCode, 500)
	}
	if body["message"] != errMsg {
		t.Errorf("got %s, want %s", body["status"], errMsg)
	}

	// test timestamp
	var btm int64 = int64(body["timestamp"].(float64))
	if btm != tm.Unix() {
		t.Errorf("got %d, want %d", btm, tm.Unix())
	}
}
