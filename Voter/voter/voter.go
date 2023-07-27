package voter

import (
	"encoding/json"
	"os"
	"errors"
)

type VoterItem struct {
	Id     int    `json:"id"`
	Name  string `json:"name"`
}

type DbMap map[int]VoterItem

type Voter struct {
	voterMap    DbMap
	dbFileName string
}

func New() (*Voter, error) {
	var dbFile = "./data/voter.json"
	if _, err := os.Stat(dbFile); err != nil {
		err := initDB(dbFile)
		if err != nil {
			return nil, err
		}
	}

	voter := &Voter{
		voterMap:    make(map[int]VoterItem),
		dbFileName: dbFile,
	}

	return voter, nil
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

func (t *Voter) AddItem(voter *VoterItem) error {
	err := t.loadDB()

	if err != nil {
		return errors.New("addVoter() LoadDB failed")
	}

	var id = len(t.voterMap)
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
		return nil, errors.New("GetAllItems() LoadDB failed")
	}

	v := make([]VoterItem, 0, len(t.voterMap))

	for _, value := range t.voterMap {
		v = append(v, value)
	}

	return v, nil
}

func (t *Voter) GetItem(id int) (VoterItem, error) {
	err := t.loadDB()

	if err != nil {
		return VoterItem{}, errors.New("GetItem() LoadDB failed")
	}

	if _, ok := t.voterMap[id]; ok {
		return t.voterMap[id], nil
	}

	return VoterItem{}, errors.New("Voter not found")
}

func (t *Voter) JsonToVoter(jsonString string) (VoterItem, error) {
	var voter VoterItem
	err := json.Unmarshal([]byte(jsonString), &voter)
	if err != nil {
		return VoterItem{}, err
	}

	return voter, nil
}


func (t *Voter) saveDB() error {
	var voterList []VoterItem
	for _, item := range t.voterMap {
		voterList = append(voterList, item)
	}
	data, err := json.MarshalIndent(voterList, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(t.dbFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (t *Voter) loadDB() error {
	data, err := os.ReadFile(t.dbFileName)
	if err != nil {
		return err
	}

	var voterList []VoterItem
	err = json.Unmarshal(data, &voterList)
	if err != nil {
		return err
	}

	for _, item := range voterList {
		t.voterMap[item.Id] = item
	}

	return nil
}
