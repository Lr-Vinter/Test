package dbapi

import ( //"database/sql"

	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DataBase struct {
	db *sql.DB
}

func InitializeDataBase(driverName string, sourceName string) (*DataBase, error) {
	db, err := sql.Open(driverName, sourceName)
	if err != nil {
		return nil, err
	}

	return &DataBase{db}, nil
}

func (d *DataBase) PushPost(UserID int, Message string, CreatedAt int64) error {
	tx, err := d.db.Begin()
	objectsinsert, err := tx.Prepare(`

		INSERT INTO MessageObjects(Type)
		VALUES ('post');

		`)
	if err != nil {
		tx.Rollback()
		return err
	}

	res, err := objectsinsert.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	ObjectId, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	insertpost, err := tx.Prepare(`
		
		INSERT INTO Posts(ObjectID, OwnerID, Message, CreatedAt) 
		VALUES (?, ?, ?, ?);
		
	`)
	_, err = insertpost.Exec(ObjectId, UserID, Message, CreatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (d *DataBase) PushComment(UserID int, Message string, PostID int, CreatedAt int64) error {
	tx, err := d.db.Begin()
	objectsinsert, err := tx.Prepare(`

		INSERT INTO MessageObjects(Type)
		VALUES ('comment');

		`)
	if err != nil {
		tx.Rollback()
		return err
	}

	res, err := objectsinsert.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	ObjectId, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	insertcomment, err := tx.Prepare(`
		
		INSERT INTO Comments(ObjectID, PostID, OwnerID, Message, CreatedAt) 
		VALUES (?, ?, ?, ?, ?);
		
	`)
	_, err = insertcomment.Exec(ObjectId, PostID, UserID, Message, CreatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Delete Comment
func (d *DataBase) DeleteComment(UserID int, CommentID int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	// Temp Table
	stmt, err := tx.Prepare(`CREATE TEMP TABLE IF NOT EXISTS Variables (Name TEXT PRIMARY KEY, Value INTEGER);`)
	if err != nil {
		tx.Rollback()
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}
	// find ObjectID
	stmt, err = tx.Prepare(`INSERT INTO Variables(Name, Value) Values ('object', (SELECT ObjectID from Comments where OwnerID = ? AND CommentID = ?));`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(UserID, CommentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Deleting Comments
	stmt, err = tx.Prepare(`DELETE FROM Comments WHERE OwnerID = ? AND CommentID = ?;`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(UserID, CommentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	//
	stmt, err = tx.Prepare(`DELETE FROM MessageObjects WHERE ObjectID = (SELECT Value from Variables LIMIT 1)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}
	//
	stmt, err = tx.Prepare(`DROP TABLE Variables`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Delete Post
func (d *DataBase) DeletePost(UserID int, PostID int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	// Temp Table
	stmt, err := tx.Prepare(`CREATE TEMP TABLE IF NOT EXISTS Variables (Name TEXT PRIMARY KEY, Value INTEGER);`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}
	// find ObjectID
	stmt, err = tx.Prepare(`INSERT INTO Variables(Name, Value) Values ('object', (SELECT ObjectID from Posts where OwnerID = ? AND PostID = ?));`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(UserID, PostID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Deleting Comments
	stmt, err = tx.Prepare(`DELETE FROM Posts WHERE OwnerID = ? AND PostID = ?;`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(UserID, PostID)
	if err != nil {
		tx.Rollback()
		return err
	}
	//
	stmt, err = tx.Prepare(`DELETE FROM MessageObjects WHERE ObjectID = (SELECT Value from Variables LIMIT 1)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}
	//
	stmt, err = tx.Prepare(`DROP TABLE Variables`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// May be need transaction
func (d *DataBase) Feed(UserID int, feedLength int, time int) ([]Post, int64, error) {
	stmt, err := d.db.Prepare(`
	WITH Answer AS 
	(
		select Posts.PostID, Posts.ObjectID, Posts.OwnerID, Message, Posts.CreatedAt from Posts
			left join Followers 
				ON Posts.OwnerID = Followers.TargetID AND followers.SubID = ?
	
			left join Likes 
				ON Posts.ObjectID = Likes.ObjectID
	
			left join Followers as LikeFromTarget 
				ON Likes.UserID = LikeFromTarget.TargetID AND LikeFromTarget.SubID = ?
	
		where (followers.SubID = ? OR LikeFromTarget.SubID = ?) AND Posts.CreatedAt > ?
	
		order by Posts.CreatedAt 
		limit ?
	)
	
	select  Answer.PostID, Answer.OwnerID, Answer.Message, Answer.CreatedAt, 
			count(Likes.ObjectID) as LikeNumber from Answer
				left join Likes ON Answer.ObjectID = Likes.ObjectID
	
	group by Answer.ObjectID
	order by Answer.CreatedAt
	`)

	if err != nil {
		return nil, 0, err
	}
	//
	t := Post{} // no alternatives to check "no row case" without this or doing 1-row query in cycle (?)
	checknorows := stmt.QueryRow(UserID, UserID, UserID, UserID, time, feedLength)
	if err != nil {
		return nil, 0, err
	}
	err = checknorows.Scan(&t.PostID, &t.OwnerID, &t.Message, &t.CreatedAt, &t.LikeNumber)
	if err == sql.ErrNoRows {
		return nil, 0, err
	}
	//
	getfeed, err := stmt.Query(UserID, UserID, UserID, UserID, time, feedLength)
	if err != nil {
		fmt.Println("hello")
		return nil, 0, err
	}

	var feed []Post
	currpost := Post{}
	for getfeed.Next() {
		err = getfeed.Scan(&currpost.PostID, &currpost.OwnerID, &currpost.Message, &currpost.CreatedAt, &currpost.LikeNumber)
		feed = append(feed, currpost)
	}

	lastPostTime := feed[len(feed)-1].CreatedAt
	return feed, lastPostTime, nil
}

func (d *DataBase) ViewComments(UserID int, PostID int) ([]Post, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(`
	WITH Answer AS 
	(
		select Posts.PostID from Posts
			left join Followers 
				ON Posts.OwnerID = Followers.TargetID AND followers.SubID = ? AND Posts.PostID = ?
	)
	select Comments.CommentID, Comments.OwnerID, Comments.Message, Comments.CreatedAt, count(Likes.ObjectID) as LikeNumber from Answer 
		join Comments 
			ON Comments.PostID = Answer.PostID
		left join Likes 
			ON Comments.ObjectID = Likes.ObjectID
	group by Comments.ObjectID
	`)
	defer stmt.Close()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	t := Post{}
	checknorows := stmt.QueryRow(UserID, PostID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = checknorows.Scan(&t.PostID, &t.OwnerID, &t.Message, &t.CreatedAt, &t.LikeNumber)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return nil, err
	}

	getcomments, err := stmt.Query(UserID, PostID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var comments []Post
	currpost := Post{}
	for getcomments.Next() {
		err = getcomments.Scan(&currpost.PostID, &currpost.OwnerID, &currpost.Message, &currpost.CreatedAt, &currpost.LikeNumber)
		comments = append(comments, currpost)
	}

	tx.Commit()
	return comments, nil
}

func (d *DataBase) Follow(UserID int, TargetID int, CreatedAt int64) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO Followers(SubID, TargetID, CreatedAt) Values (?, ?, ?);`)

	defer stmt.Close()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(UserID, TargetID, CreatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (d *DataBase) UnFollow(UserID int, TargetID int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`DELETE FROM Followers WHERE SubID = ? AND TargetID = ?`)

	defer stmt.Close()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(UserID, TargetID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
