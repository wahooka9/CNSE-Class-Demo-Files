package main

import (
	"fmt"
	//"log"
	//"net/http"
	"flag"
	"os"
	//"strconv"
	//"sync"
	//"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	//"github.com/go-resty/resty/v2"
	//"drexel.edu/voter/voter"
	//"drexel.edu/voter/poll"
	"drexel.edu/votes-api/votes"
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

/*
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

*/

func main() {

	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	votesRepo, err := votes.NewVotes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	r.GET("/votes/:id", votesRepo.GetVotesHandler)
	r.GET("/votes/:id/polls/:poll_id", votesRepo.GetVotesFromVoterOnPollHandler)
	//r.POST("/vote", AddVoteHandler)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)

}



