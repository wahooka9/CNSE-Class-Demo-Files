package schema

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

