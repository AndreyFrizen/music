-- +goose Up
-- +goose StatementBegin
CREATE TABLE artists (
	id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE artists_tracks (
    artist_id INTEGER NOT NULL,
    track_id INTEGER NOT NULL,
	FOREIGN KEY (artist_id) REFERENCES artists(id),
	PRIMARY KEY (track_id, artist_id)
);

CREATE TABLE albums (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	artist_id INTEGER NOT NULL,
	release_date TEXT NOT NULL,
	FOREIGN KEY (artist_id) REFERENCES artists(id)
);

CREATE TABLE album_tracks (
	album_id INTEGER NOT NULL,
	track_id INTEGER NOT NULL,
	FOREIGN KEY (album_id) REFERENCES albums(id),
	PRIMARY KEY (album_id, track_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS artists;
DROP TABLE IF EXISTS albums;
DROP TABLE IF EXISTS album_tracks;
DROP TABLE IF EXISTS artists_tracks;
-- +goose StatementEnd
