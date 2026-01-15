-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    password TEXT NOT NULL ,
    email TEXT NOT NULL UNIQUE
);
CREATE TABLE artists (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);
CREATE TABLE albums (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	artist_id INTEGER NOT NULL,
	release_date TEXT NOT NULL,
	FOREIGN KEY (artist_id) REFERENCES artists(id)
);
CREATE TABLE tracks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	duration INTEGER NOT NULL,
	audio_url TEXT NOT NULL,
	artist_id INTEGER NOT NULL,
	FOREIGN KEY (artist_id) REFERENCES artists(id)
);
CREATE TABLE playlists (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE playlist_tracks (
	playlist_id INTEGER NOT NULL,
	track_id INTEGER NOT NULL,
	position INTEGER NOT NULL,
	FOREIGN KEY (playlist_id) REFERENCES playlists(id),
	FOREIGN KEY (track_id) REFERENCES tracks(id),
	PRIMARY KEY (playlist_id, track_id)
);
CREATE TABLE user_albums (
	user_id INTEGER NOT NULL,
	album_id INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (album_id) REFERENCES albums(id),
	PRIMARY KEY (user_id, album_id)
);
CREATE TABLE album_tracks (
	album_id INTEGER NOT NULL,
	track_id INTEGER NOT NULL,
	FOREIGN KEY (album_id) REFERENCES albums(id),
	FOREIGN KEY (track_id) REFERENCES tracks(id),
	PRIMARY KEY (album_id, track_id)
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
DROP TABLE IF EXISTS user_albums;
DROP TABLE IF EXISTS album_tracks;
-- +goose StatementEnd
