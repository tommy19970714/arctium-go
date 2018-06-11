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

func (o *OauthTwitter) AccessToken() string {
	return o.access_token
}

func (o *OauthTwitter) AccessTokenSecret() string {
	return o.access_token_secret
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

func (t *Task) GroupId() int {
	return t.group_id
}

func (t *Task) Name() string {
	return t.task_name
}

func (t *Task) Description() string {
	return t.description
}

func (t *Task) Time() time.Time {
	return t.time
}

type Members []*Member

type Member struct {
	id         int
	group_id   int
	account_id int
	created_at time.Time
	updated_at time.Time
}

type Notifications []*Notification

type Notification struct {
	id         int
	task_id    int
	date       time.Time
	created_at time.Time
	updated_at time.Time
}

func (n *Notification) Date() time.Time {
	return n.date
}

type Group struct {
	id              int
	group_name      string
	description     string
	publish_setting string
	created_by      int
	updated_by      int
	created_at      time.Time
	updated_at      time.Time
}

func (g *Group) Id() int {
	return g.id
}

func (g *Group) GroupName() string {
	return g.group_name
}

func (g *Group) Description() string {
	return g.description
}

func (g *Group) IsPublic() bool {
	if g.publish_setting == "public" {
		return true
	}
	return false
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

func SelectNotifications(task_id int) Notifications {
	sql := fmt.Sprintf("select * from notifications where task_id = %d;", task_id)
	rows, sqlErr := DBConn.Query(sql)
	if sqlErr != nil {
		log.Fatal(sqlErr)
	}
	var notifications Notifications
	for rows.Next() {
		note := Notification{}
		err := rows.Scan(&note.id, &note.task_id, &note.date, &note.created_at, &note.updated_at)
		if err != nil {
			log.Fatal(err)
		}
		notifications = append(notifications, &note)
	}
	return notifications
}

func SelectGroup(task_id int) Group {
	sql := fmt.Sprintf("select * from groups where id = %d;", task_id)
	rows, sqlErr := DBConn.Query(sql)
	if sqlErr != nil {
		log.Fatal(sqlErr)
	}
	group := Group{}
	if rows.Next() {
		err := rows.Scan(&group.id, &group.group_name, &group.description, &group.publish_setting, &group.created_by, &group.updated_by, &group.created_at, &group.updated_at)
		if err != nil {
			log.Fatal(err)
		}
	}
	return group
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
	Connect()
	user := SelectOauthTwitter(1)
	fmt.Println(user.access_token)

	task := SelectTask(1)
	fmt.Println(task.description)

	members := SelectMembers(1)
	for _, m := range members {
		fmt.Println(m.account_id)
	}
}
