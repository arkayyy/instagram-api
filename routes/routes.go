package api

import (
	//PRE-DEFINED PACKAGE IMPORTS:-
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//CUSTOM PACKAGE IMPORTS:-
	"main/config"
	sc "main/schemas"
  pw "main/security"
)

var mutex sync.Mutex

type Res struct{
  Err string `json:"Error"`
  Pos string `json:"SuccessMessage"`
}


func Jsonify(res *Res) []byte{
  jsonD1,er := json.Marshal(res)
  if er!=nil{ fmt.Println("ERROR: "+er.Error()) } else { fmt.Print("") }

  return jsonD1
}


//ROUTES:-

//ROUTE 1 (SET OF ROUTES): for Testing Purposes
func Demo(res http.ResponseWriter, req *http.Request) string{
  io.WriteString(res,"Hello from routes")
  return "Hello from routes"
}




//ROUTE 2: For storing a NEW POST in the DB
func NewPost(res http.ResponseWriter, req *http.Request){
      
      log.Println(req.RemoteAddr,req.Method, req.URL)

      if(req.Method=="GET"){
        erro := &Res{Err: "BAD REQUEST", Pos: ""}
        
        jsonD := Jsonify(erro)

        res.Header().Set("Content-Type","application/json")
        res.Write(jsonD)
        http.Error(res,"",http.StatusBadRequest)
        return
      }

      mutex.Lock()
      var client *mongo.Client = config.Mongo()
     // fmt.Println(reflect.TypeOf(client))

      ctx, _  := context.WithTimeout(context.Background(), 15*time.Second)

      db := config.DB
      col := config.PostsCol

      c := client.Database(db).Collection(col)
      //fmt.Println("Collection type: ", reflect.TypeOf(c))

      if err:=req.ParseForm(); err!=nil{
        io.WriteString(res,"Parsing Error: "+err.Error())
      }

      uid := req.FormValue("id")
      cn := req.FormValue("caption")
      iurl := req.FormValue("imgurl")
      pt := req.FormValue("posttime")

      // fmt.Println(uid+" "+cn+" "+iurl+" "+pt)

      if uid=="" || cn=="" || iurl=="" || pt==""{
        erro := &Res{Err: "BAD REQUEST", Pos: ""}
        
        jsonD := Jsonify(erro)

        res.Header().Set("Content-Type","application/json")
        res.Write(jsonD)
        http.Error(res,"",http.StatusBadRequest)
        return
      }

      o:= options.Count()
      cnt,err:= c.CountDocuments(ctx,bson.M{"id":uid},o)

      if err!=nil{erro := &Res{Err: err.Error(), Pos: ""} 
      jsonD := Jsonify(erro)
      res.Header().Set("Content-Type","application/json")
      res.Write(jsonD)
      return
      }

      if cnt>0 {
        uid = config.Uid()
      }

      oneDoc := sc.PostSchema{

        Id: config.Uid(),
        UserId: uid,
        Caption: cn,
        ImgUrl: iurl,
        PostTime: pt,
      }

      result,err := c.InsertOne(ctx,oneDoc)
      mutex.Unlock()

      if err!=nil{
        fmt.Println(err.Error())

        erro := &Res{Err: err.Error(), Pos: ""}
        
        jsonD := Jsonify(erro)

        res.Header().Set("Content-Type","application/json")
        res.Write(jsonD)

        //io.WriteString(res,"ERROR: "+err.Error())
        //os.Exit(1)
      }else{
        
        pres := &Res{Err: "",Pos: "SUCCESSFULLY POSTED!"}

        jsonD := Jsonify(pres)
        res.Header().Set("Content-Type","application/json")
        res.Write(jsonD)

        fmt.Println("SUCCESS POST: ",result.InsertedID)
      }

      
}







//ROUTE 3: For user creation in MongoDB
func NewUser(res http.ResponseWriter, req *http.Request){
      
      log.Println(req.RemoteAddr,req.Method, req.URL)

      if(req.Method=="GET"){
        erro := &Res{Err: "BAD REQUEST", Pos: ""}
        
        jsonD := Jsonify(erro)

        res.Header().Set("Content-Type","application/json")
        res.Write(jsonD)
        http.Error(res,"",http.StatusBadRequest)
        return
      }

      mutex.Lock()
      var client *mongo.Client = config.Mongo()
      //fmt.Println(reflect.TypeOf(client))

      ctx, _  := context.WithTimeout(context.Background(), 15*time.Second)

      db := config.DB
      col := config.UsersCol

      c := client.Database(db).Collection(col)
      //fmt.Println("Collection type: ", reflect.TypeOf(c))
      if err:=req.ParseForm(); err!=nil{
        io.WriteString(res,"Parsing Error: "+err.Error())

      }

      id := config.Uid()
      name := req.FormValue("name")
      email := req.FormValue("email")
      pass := req.FormValue("password")



      if name=="" || email=="" || pass==""{
        erro := &Res{Err: "BAD REQUEST", Pos: ""}
        
        jsonD := Jsonify(erro)

        res.Header().Set("Content-Type","application/json")
        res.Write(jsonD)
        http.Error(res,"",http.StatusBadRequest)
        return
      }


      o:= options.Count()
      cnt,err:= c.CountDocuments(ctx,bson.M{"email":email},o)

      if err!=nil{erro := &Res{Err: err.Error(), Pos: ""} 
      jsonD := Jsonify(erro)
      res.Header().Set("Content-Type","application/json")
      res.Write(jsonD)
      http.Error(res,err.Error(),http.StatusInternalServerError)
      return
      }

      if cnt>0 {
        erro := &Res{Err: "User already exists with this email!", Pos: ""} 
        jsonD := Jsonify(erro)
        res.Header().Set("Content-Type","application/json")
        res.Write(jsonD)
        return
      }

      var encp = pw.Encrypt(pass)
      
      oneDoc := sc.UserSchema{
        Id: id,
        Name: name,
        Email: email,
        Password: encp,
      }

      result, insertErr := c.InsertOne(ctx,oneDoc)

      mutex.Unlock()

      if insertErr!=nil{
        fmt.Println(insertErr)
        //os.Exit(1)
        erro := &Res{Err: insertErr.Error(), Pos: ""}
        
        jsonD := Jsonify(erro)
        res.Header().Set("Content-Type","application/json")
        res.Write(jsonD)
        http.Error(res,insertErr.Error(), http.StatusInternalServerError)
        return

      }else{
        fmt.Println( "SUCCESS USER CREATE: ", result.InsertedID)
        pres := &Res{Err: "",Pos: "SUCCESSFULLY CREATED USER!"}

        jsonD := Jsonify(pres)

        res.Header().Set("Content-Type","application/json")
        res.Write(jsonD)
      }

      //io.WriteString(res,"User Created!")
}








//ROUTE 4: GET the user details using the USER ID  
func GetUser(id string,res http.ResponseWriter, req *http.Request){


        log.Println(req.RemoteAddr,req.Method, req.URL)
        mutex.Lock()
        var client *mongo.Client = config.Mongo()
        // fmt.Println(reflect.TypeOf(client))

        //ctx, _  := context.WithTimeout(context.Background(), 15*time.Second)
        db := config.DB
        col := config.UsersCol

        c := client.Database(db).Collection(col)
        //fmt.Println("Collection type: ", reflect.TypeOf(c))

        type Fields struct{
          Id string
          Name string
          Email string
        }
        var result Fields

        err := c.FindOne(context.TODO(),bson.M{"id":id}).Decode(&result)

        mutex.Unlock()

        if err!=nil{
          //fmt.Println(reflect.TypeOf(id))
          fmt.Println("Error finding user: ", err.Error())
          erro := &Res{Err: err.Error(), Pos: ""}
        
          jsonD := Jsonify(erro)

          res.Header().Set("Content-Type","application/json")
          
          res.Write(jsonD)
          
          return

        }else{
          //io.WriteString(res,"Name: "+result.Name+"\nEmail: "+result.Email+"\nID: "+result.Id)
          jsonD,er := json.Marshal(result)
          if er!=nil{ fmt.Println("ERROR: "+er.Error())
              //io.WriteString(res,"ERROR: "+er.Error())
              http.Error(res,"INTERNAL SERVER ERROR", http.StatusInternalServerError)
              return
          }else{
              res.Header().Set("Content-Type","application/json")
              res.Write(jsonD)
          }
          
        }
}





//ROUTE 5: GET the post details using the Post ID  
func GetPost(id string,res http.ResponseWriter, req *http.Request){


        log.Println(req.RemoteAddr,req.Method, req.URL)
        mutex.Lock()
        var client *mongo.Client = config.Mongo()
        // fmt.Println(reflect.TypeOf(client))

        //ctx, _  := context.WithTimeout(context.Background(), 15*time.Second)
        db := config.DB
        col := config.PostsCol

        c := client.Database(db).Collection(col)
        //fmt.Println("Collection type: ", reflect.TypeOf(c))

        type Fields struct{
          Id string
          Caption string
          ImgUrl string
          PostTime string
        }
        var result Fields

        err := c.FindOne(context.TODO(),bson.M{"id":id}).Decode(&result)

        mutex.Unlock()

        if err!=nil{
          //fmt.Println(reflect.TypeOf(id))
          fmt.Println("Error finding user: ", err.Error())
          erro := &Res{Err: err.Error(), Pos: ""}
        
          jsonD := Jsonify(erro)

          res.Header().Set("Content-Type","application/json")
          res.Write(jsonD)
          return

        }else{
          //io.WriteString(res,"ID: "+result.Id+"\nCaption: "+result.Caption+"\nPost Time: "+result.PostTime)
          jsonD,err := json.Marshal(result)
          if err!=nil{
          //io.WriteString(res,"ERROR: "+err.Error())
          http.Error(res,"INTERNAL SERVER ERROR", http.StatusInternalServerError)
          }else{
          res.Header().Set("Content-Type","application/json")
          res.Write(jsonD)
            // fmt.Println(string(jsonD))
          }
        }

    
}







//ROUTE 6: GET the post details of a user using the USER ID  
func GetPostsUser(id string,res http.ResponseWriter, req *http.Request) string{
       
  
        log.Println(req.RemoteAddr,req.Method, req.URL)

        var skip int = 0
        var size int = 10
        var page int = 0

        if(req.URL.Query().Get("size")!="") { sz,err := strconv.Atoi(req.URL.Query().Get("size"))  
               
               if err!=nil { io.WriteString(res,"ERROR: "+err.Error()) } else { size = sz }
         }

         if(req.URL.Query().Get("skip")!="") { sk,err := strconv.Atoi(req.URL.Query().Get("skip"))  
               
               if err!=nil { 
                  fmt.Print("")
                } else {skip = sk}
         }

         if(req.URL.Query().Get("page")!="") { pg,err := strconv.Atoi(req.URL.Query().Get("page"))  
               
               if err!=nil { 
                  fmt.Print("")
                } else {page = pg}
         }

        mutex.Lock()
        var client *mongo.Client = config.Mongo()
        // fmt.Println(reflect.TypeOf(client))
        ctx, _  := context.WithTimeout(context.Background(), 15*time.Second)
        db := config.DB
        col := config.PostsCol

        c := client.Database(db).Collection(col)
        //fmt.Println("Collection type: ", reflect.TypeOf(c))

        type Fields struct{
          Id string
          UserId string
          Caption string
          ImgUrl string
          PostTime string
        }

        var results []*Fields

        findOptions := options.Find()

        
        cur,err := c.Find(ctx,bson.M{"userid":id},findOptions)

        if err!=nil{
            fmt.Println("id", id)

            fmt.Println("Error finding user: ", err.Error())
            erro := &Res{Err: err.Error(), Pos: ""}
          
            jsonD := Jsonify(erro)

            res.Header().Set("Content-Type","application/json")
            res.Write(jsonD)
            return "FAILURE"

          }else{
            fmt.Print("")
          }

        
          for cur.Next(context.TODO()){
            var elem Fields
            err:=cur.Decode(&elem)
            if err!=nil{
                  fmt.Println("id", id)
              
                  fmt.Println("Error decoding: ", err.Error())
                  erro := &Res{Err: err.Error(), Pos: ""}
              
                  jsonD := Jsonify(erro)

                  res.Header().Set("Content-Type","application/json")
                  res.Write(jsonD)
                  return "FAILURE"
              }else{
                
                results = append(results, &elem)

              }
          }
        
        mutex.Unlock()

        //fmt.Println(len(results))

        cur.Close(context.TODO())

        type Ops struct{
          Id string `json:"User ID,omitempty"`
          Posts []*Fields `json:"Posts,omitempty"`
        }

        

        //arr := config.Pagination(skip,size,page)
        if(page==0) {size = len(results)}

        if(page>0){
          skip = (len(results)/size)*(page-1)
        }
        
        
        var limit int = skip+size
        var start int = skip

        if skip+size>len(results) { limit = len(results) }
          
        if skip>len(results) {start = len(results)}

        

        // if(len(arr)>0){start = arr[0]
        // limit = arr[1]}
        

        oput := &Ops{Id: id, Posts: results[start:limit]}

        jsonD, err := json.Marshal(oput)

        if err!=nil{
          erro := &Res{Err: err.Error(), Pos: ""}

          jsonD1 := Jsonify(erro)
          
          res.Header().Set("Content-Type","application/json")
          res.Write(jsonD1)
          return "FAILURE"
        } else{

          res.Header().Set("Content-Type","application/json")
          res.Write(jsonD)
        }

        return strconv.Itoa(len(results))

}



