package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"github.com/gorilla/mux"
	"drexel.edu/voter/voter"
	"drexel.edu/voter/poll"
	"drexel.edu/voter/votes"
)

var apiCallsCounter int
var mutex sync.Mutex

var apiVotersCallsCounter int
var vmutex sync.Mutex

var apiVotersSinglsCallsCounter int
var vsmutex sync.Mutex

var apiPollsCallsCounter int
var pmutex sync.Mutex

var apiVoterPollSingleCallsCounter int
var vpsmutex sync.Mutex

func countAPICalls(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// make thread safe
		mutex.Lock()
		apiCallsCounter++
		mutex.Unlock()
		next.ServeHTTP(w, r)
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	count := apiCallsCounter
	mutex.Unlock()

	// Create a health record indicating that the API is healthy.
	healthRecord := map[string]interface{}{
		"status":       "200",
		"timestamp":    time.Now(),
		"endpoints": []map[string]string{
			{
				"url":         "http://localhost:8080/voters",
				"description": "Retrieve a list of all voters.",
				"callCount" : strconv.Itoa(apiVotersCallsCounter),
			},{
				"url":         "http://localhost:8080/voters/{id}",
				"description": "Retrieve information on a single voter.",
				"callCount" : strconv.Itoa(apiVotersSinglsCallsCounter),
			},{
				"url":         "http://localhost:8080/polls",
				"description": "Retrieve a list of all active polls",
				"callCount" : strconv.Itoa(apiPollsCallsCounter),
			},{
				"url":         "http://localhost:8080/voters/{voter_id}/polls/{poll_id}",
				"description": "Retrieve information on a single voter in a single poll",
				"callCount" : strconv.Itoa(apiVoterPollSingleCallsCounter),
			},
		},
		"apiCalls": count, // Include the API call count in the health record.
	}

	healthRecordJSON, err := json.MarshalIndent(healthRecord, "", "  ")
	if err != nil {
		http.Error(w, "Failed to generate health record JSON.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(healthRecordJSON)
}

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

	params := mux.Vars(r)
	voterID, err := strconv.Atoi(params["id"])
	if err == nil {
		voterToAdd.Id = voterID
	}

	if err := voterRepo.AddItem(&voterToAdd); err != nil {
		fmt.Println("Error: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(voterToAdd)
}


func GetVotersHandler(w http.ResponseWriter, r *http.Request) {
		vmutex.Lock()
		apiVotersCallsCounter++
		vmutex.Unlock()
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
		vsmutex.Lock()
		apiVotersSinglsCallsCounter++
		vsmutex.Unlock()

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
		pmutex.Lock()
		apiPollsCallsCounter++
		pmutex.Unlock()
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
	vpsmutex.Lock()
	apiVoterPollSingleCallsCounter++
	vpsmutex.Unlock()

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
	params := mux.Vars(r)
	voterID, err := strconv.Atoi(params["voter_id"])
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	pollID, err := strconv.Atoi(params["poll_id"])
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&voteToAdd)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	voteToAdd.VoterID = voterID
	voteToAdd.PollID = pollID

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
	router.HandleFunc("/voters/{id}", AddVoterHandler).Methods("POST")
	// curl -X POST -H "Content-Type: application/json" -d '{"name":"John Doe","id":30}' http://localhost:8080/voters
	// curl -X POST -H "Content-Type: application/json" -d '{"name":"John Doe" }' http://localhost:8080/voters/12
	router.HandleFunc("/voters", GetVotersHandler).Methods("GET")
	router.HandleFunc("/voters/{id}", GetVotersByIDHandler).Methods("GET")
	
	router.HandleFunc("/polls", AddPollHandler).Methods("POST")
	// curl -X POST -H "Content-Type: application/json" -d '{"name":"What Toppings do you prefer?","id":30, "selection":["Bacon", "Lettuce", "Tomato"]}' http://localhost:8080/polls
	router.HandleFunc("/polls", GetPollsHandler).Methods("GET")
	// voters/:id/polls/:pollid
	//router.HandleFunc("/voters", GetVotesHandler).Methods("GET")
	router.HandleFunc("/voters/{voter_id}/polls/{poll_id}", GetVotesHandler).Methods("GET")
	// {"1":{"1":"Bacon"},"3":{"1":"Tomato"}}   -  User 1 = voted in Poll "1" for "Bacon"  User 3 = voted in Poll "1" for "Tomato"
	router.HandleFunc("/voters/{voter_id}/polls/{poll_id}", AddVoteHandler).Methods("POST")
	//curl -X POST -H "Content-Type: application/json" -d '{"voters_id":1,"poll_id":1, "response":"Bacon"}' http://localhost:8080/votes

	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	router.Use(countAPICalls)
	// Start the server on port 8080.
	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))


}



