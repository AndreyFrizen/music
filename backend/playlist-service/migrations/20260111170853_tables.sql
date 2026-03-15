-- +goose Up
-- +goose StatementBegin
CREATE TABLE playlists (
	id INTEGER PRIMARY KEY,
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS playlists;
DROP TABLE IF EXISTS playlist_tracks;
-- +goose StatementEnd
