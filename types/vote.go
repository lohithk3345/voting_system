package types

// import "github.com/lohithk3345/voting_system/internal/repositories/voting"

type VoteData struct {
	Name     string            `json:"name"`
	Question string            `json:"question"`
	Options  map[string]uint64 `json:"options"`
}

type InitVoteData struct {
	Data *VoteData
	Keys string
}

// type Message struct {
// 	Command string `json:"action"`
// 	Data    []byte `json:"data"`
// }

// type WriteMessage struct {
// 	Data []byte
// }

// type RoomCreated struct {
// 	Id   string
// 	Name string
// }

// type JoinData struct {
// 	Client *voting.Client
// 	RoomId string `json:"room_id"`
// }

// type LeaveData struct {
// 	Client *Client
// 	RoomId string `json:"room_id"`
// }

// type BroadcastData struct {
// 	RoomId string `json:"room_id"`
// 	Data   []byte `json:"data"`
// }

// type Vote struct {
// 	RoomId string `json:"room_id"`
// 	Option string `json:"option"`
// }
