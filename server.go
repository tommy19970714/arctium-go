package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"./mydatabase"
	"./twitter"

	"github.com/tommy19970714/gocron"
)

type Task struct {
	Id   uint64
	Time string
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("layout.html")
	if err != nil {
		panic(err)
	}
	jobs := gocron.AllJobs()
	tasks := make([]*Task, len(jobs)+1)

	for i, job := range jobs {
		tasks[i] = &Task{job.Id(), job.RunTime().String()}
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
		oauth := mydatabase.SelectOauthTwitter(m.AccountId())
		user := mydatabase.SelectUser(m.AccountId())
		token := twitter.TwitterToken{AccessToken: oauth.AccessToken(), AccessTokenSecret: oauth.AccessTokenSecret()}
		text := fmt.Sprintf("@%s %s を忘れないでね！", user.Name(), task.Name())
		twitter.Tweet(token, text)
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
		layout := "2006-01-02 15:04:05"
		group := mydatabase.SelectGroup(id)
		notifications := mydatabase.SelectNotifications(id)
		for _, n := range notifications {
			dateStr := n.Date().Format(layout)
			if group.IsPublic() {
				gocron.EveryOnlyId(uint64(id)).AtDate(dateStr).Do(scheduleReplay, id, group.Id())
			} else {
				gocron.EveryOnlyId(uint64(id)).AtDate(dateStr).Do(scheduleDM, id, group.Id())
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

func main() {
	mydatabase.Connect()
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/change", changeTaskHandler)
	http.HandleFunc("/remove", removeTaskHandler)
	http.ListenAndServe(":1955", nil)
	<-gocron.Start()
}
