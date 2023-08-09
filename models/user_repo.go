package models

type RepoUserModel struct {
	UserId   string `bson:"user_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type RepoGetUserModel struct {
	UserId   string `bson:"user_id,omitempty"`
	Username string `bson:"username,omitempty"`
}

type RepoCreateUserModel struct {
	UserId   string `bson:"user_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type RepoUpdateUserModel struct {
	Username string `bson:"username,omitempty"`
	Password string `bson:"password,omitempty"`
}
