package models

type Track struct {
	TrackId int64 `json:"track_id" db:"track_id" redis:"track_id"`
	UserId  int64 `json:"user_id" db:"user_id" redis:"user_id"`
}

type Album struct {
	AlbumId int64 `json:"album_id" db:"album_id" redis:"album_id"`
	UserId  int64 `json:"user_id" db:"user_id" redis:"user_id"`
}

type Artist struct {
	ArtistId int64 `json:"artist_id" db:"artist_id" redis:"artist_id"`
	UserId   int64 `json:"user_id" db:"user_id" redis:"user_id"`
}

type GetAlbumsRequest struct {
	UserId int64 `json:"user_id" db:"user_id" redis:"user_id"`
}

type GetAlbumsResponse struct {
	Albums []*Album `json:"albums"`
}

type AddAlbumRequest struct {
	UserId  int64 `json:"user_id" db:"user_id" redis:"user_id"`
	AlbumId int64 `json:"album_id" db:"album_id" redis:"album_id"`
}

type AddAlbumResponse struct {
	AlbumId int64 `json:"album_id"`
}

type RemoveAlbumRequest struct {
	UserId  int64 `json:"user_id" db:"user_id" redis:"user_id"`
	AlbumId int64 `json:"album_id" db:"album_id" redis:"album_id"`
}

type GetArtistsRequest struct {
	UserId int64 `json:"user_id" db:"user_id" redis:"user_id"`
}

type GetArtistsResponse struct {
	Artists []*Artist `json:"artists"`
}

type AddArtistRequest struct {
	UserId   int64 `json:"user_id" db:"user_id" redis:"user_id"`
	ArtistId int64 `json:"artist_id" db:"artist_id" redis:"artist_id"`
}

type AddArtistResponse struct {
	ArtistId int64 `json:"artist_id"`
}

type RemoveArtistRequest struct {
	UserId   int64 `json:"user_id" db:"user_id" redis:"user_id"`
	ArtistId int64 `json:"artist_id" db:"artist_id" redis:"artist_id"`
}

type GetTracksRequest struct {
	UserId int64 `json:"user_id" db:"user_id" redis:"user_id"`
}

type GetTracksResponse struct {
	Tracks []*Track `json:"tracks"`
}

type AddTrackRequest struct {
	UserId  int64 `json:"user_id" db:"user_id" redis:"user_id"`
	TrackId int64 `json:"track_id" db:"track_id" redis:"track_id"`
}

type AddTrackResponse struct {
	TrackId int64 `json:"track_id"`
}

type RemoveTrackRequest struct {
	UserId  int64 `json:"user_id" db:"user_id" redis:"user_id"`
	TrackId int64 `json:"track_id" db:"track_id" redis:"track_id"`
}
