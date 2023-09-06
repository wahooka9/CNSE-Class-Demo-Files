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
	//"drexel.edu/voter-api/voter"
	//"drexel.edu/voter-api/schema"
	"drexel.edu/poll-api/poll"
	//"drexel.edu/votes-api/votes"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 3080, "Default Port")

	flag.Parse()
}



func main() {

	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	pollRepo, err := poll.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r.GET("/polls", pollRepo.GetPollsHandler)
	r.GET("/polls/:id", pollRepo.GetPollHandler)
	r.POST("/polls", pollRepo.AddPollHandler)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)

}



