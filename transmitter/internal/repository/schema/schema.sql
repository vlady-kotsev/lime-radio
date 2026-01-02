CREATE TABLE IF NOT EXISTS songs (
    id TEXT PRIMARY KEY,
    artist TEXT,
    title TEXT,
    path TEXT,

    UNIQUE (artist, title)
);

CREATE INDEX IF NOT EXISTS idx_songs_artist ON songs(artist);
CREATE INDEX IF NOT EXISTS idx_songs_title ON songs(title);
