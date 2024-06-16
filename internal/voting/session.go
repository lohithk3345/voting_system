package voting

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/lohithk3345/voting_system/cache"
	"github.com/lohithk3345/voting_system/types"
)

type VotingRoom struct {
	id      string
	clients map[types.VoterID]*Client
	name    string
	// store   *cache.CacheService
	mu sync.Mutex
}

func RandomIDGenerator() string {
	return uuid.NewString()
}

func CreateRandomRoom(store *cache.CacheService) (*VotingRoom, string) {
	id := RandomIDGenerator()
	return &VotingRoom{
		id:      id,
		clients: make(map[types.VoterID]*Client),
		name:    "New Room",
		// store:   store,
		mu: sync.Mutex{},
	}, id
}

func InitRoom(id string, voteData *types.VoteData) *VotingRoom {
	// id := RandomIDGenerator()

	log.Println("NAME", voteData.Name)

	return &VotingRoom{
		id:      id,
		clients: make(map[types.VoterID]*Client),
		name:    voteData.Name,
		// store:   store,
		mu: sync.Mutex{},
	}
}

func CreateNewRoom(voteData *types.VoteData, manager *UserWebsocketsManager) *RoomCreated {
	id := RandomIDGenerator()

	log.Println("NAME", voteData.Name)

	room := &VotingRoom{
		id:      id,
		clients: make(map[types.VoterID]*Client),
		name:    voteData.Name,
		// store:   store,
		mu: sync.Mutex{},
	}

	// store.CreateRoom(room.id, voteData)

	manager.rooms[id] = room

	return &RoomCreated{
		Id:   id,
		Name: voteData.Name,
	}
}

func (v *VotingRoom) CreateNewId() {
	v.id = RandomIDGenerator()
}

func (v *VotingRoom) AddClient(client *Client) {
	log.Println("ADDING")
	log.Println("ADDING ID", v.id)
	v.mu.Lock()
	v.clients[client.id] = client
	v.mu.Unlock()
	log.Println("SUCCESS ADDED")
}

func (v *VotingRoom) RemoveClient(client *Client) {
	log.Println("REMOVING")
	v.mu.Lock()
	delete(v.clients, client.id)
	v.mu.Unlock()
	log.Println("SUCCESS REMOVED")
}

func (v *VotingRoom) BraodcastWrite(data []byte) {
	log.Println("WRITE", data)
	for _, client := range v.clients {
		err := client.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			return
		}
	}
}

// func (v *VotingRoom) Vote(option string) *types.VoteData {
// 	voteData, err := v.store.SetVoteByRoomId(v.id, option)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}

// 	return voteData
// }
