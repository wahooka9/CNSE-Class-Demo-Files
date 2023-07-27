package poll

import (
	"encoding/json"
	"os"
	"errors"
)

type PollItem struct {
	Id     int    `json:"id"`
	Name  string `json:"name"`
	Selections []string `json:"selection"`
}

type DbMap map[int]PollItem

type Poll struct {
	pollMap    DbMap
	dbFileName string
}

func New() (*Poll, error) {
	var dbFile = "./data/poll.json"
	if _, err := os.Stat(dbFile); err != nil {
		err := initDB(dbFile)
		if err != nil {
			return nil, err
		}
	}

	poll := &Poll{
		pollMap:    make(map[int]PollItem),
		dbFileName: dbFile,
	}

	return poll, nil
}

func initDB(dbFileName string) error {
	f, err := os.Create(dbFileName)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte("[]"))
	if err != nil {
		return err
	}

	f.Close()

	return nil
}

func (t *Poll) AddItem(poll *PollItem) error {
	err := t.loadDB()

	if err != nil {
		return errors.New("addpoll() LoadDB failed")
	}

	var id = len(t.pollMap)

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
		return nil, errors.New("GetAllItems() LoadDB failed")
	}

	v := make([]PollItem, 0, len(t.pollMap))

	for _, value := range t.pollMap {
		v = append(v, value)
	}

	return v, nil
}


func (t *Poll) JsonToPoll(jsonString string) (PollItem, error) {
	var poll PollItem
	err := json.Unmarshal([]byte(jsonString), &poll)
	if err != nil {
		return PollItem{}, err
	}

	return poll, nil
}


func (t *Poll) saveDB() error {
	var pollList []PollItem
	for _, item := range t.pollMap {
		pollList = append(pollList, item)
	}
	data, err := json.MarshalIndent(pollList, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(t.dbFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (t *Poll) loadDB() error {
	data, err := os.ReadFile(t.dbFileName)
	if err != nil {
		return err
	}

	var pollList []PollItem
	err = json.Unmarshal(data, &pollList)
	if err != nil {
		return err
	}

	for _, item := range pollList {
		t.pollMap[item.Id] = item
	}

	return nil
}

