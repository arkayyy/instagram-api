package api

import (
	"fmt"
	iu "io/ioutil"
	"strconv"

	//"net/http"
	h "net/http"
	ht "net/http/httptest"
	"testing"
)

func TestDemo(t *testing.T){
	req:= ht.NewRequest(h.MethodGet,"/test",nil)
	w := ht.NewRecorder()

	Demo(w,req)
	res := w.Result()
	defer res.Body.Close()
	data,err := iu.ReadAll(res.Body)

	if err!=nil {t.Errorf("Error"+err.Error())}

	if string(data)!="Hello from routes" {
		t.Errorf("Expected: Hello from routes, and got: "+string(data))
	}

}


func TestPostsFetch(t *testing.T){
	req:= ht.NewRequest(h.MethodGet,"/posts/users/test12345",nil)
	w := ht.NewRecorder()

	got := GetPostsUser("12345",w,req)
	fmt.Print("len: ",got)

	i,er := strconv.Atoi(got)
	fmt.Print(er)
	if(i==0){
		t.Errorf("Error: No data found!")
	}

}
