package custom_models

type Token struct {
	Username  string `json:"username" bson:"username"`
	Token     string `json:"token" bson:"token"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
}

type User struct {
	ID          *string  `json:"_id" bson:"_id"`
	Username    *string  `json:"username"`
	Phonenumber *float64 `json:"phonenumber"`
	Password    *string  `json:"password" bson:"password"`
}
