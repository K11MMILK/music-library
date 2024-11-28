package musiclibrary

type Group struct {
	Id        int    `json:"id" db:"id"`
	GroupName string `json:"groupName" db:"groupname" binding:"required"`
}

type Song struct {
	Id       int    `json:"id" db:"id"`
	SongName string `json:"songName" db:"songname" binding:"required"`
	GroupId  int    `json:"groupId" db:"groupid"`
}

type UpdateGroupInput struct {
	GroupName *string `json:"groupName" example:"Metallica"`
}

type CreateGroupInput struct {
	GroupName string `json:"groupName" binding:"required" example:"Metallica"`
}

type UpdateSongInput struct {
	SongName *string `json:"songName" example:"Enter Sandman"`
	GroupId  *int    `json:"groupId"`
}

type CreateSongInput struct {
	SongName string `json:"songName" binding:"required" example:"Enter Sandman"`
	GroupId  int    `json:"groupId" binding:"required"`
}

type UpdateSongDetailsInput struct {
	ReleaseDate string `json:"releaseDate" example:"2024-01-01"`
	Text        string `json:"text" example:"Song lyrics here"`
	Link        string `json:"link" example:"https://example.com/song"`
}
type SongDetails struct {
	Id          int    `json:"id" db:"id"`
	SongId      int    `json:"songId" db:"songid"`
	ReleaseDate string `json:"releaseDate" db:"releasedate"`
	Text        string `json:"text" db:"text"`
	Link        string `json:"link" db:"link"`
}
type SongDetailsDL struct {
	Id          int    `json:"id" db:"id"`
	SongId      int    `json:"songId" db:"songid"`
	ReleaseDate string `json:"releaseDate" db:"releasedate"`
	Link        string `json:"link" db:"link"`
}
type SongDetailsT struct {
	Id     int    `json:"id" db:"id"`
	SongId int    `json:"songId" db:"songid"`
	Text   string `json:"text" db:"text"`
}
