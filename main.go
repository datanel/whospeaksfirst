package main

import (
	"log"
	"net/http"
)

var speakers Speakers = Speakers{
	Present: []string{
		"Angèle",
		"Arnaud",
		"David",
		"Eva",
		"Jean",
		"Mourad",
		"Pascal Vert",
		"Pascal Rouge",
		"Patoche",
		"Pierre-Etienne",
		"Vincent",
	},
	Absent: []string{},
}

func main() {
	hub := newHub(speakers)
	go hub.run()

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(hub, w, r)
	})

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
