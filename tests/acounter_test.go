/*
 * Created: 2018.01
 * Author: Maxim Molchanov <m.molchanov@vonmo.com>
 */

package tests

import (
	"encoding/json"
	"github.com/Vonmo/relgo/core"
	"github.com/Vonmo/relgo/lib"
	"github.com/Vonmo/relgo/services"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func TestIncrease(t *testing.T) {
	core.WaitCore()
	core.WaitService(services.Acounter)

	var v int
	var err error

	key := lib.GetUUIDv4()

	incr(key)
	v, err = get(key)
	if err != nil || v != 1 {
		t.Fatal("wrong val")
	}

	incr(key)
	v, err = get(key)
	if err != nil || v != 2 {
		t.Fatal("wrong val")
	}

	incr(key)
	v, err = get(key)
	if err != nil || v != 3 {
		t.Fatal("wrong val")
	}
}

func TestDecrease(t *testing.T) {
	core.WaitCore()
	core.WaitService(services.Acounter)

	var v int
	var err error

	key := lib.GetUUIDv4()

	incr(key)
	incr(key)
	incr(key)
	v, err = get(key)
	if err != nil || v != 3 {
		t.Fatal("wrong val")
	}

	decr(key)
	v, err = get(key)
	if err != nil || v != 2 {
		t.Fatal("wrong val")
	}

	decr(key)
	v, err = get(key)
	if err != nil || v != 1 {
		t.Fatal("wrong val")
	}
}

func TestReset(t *testing.T) {
	core.WaitCore()
	core.WaitService(services.Acounter)

	var v int
	var err error

	key := lib.GetUUIDv4()

	incr(key)
	incr(key)
	incr(key)
	v, err = get(key)
	if err != nil || v != 3 {
		t.Fatal("wrong val")
	}

	reset(key)
	v, err = get(key)
	if err != nil || v != 0 {
		t.Fatal("wrong val")
	}
}

//=================================================================================================
// helpers
//=================================================================================================

func incr(key string) error {
	_, err := sendRequest("/api/increase/" + key)
	return err
}

func decr(key string) error {
	_, err := sendRequest("/api/decrease/" + key)
	return err
}

func reset(key string) error {
	_, err := sendRequest("/api/reset/" + key)
	return err
}

func get(key string) (res int, err error) {
	resp, err := sendRequest("/api/get/" + key)
	if err != nil {
		return res, err
	}
	res = int(int64(resp["val"].(float64)))
	return res, err
}

func sendRequest(url string) (res map[string]interface{}, err error) {
	resp, err := http.Get(baseUrl() + url)
	if err != nil {
		return res, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	res = map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return res, err
	}
	return res, err
}

func baseUrl() string {
	return "http://127.0.0.1:" + strconv.Itoa(core.System.Config.Services.Acounter.Socket.Port)
}
