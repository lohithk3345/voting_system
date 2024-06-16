package voting

import (
	"encoding/json"

	"github.com/lohithk3345/voting_system/types"
)

// type VoteData struct {
// 	RoomName string
// 	Question string
// 	Options  []string
// }

// type VoteMessage interface {
// 	GetRoomId() string
// }

type CreateRoom struct {
	Data types.VoteData
	Join JoinData
}

func NewCreateRoom(data types.VoteData, join JoinData) *CreateRoom {
	return &CreateRoom{
		Data: data,
		Join: join,
	}
}

type Message struct {
	Command string          `json:"action"`
	Data    json.RawMessage `json:"data"`
}

type WriteMessage struct {
	Data []byte
}

type RoomCreated struct {
	Id   string
	Name string
}

type JoinData struct {
	Client *Client `json:"-"`
	RoomId string  `json:"room_id"`
}

type LeaveData struct {
	Client *Client
	RoomId string `json:"room_id"`
}

type BroadcastData struct {
	RoomId string `json:"room_id"`
	Data   []byte `json:"data"`
}

func NewBroadcastData(roomId string, data []byte) *BroadcastData {
	return &BroadcastData{
		RoomId: roomId,
		Data:   data,
	}
}

type Vote struct {
	RoomId string `json:"room_id"`
	Option string `json:"option"`
}
