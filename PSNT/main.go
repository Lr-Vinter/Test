package main

import ( //"database/sql"
	"feed/internal/dbapi"
	"fmt"
)

func main() {

	DB, _ := dbapi.InitializeDataBase("sqlite3", "test1.db")
	err := DB.UnFollow(1, 4)

	//feed, time, err := DB.Feed(2, 20, 2)
	//comm, err := DB.ViewComments(2, 1)
	//fmt.Printf("posts %v", comm)
	//err := DB.DeleteComment(4, 2)
	//err := DB.DeletePost(4, 16)
	//feed, time, err := DB.Feed(2, 20, 2)
	//fmt.Printf("posts %v", feed)
	//fmt.Println("time", time)

	fmt.Println(err)

}
