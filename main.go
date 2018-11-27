package main

import (
	"encoding/json"
	"github.com/moonlightming/simple-docker-inside-webhook/commons"
	"github.com/moonlightming/simple-docker-inside-webhook/conf"
	"github.com/moonlightming/simple-docker-inside-webhook/dockerCli"
	"github.com/moonlightming/simple-docker-inside-webhook/email"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	config = conf.NewConfig()
)

func handle(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if ok := config.IsAuthentication(key); !ok {
		email.SendMail("key不正确，请确认是否本人操作。如若不是，请及时处理。")
		w.Write([]byte("NotAuthorized"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var hookReq commons.HookRequest
	json.Unmarshal(body, &hookReq)

	if err := dockercli.PullImage(hookReq.PublicAddr(), hookReq.RepoType); err != nil {
		log.Println("PullImage Fail: ", err)
		if err := email.SendMail(err.Error()); err != nil {
			log.Println("Send Email Fail: ", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	group := r.FormValue("group")
	log.Println("ServiceName: " + group + "_" + hookReq.Name)
	if err := dockercli.UpdateService(hookReq, group); err != nil {
		log.Println("UpdateService Fail: ", err)
		email.SendMail(err.Error())
		if err := email.SendMail(err.Error()); err != nil {
			log.Println("Send Email Fail: ", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("UpdateServiceError"))
		return
	}

	if err := email.SendMail("Build Successful!"); err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func main() {
	http.HandleFunc("/webhook", handle)
	log.Println("Listen:", config.BindHost)
	log.Fatal(http.ListenAndServe(config.BindHost, nil))
}
