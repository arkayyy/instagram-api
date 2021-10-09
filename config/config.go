package config

import (
	//PRE-DEFINED PACKAGES IMPORT:-
	"context"
	"fmt"
	r "crypto/rand"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Port = ":5000"

var DB string = "Appointy"
var PostsCol = "UserPosts"
var UsersCol = "Users"

var UidSize = 6

var Pkey = "passphrasewhichneedstobe32bytes!"


//Initialization function for establising link to MongoDB
func Mongo() (*mongo.Client) {
	
	clientOptions := options.Client().ApplyURI("mongodb+srv://admin:admin@cluster0.rgwkm.mongodb.net/Appointy?retryWrites=true&w=majority")
	//fmt.Println("Client optom type: ", reflect.TypeOf(clientOptions))
	client,err := mongo.Connect(context.TODO(),clientOptions)

	if err!=nil {
		fmt.Println("ERROR", err)
		//os.Exit(1)
	}

	//fmt.Println(reflect.TypeOf(client))
	return client
}



func Uid() string{
	b := make([]byte,UidSize)
	_, err := r.Read(b)
	
	if err!=nil {fmt.Println("ERROR: ", err.Error())} else {fmt.Print("")}

	return fmt.Sprintf("%X",b[:])
}









