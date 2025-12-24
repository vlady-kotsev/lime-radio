INSERT INTO songs (
    id,
    artist,
    title,
    path
) VALUES (
    ?, ?, ?, ?
)
ON CONFLICT(artist, title) DO NOTHING;