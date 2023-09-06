package voter

import (
	"log"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"drexel.edu/voter-api/repository"
)

type VoterItem struct {
	Id     int64    `json:"id"`
	Name  string `json:"name"`
}

type DbMap map[int64]VoterItem

type Voter struct {
	voterMap    DbMap
	dbFileName string
}

func New() (*Voter, error) {
	var dbFile = "./data/voter.json"
	voter := &Voter{
		voterMap:    make(map[int64]VoterItem),
		dbFileName: dbFile,
	}

	return voter, nil
}

//////////////////////
///  API 
//////////////////////

func (t *Voter) GetVotersHandler(c *gin.Context) {
	voterList, err := t.GetAllItems()
	if err != nil {
		log.Println("Error fetching voters: ", err)
		c.JSON(http.StatusConflict, voterList)
		return
	}

	c.JSON(http.StatusOK, voterList)
}

func (t *Voter) AddVoterHandler(c *gin.Context) {
	var voterToAdd VoterItem

	if err := c.ShouldBindJSON(&voterToAdd); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := t.AddItem(&voterToAdd); err != nil {
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusConflict)
	}

	c.JSON(http.StatusOK, voterToAdd)
}

func (t *Voter) GetVoterByID(c *gin.Context) (VoterItem, error) {
	idS := c.Param("id")
	voterID, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		return VoterItem{}, err
	}

	voterInfo, err := t.GetItem(voterID)
	if err != nil {
		log.Println("Error fetching voter: ", err)
		return VoterItem{}, err
	}
	return voterInfo, nil
}

/////////////////////
///  Controller Methods
/////////////////////

func (t *Voter) AddItem(voter *VoterItem) error {
	err := t.loadDB()

	if err != nil {
		log.Println("addVoter() LoadDB failed: ", err)
		t.voterMap = DbMap{}
	}

	var id = int64(len(t.voterMap))
	var _, exists = t.voterMap[id]
	for ok := true; ok; ok = exists { 
		id++
		_, exists = t.voterMap[id]
	}

	if (voter.Id > 0) {
		id = voter.Id
		var _, exists = t.voterMap[voter.Id]
		if true == exists { 
			return errors.New("Voter Exists, addVoter() failed")
		}
	}

	voter.Id = id
	t.voterMap[id] = *voter

	dberr := t.saveDB()
	if dberr != nil {
		return errors.New("addVoter() saveDB failed")
	}

	return nil
}

func (t *Voter) GetAllItems() ([]VoterItem, error) {
	err := t.loadDB()
	if err != nil {
		return []VoterItem{}, errors.New("GetAllItems() LoadDB failed")
	}

	v := make([]VoterItem, 0, len(t.voterMap))

	for _, value := range t.voterMap {
		v = append(v, value)
	}
	return v, nil
}

func (t *Voter) GetItem(id int64) (VoterItem, error) {
	err := t.loadDB()

	if err != nil {
		return VoterItem{}, errors.New("GetItem() LoadDB failed")
	}
	if _, ok := t.voterMap[id]; ok {
		return t.voterMap[id], nil
	}
	return VoterItem{}, nil
}

func (t *Voter) JsonToVoter(jsonString string) (VoterItem, error) {
	var voter VoterItem
	err := json.Unmarshal([]byte(jsonString), &voter)
	if err != nil {
		return VoterItem{}, err
	}

	return voter, nil
}

////////////////
/// Repository 
////////////////
func (t *Voter) saveDB() error {
	var voterList []VoterItem
	for _, item := range t.voterMap {
		voterList = append(voterList, item)
	}
	data, err := json.MarshalIndent(voterList, "", "  ")
	if err != nil {
		return err
	}
	repository.SetValueForKey(t.dbFileName, data)
	if err != nil {
		return err
	}
	return nil
}


func (t *Voter) loadDB() error {
	data, err := repository.GetValueForKey(t.dbFileName)
	if err != nil {
		return err
	}

	var voterList []VoterItem

	err = json.Unmarshal([]byte(data), &voterList)
	if err != nil {
		return err
	}

	for _, item := range voterList {
		t.voterMap[item.Id] = item
	}

	return nil
}

	
