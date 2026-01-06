-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL
    );
CREATE TABLE artists (
	id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);
CREATE TABLE albums (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	artist_id TEXT NOT NULL,
	release_date TEXT NOT NULL,
	FOREIGN KEY (artist_id) REFERENCES artists(id)
);
CREATE TABLE tracks (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	duration INTEGER NOT NULL,
	audio_url TEXT NOT NULL
);
CREATE TABLE playlists (
	playlist_id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	title TEXT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE playlist_tracks (
	playlist_id TEXT NOT NULL,
	track_id TEXT NOT NULL,
	position INTEGER NOT NULL,
	FOREIGN KEY (playlist_id) REFERENCES playlists(playlist_id),
	FOREIGN KEY (track_id) REFERENCES tracks(id),
	PRIMARY KEY (playlist_id, track_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS artists;
DROP TABLE IF EXISTS albums;
DROP TABLE IF EXISTS tracks;
DROP TABLE IF EXISTS playlists;
DROP TABLE IF EXISTS playlist_tracks;
-- +goose StatementEnd
