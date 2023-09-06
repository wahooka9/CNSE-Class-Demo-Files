package main

import (
	"fmt"
	"log"
	"net/http"
	"flag"
	"os"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"drexel.edu/voter-api/voter"
	//"drexel.edu/voter-api/schema"
	//"drexel.edu/voter/poll"
	//"drexel.edu/votes-api/votes"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 2080, "Default Port")

	flag.Parse()
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
/*
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
*/
	tempMap := make(map[string]interface{})
	tempMap["VoterInfo"] = voterInto
//	tempMap["Votes"] = voterVotes

	c.JSON(http.StatusOK, tempMap)
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


	r.GET("/voters", voterRepo.GetVotersHandler)
	r.POST("/voters", voterRepo.AddVoterHandler)
	r.GET("/voters/:id", GetVotersByIDHandler)


	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)

}



