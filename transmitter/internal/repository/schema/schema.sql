CREATE TABLE IF NOT EXISTS songs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    artist TEXT,
    title TEXT,
    path TEXT,

    UNIQUE (artist, title)
);

CREATE INDEX IF NOT EXISTS idx_songs_artist ON songs(artist);