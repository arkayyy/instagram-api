package router

import (
	//PRE-DEFINED PACKAGE IMPORTS:-
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"

	//CUSTOM PACKAGE IMPORTS:-
	api "main/routes"
	//pass "main/security"
)

var mutex sync.Mutex

//Extracting URL path specifications through regex
var U3 = regexp.MustCompile(`/(?P<c1>[a-zA-Z]+)/(?P<c2>[a-zA-Z]+)/(?P<c3>[a-zA-Z0-9]+)`)
var U2 = regexp.MustCompile(`/(?P<c1>[a-zA-Z]+)/(?P<c2>[a-zA-Z0-9]+)`)
var U1 = regexp.MustCompile(`/(?P<c1>[a-zA-Z]+)`)


//Starting the main central routing function, which executes necessary functions as per route
func Start(){


	http.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {io.WriteString(rw,"Hello from Router!")}) //demo route for testing

	http.HandleFunc("/", Direct)

	

}


func Direct(rw http.ResponseWriter, r *http.Request){
	
	mutex.Lock()
	match1 := U1.FindStringSubmatch(r.URL.Path)
	match2 := U2.FindStringSubmatch(r.URL.Path)
	match3 := U3.FindStringSubmatch(r.URL.Path)

	//storing length of resultant arrays after matching with regex's
	var l1 int = len(match1)
	var l2 int = len(match2)
	var l3 int = len(match3)

	//Since at max there can be 3 parameters embedded in the URL
	var c1 string= ""
	var c2 string= ""
	var c3 string= ""
	
	if(l1==2 && l2==0 && l3==0){
		
			c1 = match1[1]	
			
			if(c1=="users"){
				api.NewUser(rw,r)
			}else if(c1=="posts"){
				api.NewPost(rw,r)
			}else{
				// io.WriteString(rw,"404 NOT FOUND!")
				http.Error(rw,"404 NOT FOUND!", http.StatusNotFound)
			}
	
	}else if(l2==3 && l3==0){

			c1 = match2[1]

			if(c1=="users"){
				c2=match2[2]
				// fmt.Println("User data fetched: ",c2)
				// io.WriteString(rw,"User data fetched!")
				api.GetUser(c2,rw,r)
			}else if(c1=="posts"){
				c2=match2[2]
				fmt.Println("Post Data Fetched: ", c2)
				//  io.WriteString(rw,"Post data detched!")
				api.GetPost(c2,rw,r)
			}else{
				http.Error(rw,"404 NOT FOUND!", http.StatusNotFound)
			}
	
	}else if(l3==4){
		
			c1 = match3[1]
			c2 = match3[2]
			c3 = match3[3]

			if(c1=="posts"){
				if(c2=="users"){
					fmt.Println("Fetched post data for: ",c3)
					api.GetPostsUser(c3,rw,r)
				}else{
					http.Error(rw,"404 NOT FOUND!", http.StatusNotFound)
				}
			}else{
				http.Error(rw,"404 NOT FOUND!", http.StatusNotFound)
			}
	}else{
		http.Error(rw,"404 NOT FOUND!", http.StatusNotFound)
	}

	mutex.Unlock()

}
