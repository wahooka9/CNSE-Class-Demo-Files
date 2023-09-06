package schema

type DbMapVoter map[int64]VoterItem

type VoterItem struct {
	Id     int64    `json:"id"`
	Name  string `json:"name"`
}

type Voter struct {
	voterMap    DbMapVoter
	dbFileName string
}

type PollItem struct {
	Id     int64    `json:"id"`
	Name  string `json:"name"`
	Selections []string `json:"selection"`
}

type DbMapPoll map[int64]PollItem

type Poll struct {
	pollMap    DbMapPoll
	dbFileName string
}