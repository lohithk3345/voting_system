package voting

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lohithk3345/voting_system/cache"
	"github.com/lohithk3345/voting_system/types"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

type UserWebsocketsManager struct {
	clients    map[types.VoterID]*Client
	broadcast  chan *BroadcastData
	join       chan *JoinData
	leave      chan *LeaveData
	createRoom chan *CreateRoom
	vote       chan *Vote
	rooms      map[string]*VotingRoom
	store      *cache.CacheService
	mu         sync.Mutex
}

// func NewWebsocketManager() *UserWebsocketsManager {
// 	return &UserWebsocketsManager{
// 		clients:    make(map[string]*Client),
// 		broadcast:  make(chan *BroadcastData, 256),
// 		join:       make(chan *JoinData, 256),
// 		leave:      make(chan *LeaveData, 512),
// 		createRoom: make(chan *CreateRoom, 256),
// 		vote:       make(chan *Vote, 1024),
// 		rooms:      make(map[string]*VotingRoom),
// 	}
// }

func Initialize(cache *cache.CacheService) (*UserWebsocketsManager, error) {
	voteData, err := cache.InitWSManager()
	if err != nil {
		log.Println("Error in Initializing")
		return nil, err
	}

	var rooms map[string]*VotingRoom = make(map[string]*VotingRoom)

	for _, data := range voteData {
		rooms[data.Keys] = InitRoom(data.Keys, data.Data)
	}

	if len(rooms) <= 0 {
		return &UserWebsocketsManager{
			clients:    make(map[string]*Client),
			broadcast:  make(chan *BroadcastData, 256),
			join:       make(chan *JoinData, 256),
			leave:      make(chan *LeaveData, 512),
			createRoom: make(chan *CreateRoom, 256),
			vote:       make(chan *Vote, 1024),
			rooms:      make(map[string]*VotingRoom),
			store:      cache,
		}, nil
	}

	return &UserWebsocketsManager{
		clients:    make(map[string]*Client),
		broadcast:  make(chan *BroadcastData, 256),
		join:       make(chan *JoinData, 256),
		leave:      make(chan *LeaveData, 512),
		createRoom: make(chan *CreateRoom, 256),
		vote:       make(chan *Vote, 1024),
		rooms:      rooms,
		store:      cache,
		mu:         sync.Mutex{},
	}, nil
}

func (m *UserWebsocketsManager) HandleWS(ctx *gin.Context) {
	// store := cache.NewCacheService()
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Error in upgrading")
		return
	}
	defer conn.Close()

	id, ok := ctx.Get("userId")

	if !ok {
		return
	}

	log.Println("TOKEN", id)

	conn.WriteMessage(websocket.TextMessage, []byte(id.(string)))

	// room, _ := CreateRandomRoom()
	client, id := NewClient(conn)

	go m.Run()

	m.Read(client)
	// for {
	// }
}

func (m *UserWebsocketsManager) Read(client *Client) {
	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		var message Message
		json.Unmarshal(msg, &message)
		log.Println(message.Command, message.Data)
		if message.Command == "JOIN" || message.Command == "join" {
			var joinData JoinData
			json.Unmarshal(message.Data, &joinData)
			log.Println("JOIN DATA", joinData, message.Data)
			joinData.Client = client
			log.Println("NIL ERROR", m.rooms, joinData)
			if m.rooms != nil {
				m.join <- &joinData
			}
			// if m.rooms[joinData.RoomId].clients[client.id] != nil {
			// 	log.Println("Already Client Present")
			// } else {
			// 	m.join <- &joinData
			// }
		}

		if message.Command == "CREATE" || message.Command == "create" {
			var voteData types.VoteData
			err := json.Unmarshal(message.Data, &voteData)
			if err != nil {
				log.Println("Error Unmarshal", err)
			}
			log.Println("JOIN DATA", voteData)
			createRoom := NewCreateRoom(voteData, JoinData{Client: client})
			log.Println("CREATE", m.rooms)
			m.createRoom <- createRoom
		}

		if message.Command == "leave" || message.Command == "LEAVE" {
			var leaveData LeaveData
			json.Unmarshal(message.Data, &leaveData)
			log.Println("Leaving", m.rooms[leaveData.RoomId].clients[client.id])
			if m.rooms[leaveData.RoomId].clients[client.id] != nil {
				m.leave <- &leaveData
			} else {
				log.Println("Already Client Left")
			}
		}

		if message.Command == "VOTE" || message.Command == "vote" {
			var vote Vote
			json.Unmarshal(message.Data, &vote)
			m.vote <- &vote
		}
		msg, _ = json.Marshal(message)
		client.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func (m *UserWebsocketsManager) Run() {
	for {
		select {
		case voteData := <-m.createRoom:
			created := CreateNewRoom(&voteData.Data, m)

			wData, err := json.Marshal(created)
			if err != nil {
				log.Println("Error Marshal In CreateRoom", err)
				continue
			}

			m.mu.Lock()
			m.store.CreateRoom(created.Id, &voteData.Data)
			m.mu.Unlock()

			log.Println("CREATED", created)

			voteData.Join.RoomId = created.Id

			m.join <- &voteData.Join

			m.broadcast <- NewBroadcastData(created.Id, []byte(wData))

			// &BroadcastData{
			// 	RoomId: created.Id,
			// 	Data:   []byte(wData),
			// }

			// log.Println("Room Created")

		case joinData := <-m.join:
			// m.rooms[joinData.RoomId].AddClient(joinData.Client)
			v := m.rooms[joinData.RoomId]
			v.AddClient(joinData.Client)

		case leaveData := <-m.leave:
			m.rooms[leaveData.RoomId].RemoveClient(leaveData.Client)

		case bData := <-m.broadcast:
			log.Println("BROADCASTING", bData)
			room := m.rooms[bData.RoomId]
			log.Println("ID", room.id)
			room.BraodcastWrite(bData.Data)

		case vote := <-m.vote:
			m.mu.Lock()
			data, err := m.store.SetVoteByRoomId(vote.RoomId, vote.Option, vote.Message)
			m.mu.Unlock()
			if err != nil {
				log.Println("Error SetVote", err)
				continue
			}

			wData, err := json.Marshal(data)
			if err != nil {
				log.Println("Error Marshal In Vote", err)
				continue
			}
			m.broadcast <- NewBroadcastData(vote.RoomId, []byte(wData))
		}
	}
	// log.Println("Exitting Run")
}
