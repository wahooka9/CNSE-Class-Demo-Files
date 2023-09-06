package main

import (
	"fmt"
	"log"
	"net/http"
	"flag"
	"os"
	"strconv"
	//"sync"
	//"time"
	//"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"drexel.edu/votes-api/votes"
	"drexel.edu/votes-api/schema"
)

var (
	hostFlag string
	voterAPIURL string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")
	flag.StringVar(&voterAPIURL, "voterapi", "http://localhost:1080", "Default endpoint for voter API")
	flag.Parse()
}


func AddVoteHandler(c *gin.Context) {
	var voteToAdd votes.VoteData
	if err := c.ShouldBindJSON(&voteToAdd); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fmt.Println(voteToAdd)
    voterURL := os.Getenv("VOTER_URL")
    
    if voterURL == "" {
        voterURL = "host.docker.internal:2080"
    } else {
    	voterURL = "http://" + voterURL
    }
	var votetest =  voterURL + "/voters/" + strconv.Itoa(int(voteToAdd.VoterID))
	fmt.Println(voterURL)
	var pub = schema.VoterItem {}
	var apiClient = resty.New()
	_, err := apiClient.R().SetResult(&pub).Get(votetest)
	if err != nil {
		emsg := "Could not get voter from API: (" + voterURL + ")" + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": emsg})
		return
	}
	fmt.Println(pub)
	if pub.Id < 1 {
		c.JSON(http.StatusNotFound, "Voter not found")
		return
	}

    pollURL := os.Getenv("POLL_URL")
    
    if pollURL == "" {
        pollURL = "host.docker.internal:3080"
    } else {
    	pollURL = "http://" + pollURL
    }
	var pollTest =  pollURL + "/polls/" + strconv.Itoa(int(voteToAdd.PollID))
	fmt.Println(pollURL)
	var poll = schema.PollItem {}
	_, err = apiClient.R().SetResult(&poll).Get(pollTest)
	fmt.Println(poll)
	if err != nil {
		emsg := "Could not get voter from API: (" + voterURL + ")" + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": emsg})
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

	votesRepo, err := votes.NewVotes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r.GET("/votes/voter/:id", votesRepo.GetVoterVotesByID)
	
	r.GET("/votes/:id", votesRepo.GetVotesHandler)
	r.GET("/votes/:id/polls/:poll_id", votesRepo.GetVotesFromVoterOnPollHandler)
	r.POST("/vote", AddVoteHandler)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)

}



