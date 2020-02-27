package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "goy/back/controller/vote"
    "google.golang.org/appengine"
)

func main() {
    
    // 路由與controller的對應配置
    router := mux.NewRouter()
    router.HandleFunc("/vote", vote.GetVote).Methods("GET")
    router.HandleFunc("/vote/info", vote.GetVoteInfo).Methods("GET")

    http.Handle("/", router)
    appengine.Main()
}
