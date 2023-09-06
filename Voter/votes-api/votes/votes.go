package votes

import (
	"log"
	//"os"
	"fmt"
	"encoding/json"
	"errors"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"drexel.edu/votes-api/repository"
	"drexel.edu/votes-api/schema"
)

type VoteData struct {
	VoterID int64 `json:"voter_id"`
	PollID int64 `json:"poll_id"`
	Selection string `json:"response"`
}


type ResponsesVoterData map[string][]int64
type PollResponsesData map[int64]ResponsesVoterData

type VoterPollResponsData map[int64]string
type VoterPollData map[int64]VoterPollResponsData


type HasVotedData map[int64][]int64 

type VotesData struct {
	VoterVotedData    HasVotedData `json:"has_vote_data"`
	FullPollResultsData PollResponsesData `json:"poll_data"`
	FullVoterResultsData VoterPollData `json:"vote_data"`
	DBFileName string `json:"filename"`
}


/////////////////
///  API
////////////////

func (t *VotesData) GetVoterVotesByID(c *gin.Context) (VoterPollResponsData, error) {
	idS := c.Param("id")
	voterID, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return nil, err
	}

	voterInfo, err := t.GetItem(voterID)
	if err != nil {
		log.Println("Error fetching voters votes: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return nil, err
	}

	var test = "http://localhost:2080" + "/voters"
	var pub []schema.VoterItem
	var apiClient = resty.New()
	_, err = apiClient.R().SetResult(&pub).Get(test)
	if err != nil {
		emsg := "Could not get publication from API: (" + test + ")" + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": emsg})
		return nil, err
	}

	return voterInfo, err
}

func (t *VotesData) GetVotesHandler(c *gin.Context) {
	idS := c.Param("id")
	pollID, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollInfo, err := t.GetPollItems(pollID)
	if err != nil {
		fmt.Println("Error: ", err)
	}

 //    voterURL := os.Getenv("VOTER_URL")
 //    fmt.Println(voterURL)
 //    if voterURL == "" {
 //        voterURL = "host.docker.internal:2080"
 //    } else {
 //    	voterURL = "http://" + voterURL
 //    }
	// var votetest =  voterURL + "/voters"
	// var pub = [] schema.VoterItem {}
	// fmt.Println("list: ", pub)
	// var apiClient = resty.New()
	// _, err = apiClient.R().SetResult(&pub).Get(votetest)
	// fmt.Println("list: ", pub)
	// if err != nil {
	// 	emsg := "Could not get publication from API: (" + voterURL + ")" + err.Error()
	// 	c.JSON(http.StatusNotFound, gin.H{"error": emsg})
	// 	return
	// }

	c.JSON(http.StatusOK, pollInfo)
}

func (t *VotesData) GetVotesFromVoterOnPollHandler(c *gin.Context) {
	idS := c.Param("id")
	_, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
/*
	idP := c.Param("poll_id")
	pollID, err := strconv.ParseInt(idP, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollList, err := t.GetVoterDataOnPoll(voterID, pollID)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	tempMap := make(map[string]interface{})
	tempMap["voter_id"] = idS
	tempMap["poll_id"] = idP
	tempMap["response"] = pollList
*/
	c.JSON(http.StatusOK, "")
}

//////////////
///  Votes Controller
//////////////

func contains(s []int64, e int64) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func (t *VotesData) AddItem(voter VoteData) error {
	err := t.loadDB()
	if err != nil {
		return errors.New("addVote() LoadDB failed")
	}

	// some of these are redundant checks 
	if ! contains(t.VoterVotedData[voter.PollID], voter.VoterID) {

		var _, exists = t.VoterVotedData[voter.PollID]
		if exists == false { 
		 	t.VoterVotedData[voter.PollID] = []int64{}
		}

		_, exists = t.FullPollResultsData[voter.PollID]
		if exists == false { 
			t.FullPollResultsData = make(PollResponsesData)
			t.FullPollResultsData[voter.PollID] = map[string][]int64{}
		}

		_, exists = t.FullPollResultsData[voter.PollID][voter.Selection]
		if exists == false { 
		 	t.FullPollResultsData[voter.PollID][voter.Selection] = []int64{}
		}

		_, exists = t.FullVoterResultsData[voter.VoterID]
		if exists == false { 
		 	t.FullVoterResultsData[voter.VoterID] = map[int64]string{}
		}

		_, exists = t.FullVoterResultsData[voter.VoterID][voter.PollID]
		if exists == false { 
		 	t.FullVoterResultsData[voter.VoterID][voter.PollID] = voter.Selection
		}

		t.VoterVotedData[voter.PollID] = append(t.VoterVotedData[voter.PollID], voter.VoterID)
		t.FullPollResultsData[voter.PollID][voter.Selection] = append(t.FullPollResultsData[voter.PollID][voter.Selection], voter.VoterID)
		
	}

	dberr := t.saveDB()
	if dberr != nil {
		return errors.New("addVoter() saveDB failed")
	}

	return nil
}


func (t *VotesData) GetAllItems() (PollResponsesData, error) {
	err := t.loadDB()
	if err != nil {
		return PollResponsesData{}, errors.New("GetAllItems() LoadDB failed")
	}

	return t.FullPollResultsData, nil
}

func (t *VotesData) GetPollItems(id int64) (ResponsesVoterData, error) {
	err := t.loadDB()

	if err != nil {
		return ResponsesVoterData{}, errors.New("GetItem() LoadDB failed")
	}

	if _, ok := t.FullPollResultsData[id]; ok {
		return t.FullPollResultsData[id], nil
	}

	return ResponsesVoterData{}, errors.New("Voter not found")
}

func (t *VotesData) GetVoterDataOnPoll(voter_id int64, poll_id int64) (string, error) {
	//err := t.loadDB()
/*
	if err != nil {
		return "", errors.New("GetItem() LoadDB failed")
	}

	if _, ok := t.FullVoterResultsData[voter_id][poll_id]; ok {
		return t.FullVoterResultsData[voter_id][poll_id], nil
	}
*/
	return "", errors.New("Voter not found")
}




func (t *VotesData) GetItem(id int64) (VoterPollResponsData, error) {
	err := t.loadDB()
	if err != nil {
		return VoterPollResponsData{}, errors.New("GetItem() LoadDB failed")
	}

	if _, ok := t.FullVoterResultsData[id]; ok {
		return t.FullVoterResultsData[id], nil
	}

	return VoterPollResponsData{}, nil
}


func NewVotes() (*VotesData, error) {
	var dbFile = "./data/votes.json"

	votes := &VotesData{
		VoterVotedData:    make(HasVotedData),
		FullPollResultsData: make(PollResponsesData),
		FullVoterResultsData: make(VoterPollData),
		DBFileName: dbFile,
	}

	return votes, nil
}

/////////////////
///  Repository
////////////////

func (t *VotesData) saveDB() error {
	temp := VotesData(*t)
	fmt.Println(temp)
	data, err := json.MarshalIndent(temp, "", "  ")
	if err != nil {
		log.Println("Error marshaling Json: ", err)
		return err
	}
fmt.Println(data)
	repository.SetValueForKey(t.DBFileName, data)
	if err != nil {
		log.Println("Error setting value in redis: ", err)
		return err
	}
	return nil
}


func (t *VotesData) loadDB() error {
	data, err := repository.GetValueForKey(t.DBFileName)
	if err != nil {
		log.Println("Error fetching redis values: ", err)
		return err
	}

	err = json.Unmarshal([]byte(data), t)
	if err != nil {
		log.Println("Error unmarsheling data: ", err)
		t, err = NewVotes()
		//return err
	}

	return nil
}


