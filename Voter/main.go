package main

import (
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

	c.JSON(http.StatusOK, healthRecord)
}

func GetVotersByIDHandler(c *gin.Context) {
	voterRepo, err := voter.New()
	if err != nil {
		log.Println("Error creating voter object: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voterInto, err := voterRepo.GetVoterByID(c)
	if err != nil {
		log.Println("Error getting voter information: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return 
	}

	votesRepo, err := votes.NewVotes()
	if err != nil {
		log.Println("Error creating votes object: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	voterVotes, err := votesRepo.GetVoterVotesByID(c)
	if err != nil {
		log.Println("Error getting votes: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	tempMap := make(map[string]interface{})
	tempMap["VoterInfo"] = voterInto
	tempMap["Votes"] = voterVotes

	c.JSON(http.StatusOK, tempMap)
}

func AddVoteHandler(c *gin.Context) {
	var voteToAdd votes.VoteData
	if err := c.ShouldBindJSON(&voteToAdd); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voterRepo, err := voter.New()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	voter, err := voterRepo.GetItem(voteToAdd.VoterID)
	if voter.Id == 0 {
		log.Println("User does not exist: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollRepo, err := poll.New()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	poll, err := pollRepo.GetItem(voteToAdd.PollID)
	if err != nil  {
		log.Println("Poll does not exist: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voteRepo, err := votes.NewVotes()
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	for _, selection := range poll.Selections {
		fmt.Println("check selection: ", selection)
		if selection == voteToAdd.Selection {
			fmt.Println("selection: ", selection)
			err = voteRepo.AddItem(voteToAdd)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			c.JSON(http.StatusOK, voteToAdd)
			return
		}
	}
	c.JSON(http.StatusBadRequest, voteToAdd)
}



func main() {

	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	voterRepo, err := voter.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pollRepo, err := poll.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	votesRepo, err := votes.NewVotes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}



	r.GET("/voters", voterRepo.GetVotersHandler)
	r.POST("/voters", voterRepo.AddVoterHandler)
	r.GET("/voters/:id", GetVotersByIDHandler)

	r.GET("/polls", pollRepo.GetPollsHandler)
	r.POST("/polls", pollRepo.AddPollHandler)

	r.GET("/votes/:id", votesRepo.GetVotesHandler)
	r.GET("/votes/:id/polls/:poll_id", votesRepo.GetVotesFromVoterOnPollHandler)
	r.POST("/vote", AddVoteHandler)

	r.GET("/health", healthCheckHandler)

	//v2 := r.Group("/v2")
	//v2.GET("/todo", apiHandler.ListSelectTodos)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)

}



