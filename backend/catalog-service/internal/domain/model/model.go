package model

// Artist represents an artist in the system
type Artist struct {
	ID   int64  `json:"id" db:"id" redis:"id"`
	Name string `json:"name" db:"name" redis:"name"`
}

// Album represents an album in the system
type Album struct {
	ID          int64  `json:"id" db:"id" redis:"id"`
	Title       string `json:"title" db:"title" redis:"title"`
	ArtistID    int64  `json:"artist_id" db:"artist_id" redis:"artist_id"`
	ReleaseDate string `json:"release_date" db:"release_date" redis:"release_date"`
}

func (a *Album) GetID() int64 {
	return a.ID
}

func (a *Artist) GetID() int64 {
	return a.ID
}

type GetAlbumRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type GetAlbumResponse struct {
	Album *Album `json:"album"`
}

type CreateAlbumRequest struct {
	Title       string `json:"title" validate:"required"`
	ArtistID    int64  `json:"artist_id" validate:"required"`
	ReleaseDate string `json:"release_date" validate:"required"`
}

type CreateAlbumResponse struct {
	ID int64 `json:"id"`
}

type DeleteAlbumRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type DeleteAlbumResponse struct {
	ID int64 `json:"id"`
}

type GetArtistRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type GetArtistResponse struct {
	Artist *Artist `json:"artist"`
}

type CreateArtistRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateArtistResponse struct {
	ID int64 `json:"id"`
}

type DeleteArtistRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type DeleteArtistResponse struct {
	ID int64 `json:"id"`
}
