CREATE TABLE songs
(
    id serial PRIMARY KEY,
    "group" VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL
);

CREATE TABLE songDetails
(
    id serial PRIMARY KEY,
    releaseDate VARCHAR(255) DEFAULT 'N/A',
    text TEXT DEFAULT 'N/A',
    link VARCHAR(255) DEFAULT 'N/A',
    songId INT NOT NULL UNIQUE,
    FOREIGN KEY (songId) REFERENCES songs(id) ON DELETE CASCADE
);