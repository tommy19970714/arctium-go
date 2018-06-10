package mydatabase

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	DBConn *sql.DB
)

type OauthTwitter struct {
	id                  int
	account_id          int
	access_token        string
	access_token_secret string
	created_at          time.Time
	updated_at          time.Time
}

type Tasks []*Task

type Task struct {
	id          int
	group_id    int
	task_name   string
	description string
	time        time.Time
	created_by  int
	updated_by  int
	created_at  time.Time
	updated_at  time.Time
}

type Members []*Member

type Member struct {
	id         int
	group_id   int
	account_id int
	created_at time.Time
	updated_at time.Time
}

func SelectOauthTwitter(account_id int) OauthTwitter {
	sql := fmt.Sprintf("select * from oauthtwitters where account_id = %d;", account_id)
	rows, sqlErr := DBConn.Query(sql)
	if sqlErr != nil {
		log.Fatal(sqlErr)
	}
	oauth := OauthTwitter{}
	if rows.Next() {
		err := rows.Scan(&oauth.id, &oauth.account_id, &oauth.access_token, &oauth.access_token_secret, &oauth.created_at, &oauth.updated_at)
		if err != nil {
			log.Fatal(err)
		}
	}
	return oauth
}

func SelectTask(task_id int) Task {
	sql := fmt.Sprintf("select * from tasks where id = %d;", task_id)
	rows, sqlErr := DBConn.Query(sql)
	if sqlErr != nil {
		log.Fatal(sqlErr)
	}
	task := Task{}
	if rows.Next() {
		err := rows.Scan(&task.id, &task.group_id, &task.task_name, &task.description, &task.time, &task.created_by, &task.updated_by, &task.created_at, &task.updated_at)
		if err != nil {
			log.Fatal(err)
		}
	}
	return task
}

func SelectMembers(group_id int) Members {
	sql := fmt.Sprintf("select * from members where group_id = %d;", group_id)
	rows, sqlErr := DBConn.Query(sql)
	if sqlErr != nil {
		log.Fatal(sqlErr)
	}
	var members Members
	for rows.Next() {
		member := Member{}
		err := rows.Scan(&member.id, &member.group_id, &member.account_id, &member.created_at, &member.updated_at)
		if err != nil {
			log.Fatal(err)
		}
		members = append(members, &member)
	}
	return members
}

func Connect() {
	var myEnv map[string]string
	myEnv, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connStr := fmt.Sprintf("sslmode=disable dbname=arctium_development host=db port=5432 user=postgres password=%s", myEnv["POSTGRESQL_DATABASE_PASSWORD"])
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	DBConn = db
}

func unitTest() {
	connect()
	user := SelectOauthTwitter(1)
	fmt.Println(user.access_token)

	task := SelectTask(1)
	fmt.Println(task.description)

	members := SelectMembers(1)
	for _, m := range members {
		fmt.Println(m.account_id)
	}
}
