package poll

import (
	"log"
	"encoding/json"
	"errors"
	"strconv"
	"drexel.edu/poll-api/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PollItem struct {
	Id     int64    `json:"id"`
	Name  string `json:"name"`
	Selections []string `json:"selection"`
}

type DbMap map[int64]PollItem

type Poll struct {
	pollMap    DbMap
	dbFileName string
}

/////////////////
///  API
/////////////////
func (t *Poll) GetPollsHandler(c *gin.Context) {
	pollList, err := t.GetAllItems()
	
	if err != nil {
		log.Println("poll failed :", err)
	}

	c.JSON(http.StatusOK, pollList)
}

func (t *Poll) GetPollHandler(c *gin.Context) {
	idS := c.Param("id")
	pollID, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	polldata, err := t.GetItem(pollID)
	if err != nil {
		log.Println("poll failed :", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, polldata)
}

func (t *Poll) AddPollHandler(c *gin.Context) {
	var pollToAdd PollItem
	if err := c.ShouldBindJSON(&pollToAdd); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := t.AddItem(&pollToAdd); err != nil {
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusConflict)
	}

	c.JSON(http.StatusOK, pollToAdd)
}


/////////////////
/// Poll Controller
/////////////////

func New() (*Poll, error) {
	var dbFile = "./data/poll.json"

	poll := &Poll{
		pollMap:    make(map[int64]PollItem),
		dbFileName: dbFile,
	}

	return poll, nil
}

func (t *Poll) AddItem(poll *PollItem) error {
	err := t.loadDB()

	if err != nil {
		log.Println("Load DB Filed createing new:", err)
	}

	var id = int64(len(t.pollMap))

	var _, exists = t.pollMap[id]
	for ok := true; ok; ok = exists { 
		id++
		_, exists = t.pollMap[id]
	}
	
	poll.Id = id
	t.pollMap[id] = *poll

	dberr := t.saveDB()
	if dberr != nil {
		return errors.New("addpoll() saveDB failed")
	}

	return nil
}

func (t *Poll) GetAllItems() ([]PollItem, error) {
	err := t.loadDB()
	if err != nil {
		return []PollItem{}, errors.New("GetAllItems() LoadDB failed")
	}

	v := make([]PollItem, 0, len(t.pollMap))
	if len(t.pollMap) < 1 {
		log.Println("poll failed :", []PollItem{})
		return []PollItem{}, nil
	}

	for _, value := range t.pollMap {
		v = append(v, value)
	}

	return v, nil
}

func (t *Poll) GetItem(pollID int64) (PollItem, error) {
	err := t.loadDB()
	if err != nil {
		return PollItem{}, errors.New("GetAllItems() LoadDB failed")
	}

	var pollData = t.pollMap[pollID]
	if pollData.Id != 0 { 
		log.Println("Poll exist: ", t.pollMap[pollID])
		return t.pollMap[pollID], nil
	}
	return PollItem{}, errors.New("Poll does not exist")
}


func (t *Poll) JsonToPoll(jsonString string) (PollItem, error) {
	var poll PollItem
	err := json.Unmarshal([]byte(jsonString), &poll)
	if err != nil {
		return PollItem{}, err
	}

	return poll, nil
}

////////////////////
///  Repository 
/////////////////////
func (t *Poll) saveDB() error {
	var pollList []PollItem
	for _, item := range t.pollMap {
		pollList = append(pollList, item)
	}
	data, err := json.MarshalIndent(pollList, "", "  ")
	if err != nil {
		return err
	}
	repository.SetValueForKey(t.dbFileName, data)
	if err != nil {
		return err
	}
	return nil
}


func (t *Poll) loadDB() error {
	data, err := repository.GetValueForKey(t.dbFileName)
	if err != nil {
		return err
	}

	var pollList []PollItem

	err = json.Unmarshal([]byte(data), &pollList)
	if err != nil {
		return err
	}

	for _, item := range pollList {
		t.pollMap[item.Id] = item
	}

	return nil
}


