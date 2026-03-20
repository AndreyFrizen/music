package model

type TrackCollection struct {
	TrackID int64 `json:"track_id" db:"track_id" redis:"track_id"`
	UserID  int64 `json:"user_id" db:"user_id" redis:"user_id"`
}

type AlbumCollection struct {
	AlbumID int64 `json:"album_id" db:"album_id" redis:"album_id"`
	UserID  int64 `json:"user_id" db:"user_id" redis:"user_id"`
}

type ArtistCollection struct {
	ArtistID int64 `json:"artist_id" db:"artist_id" redis:"artist_id"`
	UserID   int64 `json:"user_id" db:"user_id" redis:"user_id"`
}

type PlaylistCollection struct {
	PlaylistID int64 `json:"playlist_id" db:"playlist_id" redis:"playlist_id"`
	UserID     int64 `json:"user_id" db:"user_id" redis:"user_id"`
}
