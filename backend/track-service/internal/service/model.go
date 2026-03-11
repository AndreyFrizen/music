package services

type CreateTrackRequest struct {
	Title    string `json:"title" validate:"required"`
	Duration int64  `json:"duration" validate:"required"`
	AudioURL string `json:"audio_url"`
	ArtistID int64  `json:"artist_id"`
	AlbumID  int64  `json:"album_id"`
}

type CreateTrackResponse struct {
	ID int64
}

type GetTrackRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type GetTrackResponse struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Duration int64  `json:"duration"`
	AudioURL string `json:"audio_url"`
	ArtistID int64  `json:"artist_id"`
	AlbumID  int64  `json:"album_id"`
}

type UpdateTrackRequest struct {
	ID       int64  `json:"id" validate:"required"`
	Title    string `json:"title"`
	Duration int64  `json:"duration"`
	AudioURL string `json:"audio_url"`
	ArtistID int64  `json:"artist_id"`
	AlbumID  int64  `json:"album_id"`
}

type UpdateTrackResponse struct {
	ID int64
}

type DeleteTrackRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type DeleteTrackResponse struct {
	Success bool `json:"success"`
}
