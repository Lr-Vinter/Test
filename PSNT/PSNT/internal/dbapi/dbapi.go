package dbapi

import ( //"database/sql"

	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type DataBase struct {
	db *sql.DB

	answer []Post
}

func InitializeDataBase(driverName string, sourceName string) (*DataBase, error) {
	db, err := sql.Open(driverName, sourceName)
	if err != nil {
		return nil, err
	}

	return &DataBase{db, nil}, nil
}

func (d *DataBase) PushPost(UserID int, Message string, CreatedAt int64) error {
	tx, err := d.db.Begin()

	res, err := tx.Exec(`INSERT INTO MessageObjects(Type) VALUES ('post');`)
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

	res, err := tx.Exec(`INSERT INTO MessageObjects(Type) VALUES ('comment');`)
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

	_, err = tx.Exec(`DELETE FROM Comments WHERE OwnerID = ? AND CommentID = ?;`, UserID, CommentID)
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
	_, err = tx.Exec(`DELETE FROM Posts WHERE OwnerID = ? AND PostID = ?;`, UserID, PostID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

type modifier func(*sql.Tx, int) error

func (d *DataBase) loadAnswer(result *sql.Rows, likes bool) error {
	currpost := Post{}
	objectid := 0

	var err error
	for result.Next() {
		if likes {
			err = result.Scan(&currpost.PostID, &objectid, &currpost.OwnerID, &currpost.Message, &currpost.CreatedAt, &currpost.LikeNumber)
		} else {
			err = result.Scan(&currpost.PostID, &objectid, &currpost.OwnerID, &currpost.Message, &currpost.CreatedAt)
		}
		if err != nil {
			return err
		}
		d.answer = append(d.answer, currpost)
	}

	return nil
}

func (d *DataBase) ModExcludeSeenPosts(tx *sql.Tx, UserID int) error {

	modcond := ` AS
	select * from Mods exi

	left join SeenPosts 
		ON Mods.PostID = SeenPosts.PostID AND SeenPosts.FollowerID = ` + strconv.Itoa(UserID) +
		` where SeenPosts.FollowerID is NULL;`

	_, err := tx.Exec(`create temp table Mods` + modcond)

	if err != nil {
		_, err = tx.Exec(`create temp table TempData` + modcond)
		_, err = tx.Exec(`drop table Mods;`)

		_, err = tx.Exec(`alter table TempData RENAME TO Mods;`)

	}

	return nil
}

func (d *DataBase) ModExcludeFromUser(tx *sql.Tx, UserID int) error {

	modcond := ` AS 
	select * from Mods 
	where Mods.OwnerID != ?;`

	_, err := tx.Exec(`create temp table Mods`+modcond, 2)

	if err != nil {
		_, err = tx.Exec(`create temp table TempData`+modcond, 2)

		_, err = tx.Exec(`drop table Mods;`)
		_, err = tx.Exec(`alter table TempData RENAME TO Mods;`)
	}

	return nil
}

func (d *DataBase) ModLikesCount(tx *sql.Tx, UserID int) error {

	modcond := ` AS select Mods.PostID, Mods.ObjectID, Mods.OwnerID, Mods.Message, Mods.CreatedAt, count(Likes.ObjectID) as LikeNumber from Mods 
	left join Likes ON Mods.ObjectID = Likes.ObjectID group by Mods.ObjectID`

	_, err := tx.Exec(`create temp table Mods`+modcond, 2)

	if err != nil {
		_, err = tx.Exec(`create temp table TempData`+modcond, 2)

		_, err = tx.Exec(`drop table Mods;`)
		_, err = tx.Exec(`alter table TempData RENAME TO Mods;`)
	}

	return nil
}

func (d *DataBase) LoadFromTemp(tx *sql.Tx) *sql.Rows {
	res, err := tx.Query(`select Mods.PostID, Mods.ObjectID, Mods.OwnerID, Mods.Message, Mods.CreatedAt, LikeNumber from Mods`)
	if err != nil {
		panic(err)
	}
	return res
}

//
type LogicFunc func(int, int64, int, ...modifier) ([]Post, error)

func (d *DataBase) GetFollowerPosts(UserID int, time int64, number int, modfuncs ...modifier) ([]Post, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	//
	queryParam := ` select Posts.PostID, Posts.ObjectID, Posts.OwnerID, Posts.Message, Posts.CreatedAt from Posts join Followers 
		ON Posts.OwnerID = Followers.TargetID AND Followers.SubID = ? 
		where Posts.CreatedAt > ?
		limit ?`

	var res *sql.Rows

	len := len(modfuncs)
	switch len {
	case 0:
		res, err = tx.Query(queryParam, UserID, time, number)
		d.loadAnswer(res, false)

	default:
		_, err = tx.Exec(`create temp table Mods AS`+queryParam, UserID, time, number)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		for i := range modfuncs {
			err = modfuncs[i](tx, UserID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		res = d.LoadFromTemp(tx)
		d.loadAnswer(res, true)
	}
	//_, err = tx.Exec(fmt.Sprint("drop table ", "temp.Mods"))
	tx.Commit()
	return d.answer, nil // *sql.Rows link and ?
}

// Rework
// Posts with likes from target
func (d *DataBase) GetLikedByFollowerPosts(UserID int, time int64, number int, modfuncs ...modifier) ([]Post, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	//
	queryParam := ` select distinct Posts.PostID, Posts.ObjectID, Posts.OwnerID, Posts.Message, Posts.CreatedAt from Posts left join Likes
		ON Posts.ObjectID = Likes.ObjectID
			left join Followers as LikeFromTarget
				ON Likes.UserID = LikeFromTarget.TargetID AND LikeFromTarget.SubID = ?
					where LikeFromTarget.SubID = ? AND Posts.CreatedAt > ?
					limit ?`

	var res *sql.Rows

	len := len(modfuncs)
	switch len {
	case 0:
		res, err = tx.Query(queryParam, UserID, UserID, time, number)
		d.loadAnswer(res, false)

	default:
		_, err = tx.Exec(`create temp table Mods AS`+queryParam, UserID, UserID, time, number)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		for i := range modfuncs {
			err = modfuncs[i](tx, UserID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		res = d.LoadFromTemp(tx)
		d.loadAnswer(res, true)
	}

	tx.Commit()
	return d.answer, nil
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
