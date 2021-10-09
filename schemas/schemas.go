package schemas

type PostSchema struct{
	Id string `json: "Id"`
	UserId string `json: "UserId"`
	Caption string `json: "Caption"`
	ImgUrl string `json: "ImageUrl"`
	PostTime string `json: "PostTime"`
}

type UserSchema struct{
	Id string `json: "Id"`
	Name string `json: "Name"`
	Email string `json: "Email"`
	Password string `json: "Password"`
}


