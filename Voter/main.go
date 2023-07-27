package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/gorilla/mux"
	"drexel.edu/voter/voter"
	"drexel.edu/voter/poll"
	"drexel.edu/voter/votes"
)


func AddVoterHandler(w http.ResponseWriter, r *http.Request) {
	var voterToAdd voter.VoterItem
	err := json.NewDecoder(r.Body).Decode(&voterToAdd)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	voterRepo, err := voter.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := voterRepo.AddItem(&voterToAdd); err != nil {
		fmt.Println("Error: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(voterToAdd)
}


func GetVotersHandler(w http.ResponseWriter, r *http.Request) {
	voterRepo, err := voter.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	voterList, err := voterRepo.GetAllItems()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(voterList)
}

func GetVotersByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	voterID, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	voteRepo, err := votes.NewVotes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	v, err := voteRepo.GetItem(voterID)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	voterRepo, err := voter.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	d, err := voterRepo.GetItem(voterID)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	tempMap := make(map[string]interface{})
	tempMap["VoterInfo"] = d
	tempMap["Votes"] = v
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tempMap)
}

func AddPollHandler(w http.ResponseWriter, r *http.Request) {
	var pollToAdd poll.PollItem
	err := json.NewDecoder(r.Body).Decode(&pollToAdd)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	pollRepo, err := poll.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := pollRepo.AddItem(&pollToAdd); err != nil {
		fmt.Println("Error: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pollToAdd)
}

func GetPollsHandler(w http.ResponseWriter, r *http.Request) {
	pollRepo, err := poll.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pollList, err := pollRepo.GetAllItems()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pollList)
}

func GetVotesHandler(w http.ResponseWriter, r *http.Request) {
	voteRepo, err := votes.NewVotes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pollList, err := voteRepo.GetAllItems()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pollList)
}

func AddVoteHandler(w http.ResponseWriter, r *http.Request) {
	var voteToAdd votes.VoteData
	err := json.NewDecoder(r.Body).Decode(&voteToAdd)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}


	voteRepo, err := votes.NewVotes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Adding Data")
	err = voteRepo.AddItem(voteToAdd)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("success")
}



func main() {
	router := mux.NewRouter()

	// Register the API endpoints.
	router.HandleFunc("/voters", AddVoterHandler).Methods("POST")
	// curl -X POST -H "Content-Type: application/json" -d '{"name":"John Doe","id":30}' http://localhost:8080/voters
	router.HandleFunc("/voters", GetVotersHandler).Methods("GET")
	router.HandleFunc("/voters/{id}", GetVotersByIDHandler).Methods("GET")

	router.HandleFunc("/polls", AddPollHandler).Methods("POST")
	// curl -X POST -H "Content-Type: application/json" -d '{"name":"What Toppings do you prefer?","id":30, "selection":["Bacon", "Lettuce", "Tomato"]}' http://localhost:8080/polls
	router.HandleFunc("/polls", GetPollsHandler).Methods("GET")

	router.HandleFunc("/votes", GetVotesHandler).Methods("GET")
	// {"1":{"1":"Bacon"},"3":{"1":"Tomato"}}   -  User 1 = voted in Poll "1" for "Bacon"  User 3 = voted in Poll "1" for "Tomato"
	router.HandleFunc("/vote", AddVoteHandler).Methods("POST")
	//curl -X POST -H "Content-Type: application/json" -d '{"voters_id":1,"poll_id":1, "response":"Bacon"}' http://localhost:8080/votes

	// Start the server on port 8080.
	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}