package models

type Room struct {
	Id        int
	Title     string
	CreatorId int
}

func GetRoom(roomId int) *Room {
	return &Room{
		Id:    roomId,
		Title: "Room " + string(roomId),
	}
}
