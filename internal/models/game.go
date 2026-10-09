package models

type Game struct {
	AppID    int
	Title    string
	Playtime *int // can be NULL for untracked games
	Platform string
}
