INSERT INTO songs (
    artist,
    title,
    path
) VALUES (
    ?, ?, ?
)
ON CONFLICT(artist, title) DO NOTHING;