package votes

import (
	"encoding/json"
	"os"
	//"fmt"
	"errors"
)

type VoteData struct {
	VoterID int `json:"voter_id"`
	PollID int `json:"poll_id"`
	Selection string `json:"response"`
}


type ResponsesVoterData map[string][]int
type PollResponsesData map[int]ResponsesVoterData

type VoterPollResponsData map[int]string
type VoterPollData map[int]VoterPollResponsData


type HasVotedData map[int][]int 

type VotesData struct {
	VoterVotedData    HasVotedData `json:"has_vote_data"`
	FullPollResultsData PollResponsesData `json:"poll_data"`
	FullVoterResultsData VoterPollData `json:"vote_data"`
	DBFileName string `json:"filename"`
}



func contains(s []int, e int) bool {
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
		 	t.VoterVotedData[voter.PollID] = []int{}
		}

		_, exists = t.FullPollResultsData[voter.PollID]
		if exists == false { 
			t.FullPollResultsData = make(PollResponsesData)
			t.FullPollResultsData[voter.PollID] = map[string][]int{}
		}

		_, exists = t.FullPollResultsData[voter.PollID][voter.Selection]
		if exists == false { 
		 	t.FullPollResultsData[voter.PollID][voter.Selection] = []int{}
		}

		_, exists = t.FullVoterResultsData[voter.VoterID]
		if exists == false { 
		 	t.FullVoterResultsData[voter.VoterID] = map[int]string{}
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

func (t *VotesData) GetPollItems(id int) (ResponsesVoterData, error) {
	err := t.loadDB()

	if err != nil {
		return ResponsesVoterData{}, errors.New("GetItem() LoadDB failed")
	}

	if _, ok := t.FullPollResultsData[id]; ok {
		return t.FullPollResultsData[id], nil
	}

	return ResponsesVoterData{}, errors.New("Voter not found")
}

func (t *VotesData) GetVoterDataOnPoll(voter_id int, poll_id int) (string, error) {
	err := t.loadDB()

	if err != nil {
		return "", errors.New("GetItem() LoadDB failed")
	}

	if _, ok := t.FullVoterResultsData[voter_id][poll_id]; ok {
		return t.FullVoterResultsData[voter_id][poll_id], nil
	}

	return "", errors.New("Voter not found")
}




func (t *VotesData) GetItem(id int) (VoterPollResponsData, error) {
	err := t.loadDB()

	if err != nil {
		return VoterPollResponsData{}, errors.New("GetItem() LoadDB failed")
	}

	if _, ok := t.FullVoterResultsData[id]; ok {
		return t.FullVoterResultsData[id], nil
	}

	return VoterPollResponsData{}, errors.New("Voter not found")
}


func NewVotes() (*VotesData, error) {
	var dbFile = "./data/votes.json"
	if _, err := os.Stat(dbFile); err != nil {
		err := initDB(dbFile)
		if err != nil {
			return nil, err
		}
	}

	votes := &VotesData{
		VoterVotedData:    make(HasVotedData),
		FullPollResultsData: make(PollResponsesData),
		FullVoterResultsData: make(VoterPollData),
		DBFileName: dbFile,
	}

	return votes, nil
}




func initDB(dbFileName string) error {
	f, err := os.Create(dbFileName)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte("{}"))
	if err != nil {
		return err
	}

	f.Close()

	return nil
}


func (t *VotesData) saveDB() error {
	temp := VotesData(*t)
	data, err := json.MarshalIndent(temp, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(t.DBFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (t *VotesData) loadDB() error {
	data, err := os.ReadFile(t.DBFileName)
	
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, t)
	if err != nil {
		return err
	}

	return nil
}
