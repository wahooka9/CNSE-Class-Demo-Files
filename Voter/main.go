package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"flag"
	"os"
	"strconv"
	"sync"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"drexel.edu/voter/voter"
	"drexel.edu/voter/poll"
	"drexel.edu/voter/votes"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

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

func healthCheckHandler(c *gin.Context) {
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
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	//w.Write(healthRecordJSON)
	c.JSON(http.StatusOK, healthRecordJSON)
}

func AddVoterHandler(c *gin.Context) {
	var voterToAdd voter.VoterItem

	voterRepo, err := voter.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := c.ShouldBindJSON(&voterToAdd); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := voterRepo.AddItem(&voterToAdd); err != nil {
		fmt.Println("Error: ", err)
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusConflict)
	}

	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(voterToAdd)
	c.JSON(http.StatusOK, voterToAdd)
}


func GetVotersHandler(c *gin.Context) {
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

	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(voterList)
	c.JSON(http.StatusOK, voterList)
}

func GetVotersByIDHandler(c *gin.Context) {
		vsmutex.Lock()
		apiVotersSinglsCallsCounter++
		vsmutex.Unlock()

	idS := c.Param("id")
	voterID, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
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
	//tempMap["Votes"] = v
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(tempMap)
	c.JSON(http.StatusOK, tempMap)
}

func AddPollHandler(c *gin.Context) {
	var pollToAdd poll.PollItem

	pollRepo, err := poll.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := c.ShouldBindJSON(&pollToAdd); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := pollRepo.AddItem(&pollToAdd); err != nil {
		fmt.Println("Error: ", err)
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusConflict)
	}

	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(pollToAdd)
	c.JSON(http.StatusOK, pollToAdd)
}

func GetPollsHandler(c *gin.Context) {
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

	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(pollList)
	c.JSON(http.StatusOK, pollList)
}

func GetVotesHandler(c *gin.Context) {
	vpsmutex.Lock()
	apiVoterPollSingleCallsCounter++
	vpsmutex.Unlock()

	voteRepo, err := votes.NewVotes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	idS := c.Param("poll_id")
	pollID, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollInfo, err := voteRepo.GetPollItems(pollID)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	c.JSON(http.StatusOK, pollInfo)
}

func GetVotesFromVoterOnPollHandler(c *gin.Context) {
	vpsmutex.Lock()
	apiVoterPollSingleCallsCounter++
	vpsmutex.Unlock()

	voteRepo, err := votes.NewVotes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	idS := c.Param("id")
	voterID, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	idP := c.Param("poll_id")
	pollID, err := strconv.ParseInt(idP, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollList, err := voteRepo.GetVoterDataOnPoll(voterID, pollID)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	c.JSON(http.StatusOK, pollList)
}


func AddVoteHandler(c *gin.Context) {
	var voteToAdd votes.VoteData

	idS := c.Param("voter_id")
	voterID, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	idP := c.Param("poll_id")
	pollID, err := strconv.ParseInt(idP, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := c.ShouldBindJSON(&voteToAdd); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
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

	c.JSON(http.StatusOK, voteToAdd)
}



func main() {

	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/voters", GetVotersHandler)
	r.POST("/voters", AddVoterHandler)
	r.POST("/voters/:id", AddVoterHandler)
	r.GET("/voters/:id", GetVotersByIDHandler)

	r.GET("/polls", GetPollsHandler)
	r.POST("/polls", AddPollHandler)

	r.GET("/polls/:id", GetVotesHandler)
	r.GET("/voters/:id/polls/:poll_id", GetVotesFromVoterOnPollHandler)
	r.POST("/voters/:id/polls/:poll_id", AddVoteHandler)

	r.GET("/health", healthCheckHandler)

	//We will now show a common way to version an API and add a new
	//version of an API handler under /v2.  This new API will support
	//a path parameter to search for todos based on a status
	//v2 := r.Group("/v2")
	//v2.GET("/todo", apiHandler.ListSelectTodos)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)

}



