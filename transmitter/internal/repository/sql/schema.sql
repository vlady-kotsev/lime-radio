CREATE TABLE IF NOT EXISTS tracks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    artist TEXT,
    title TEXT
);

CREATE INDEX IF NOT EXISTS idx_tracks_path ON tracks(artist);