package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"./mydatabase"
	"./twitter"

	"github.com/googollee/go-socket.io"
	"github.com/tommy19970714/gocron"
)

type Task struct {
	Id   uint64
	Time string
}

type Message struct {
	UserId  string `json:"user_id"`
	TaskId  string `json:"task_id"`
	Message string `json:"message"`
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("layout.html")
	if err != nil {
		panic(err)
	}
	jobs := gocron.AllJobs()
	tasks := make([]*Task, len(jobs)+1)

	for i, job := range jobs {
		tasks[i] = &Task{job.Id(), job.NextScheduledTime().String()}
	}

	err = tmpl.Execute(w, tasks)
	if err != nil {
		panic(err)
	}
}

func scheduleDM(task_id int, group_id int) {
	task := mydatabase.SelectTask(task_id)
	members := mydatabase.SelectMembers(group_id)
	for _, m := range members {
		oauth := mydatabase.SelectOauthTwitter(m.AccountId())
		user := mydatabase.SelectUser(m.AccountId())
		token := twitter.TwitterToken{AccessToken: oauth.AccessToken(), AccessTokenSecret: oauth.AccessTokenSecret()}
		text := fmt.Sprintf("%s を忘れないでね！", task.Name())
		twitter.DirectMessageWithName(token, text, user.Name())
	}
}

func scheduleReplay(task_id int, group_id int) {
	task := mydatabase.SelectTask(task_id)
	members := mydatabase.SelectMembers(group_id)
	for _, m := range members {
		user := mydatabase.SelectUser(m.AccountId())
		text := fmt.Sprintf("@%s %s を忘れないでね！", user.Name(), task.Name())
		twitter.BotTweet(text)
	}
}

func changeTaskHandler(w http.ResponseWriter, r *http.Request) {
	if param, ok := r.URL.Query()["id"]; ok {
		id, err := strconv.Atoi(param[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error")
			return
		}
		task := mydatabase.SelectTask(id)
		group := mydatabase.SelectGroup(task.GroupId())
		notifications := mydatabase.SelectNotificationsWithTask(id)
		for _, n := range notifications {
			date := n.Date().In(time.UTC)
			if group.IsPublic() {
				gocron.EveryOnlyId(uint64(id)).AtDateWithTime(date).Do(scheduleReplay, id, group.Id())
			} else {
				gocron.EveryOnlyId(uint64(id)).AtDateWithTime(date).Do(scheduleDM, id, group.Id())
			}
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "task is setted!")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error")
	}
}

func removeTaskHandler(w http.ResponseWriter, r *http.Request) {
	if param, ok := r.URL.Query()["id"]; ok {
		id, err := strconv.Atoi(param[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error")
			return
		}
		gocron.RemoveFromId(uint64(id))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "removed!")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error")
	}
}

func routineTask() {
	fmt.Println("routineTask")
	notifications := mydatabase.SelectNotificationsWithMin(60)
	for _, n := range notifications {
		id := n.TaskId()
		task := mydatabase.SelectTask(id)
		group := mydatabase.SelectGroup(task.GroupId())
		date := n.Date().In(time.UTC)
		gocron.RemoveFromId(uint64(id))
		if group.IsPublic() {
			gocron.EveryOnlyId(uint64(id)).AtDateWithTime(date).Do(scheduleReplay, id, group.Id())
		} else {
			gocron.EveryOnlyId(uint64(id)).AtDateWithTime(date).Do(scheduleDM, id, group.Id())
		}
	}
}

func setupSocket() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.On("chat join", func(msg string) {
			so.Join(msg)
			log.Println("joined")
		})
		so.On("chat message", func(msg string) {
			var message Message
			err := json.Unmarshal([]byte(msg), &message)
			if err == nil {
				so.BroadcastTo(message.TaskId, "chat message", msg)
			} else {
				log.Println(err)
			}
		})
		so.On("chat message with ack", func(msg string) string {
			return msg
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	http.Handle("/socket.io/", server)
	log.Println("open socket")
}

func main() {
	mydatabase.Connect()
	twitter.SetupTwitter()
	gocron.EveryWithId(0, 1).Hour().Do(routineTask)
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/change", changeTaskHandler)
	http.HandleFunc("/remove", removeTaskHandler)
	setupSocket()
	gocron.Start()
	routineTask()
	http.ListenAndServe(":1955", nil)
}
