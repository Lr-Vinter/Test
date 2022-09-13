package main

import ( //"database/sql"
	service "feed/internal/Service"
	"feed/internal/dbapi"
	"fmt"
)

func main() {

	DB, _ := dbapi.InitializeDataBase("sqlite3", "test1.db")

	posts, err := DB.GetLikedByFollowerPosts(4, 3, 100)

	S := service.NewService(DB, 3)
	S.RegisterLogicFunc((*dbapi.LogicFunc)(&DB.GetFollowerPosts))
	S.GetFeed(2, 2, 20)
	//posts, err := DB.GetFollowerPosts(3, 2, 20)

	//posts, err := DB.GetLikedByFollowerPosts(4, 100, DB.ModLikesCount)
	//posts, err := DB.GetFollowerPosts(3, 100, DB.ModExcludeSeenPosts, DB.ModExcludeFromUser, DB.ModLikesCount)

	//posts, err := DB.GetFollowerPosts(3, 100, DB.ModExcludeFromUser, DB.ModExcludeSeenPosts)
	//posts, err := DB.GetFollowerPosts(3, 100)

	fmt.Println(posts)

	fmt.Println(err)

}
