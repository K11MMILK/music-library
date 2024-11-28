SET datestyle = 'German, YMD';

CREATE TABLE groupss
(
    id serial PRIMARY KEY,
    groupName VARCHAR(255) NOT NULL
);
CREATE TABLE songs
(
    id serial PRIMARY KEY,
    songName VARCHAR(255) NOT NULL,
    groupId INT NOT NULL,
    FOREIGN KEY (groupId) REFERENCES groupss(id) ON DELETE CASCADE
);

CREATE TABLE songDetails
(
    id serial PRIMARY KEY,
    releaseDate DATE,
    text TEXT DEFAULT 'N/A',
    link VARCHAR(255) DEFAULT 'N/A',
    songId INT NOT NULL UNIQUE,
    FOREIGN KEY (songId) REFERENCES songs(id) ON DELETE CASCADE
);
CREATE INDEX idx_group_name ON groupss(groupName);

CREATE INDEX idx_song_name ON songs(songName);

CREATE INDEX idx_song_group_id ON songs(groupId);

CREATE INDEX idx_song_details_release_date ON songDetails(releaseDate);

CREATE INDEX idx_song_details_link ON songDetails(link);

CREATE INDEX idx_song_details_song_id ON songDetails(songId);
