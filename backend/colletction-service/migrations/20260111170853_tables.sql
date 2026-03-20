-- +goose Up
-- +goose StatementBegin
CREATE TABLE track_collection (
    track_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (track_id) REFERENCES tracks(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (track_id, user_id)
);

CREATE TABLE album_collection (
    album_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (album_id) REFERENCES albums(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE artist_collection (
    artist_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (artist_id) REFERENCES artists(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (artist_id, user_id)
);

create table playlist_collection (
    playlist_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (playlist_id) REFERENCES playlists(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (playlist_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS track_collection;
DROP TABLE IF EXISTS album_collection;
DROP TABLE IF EXISTS artist_collection;
DROP TABLE IF EXISTS playlist_collection;
-- +goose StatementEnd
