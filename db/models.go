package db

import (
	"time"
)

type Base struct {
	ID        string    `gorm:"primaryKey" json:"id" sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserInfo struct {
	Base
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Id           string `json:"id"`
}
type SignUpUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Player struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type JWTInfo struct {
	JWT string `json:"jwt"`
}

type Room struct {
	Room        string
	Id1         string
	Id2         string
	IsStreaming bool
}

type Point struct {
	Id    string `json:"id"`
	Point string `json:"point"`
}

type Pointer struct {
	PointerX float64 `json:"x"`
	PointerY float64 `json:"y"`
}

type Post struct {
	Room  int     `json:"room"`
	UUID  string  `json:"name"`
	Pos   Pointer `json:"pos"`
	Score int     `json:"score"`
}

type SendData struct {
	Room     int                `json:"room"`
	Question []int              `json:"question"`
	Score    map[string]int     `json:"score"`
	Pos      map[string]Pointer `json:"pos"`
}

type PostPlayer struct {
	Player1 map[string]string `json:"player1"`
	Player2 map[string]string `json:"player2"`
}

func InitSendData() *SendData {
	return &SendData{
		Question: make([]int, 10),
		Score:    make(map[string]int),
		Pos:      make(map[string]Pointer),
	}
}
