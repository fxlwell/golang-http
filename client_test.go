package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

type TestJson struct {
	Name string
	Age  int
}

func TestClient(t *testing.T) {
	httpUrl := fmt.Sprintf("http://127.0.0.1:%d/get?v=testget123456", testPort)
	if s, resp, err := DefaultClient.Get(httpUrl).String(); err != nil || s != "testget123456" || resp.StatusCode != http.StatusOK {
		t.Fatal(s, resp, err)
	}
	if s, resp, err := DefaultClient.Get(httpUrl).Bytes(); err != nil || string(s) != "testget123456" || resp.StatusCode != http.StatusOK {
		t.Fatal(s, resp, err)
	}

	jsonObj := &TestJson{"testJson", 18}
	b, err := json.Marshal(jsonObj)
	if err != nil {
		t.Fatal(err)
	}

	httpUrl = fmt.Sprintf("http://127.0.0.1:%d/get?v=%s", testPort, string(b))
	jsonObjOut := &TestJson{}
	if s, resp, err := DefaultClient.Get(httpUrl).JsonOBJ(jsonObjOut); err != nil || resp.StatusCode != http.StatusOK {
		t.Fatal(s, resp, err)
	} else {
		if jsonObjOut.Name != "testJson" || jsonObjOut.Age != 18 {
			t.Fatal(*jsonObjOut)
		}
	}

	httpUrl = fmt.Sprintf("http://127.0.0.1:%d/header", testPort)
	data := url.Values{}
	data.Set("v", "testget123456")
	data.Set("v1", "v1-value")
	data.Set("v2", "v2-value")

	header := http.Header{}
	header.Set("h", "testget123456")

	if s, resp, err := DefaultClient.PostForm(httpUrl, header, data).String(); err != nil || s != "testget123456" || resp.StatusCode != http.StatusOK {
		t.Fatal(s, resp, err)
	}
	if s, resp, err := DefaultClient.PostForm(httpUrl, header, data).Bytes(); err != nil || string(s) != "testget123456" || resp.StatusCode != http.StatusOK {
		t.Fatal(s, resp, err)
	}

	httpUrl = fmt.Sprintf("http://127.0.0.1:%d/post-form", testPort)
	if s, resp, err := DefaultClient.PostForm(httpUrl, header, data).String(); err != nil || s != "testget123456" || resp.StatusCode != http.StatusOK {
		t.Fatal(s, resp, err)
	}
	if s, resp, err := DefaultClient.PostForm(httpUrl, header, data).Bytes(); err != nil || string(s) != "testget123456" || resp.StatusCode != http.StatusOK {
		t.Fatal(s, resp, err)
	}

	data.Set("v", string(b))
	header.Set("h", string(b))

	jsonObjOut3 := &TestJson{}
	if s, resp, err := DefaultClient.PostForm(httpUrl, header, data).JsonOBJ(jsonObjOut3); err != nil || resp.StatusCode != http.StatusOK {
		t.Fatal(s, resp, err)
	} else {
		if jsonObjOut3.Name != "testJson" || jsonObjOut3.Age != 18 {
			t.Fatal(*jsonObjOut3)
		}
	}

	if s, resp, err := DefaultClient.PostForm(httpUrl, header, nil).String(); err != nil || resp.StatusCode != http.StatusOK {
		t.Fatal(s, resp, err)
	}

	httpUrl = fmt.Sprintf("http://127.0.0.1:%d/post-body", testPort)
	if s, resp, err := DefaultClient.Post(httpUrl, header, "testget123456").String(); err != nil || s != "testget123456" || resp.StatusCode != http.StatusOK {
		t.Fatal(s, resp, err)
	}
}

func Benchmark_ClientGet(b *testing.B) {
	httpUrl := fmt.Sprintf("http://127.0.0.1:%d/get?v=testget123456", testPort)
	for n := 0; n < b.N; n++ {
		DefaultClient.Get(httpUrl).String()
	}

}
