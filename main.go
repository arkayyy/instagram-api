package main

import (
	//PRE-DEFINED PACKAGE IMPORTS:-
	"fmt"
	"net/http"
	s "sync"

	//CUSTOM PACKAGE IMPORTS:-
	conf "main/config"
	router "main/router"
)

var mutex s.Mutex


func main() {
  mutex.Lock()
  router.Start()
  
  fmt.Println("Server running at 5000")
  http.ListenAndServe(conf.Port, nil)
  

  mutex.Unlock()
}

