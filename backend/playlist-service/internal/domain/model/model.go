package model

// Playlist represents a playlist in the system
type Playlist struct {
	ID     int64  `json:"id" db:"id" redis:"id"`
	UserID int64  `json:"user_id" db:"user_id" redis:"user_id"`
	Title  string `json:"title" db:"title" redis:"title"`
}

// PlaylistTrack represents a track in a playlist in the system
type PlaylistTrack struct {
	PlaylistID int64 `json:"id" db:"id" redis:"playlist_id"`
	TrackID    int64 `json:"track_id" db:"track_id" redis:"track_id"`
	Position   int64 `json:"position" db:"position" redis:"position"`
}

// NewPlaylist represents a new playlist to be created
type NewPlaylist struct {
	UserID int64  `json:"user_id" db:"user_id" redis:"user_id"`
	Title  string `json:"title" db:"title" redis:"title"`
}
