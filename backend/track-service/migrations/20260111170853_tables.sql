-- +goose Up
-- +goose StatementBegin
CREATE TABLE tracks (
	id INTEGER SERIAL PRIMARY KEY,
	title VARCHAR NOT NULL,
	duration INTEGER NOT NULL,
	audio_url VARCHAR NOT NULL,
	artist_id INTEGER NOT NULL,
	FOREIGN KEY (artist_id) REFERENCES artists(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tracks;
-- +goose StatementEnd
