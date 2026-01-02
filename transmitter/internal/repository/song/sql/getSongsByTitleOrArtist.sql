SELECT id, artist, title, path 
FROM songs
WHERE artist LIKE ? OR title LIKE ?;