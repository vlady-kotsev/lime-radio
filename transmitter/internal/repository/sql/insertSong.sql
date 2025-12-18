INSERT INTO tracks (
    artist,
    title
) VALUES (
    ?, ?
)
ON CONFLICT(id) DO UPDATE SET
    artist = excluded.artist,
    title = excluded.title;