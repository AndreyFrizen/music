package model

// Track represents a track in the system
type Track struct {
	ID       int64  `json:"id" db:"id" redis:"id"`
	Title    string `json:"title" db:"title" redis:"title"`
	Duration int64  `json:"duration" db:"duration" redis:"duration"`
	AudioURL string `json:"audio_url" db:"audio_url" redis:"audio_url"`
	ArtistID int64  `json:"artist_id" db:"artist_id" redis:"artist_id"`
	AlbumID  int64  `json:"album_id" db:"album_id" redis:"album_id"`
}
