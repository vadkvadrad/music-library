package mocks

import (
	"context"
	"music-lib/internal/model"
)

// Мок IArtistRepository
type MockArtistRepo struct {
	// Repository
	CreateFunc func(ctx context.Context, entity *model.Artist) error
	UpdateFunc func(ctx context.Context, entity *model.Artist) error
	DeleteFunc func(ctx context.Context, id uint) error

	// Searchable
	SearchFunc func(ctx context.Context, query string, limit, offset int) ([]model.Artist, int64, error)

	// IArtistRepository
	GetByIDFunc                func(ctx context.Context, id uint) (*model.Artist, error)
	GetByUserIDFunc            func(ctx context.Context, userID uint) (*model.Artist, error)
	GetWithAlbumsFunc          func(ctx context.Context, id uint) (*model.Artist, error)
	IsExistsFunc               func(ctx context.Context, name string) bool
	GetArtistAlbumByUserIDFunc func(ctx context.Context, userID uint, albumID uint) (*model.Album,int, error)
}

func (m *MockArtistRepo) Create(ctx context.Context, entity *model.Artist) error {
	return m.CreateFunc(ctx, entity)
}

func (m *MockArtistRepo) Update(ctx context.Context, entity *model.Artist) error {
	return m.UpdateFunc(ctx, entity)
}

func (m *MockArtistRepo) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockArtistRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.Artist, int64, error) {
	return m.SearchFunc(ctx, query, limit, offset)
}

func (m *MockArtistRepo) GetByID(ctx context.Context, id uint) (*model.Artist, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *MockArtistRepo) GetByUserID(ctx context.Context, userID uint) (*model.Artist, error) {
	return m.GetByUserIDFunc(ctx, userID)
}

func (m *MockArtistRepo) GetWithAlbums(ctx context.Context, id uint) (*model.Artist, error) {
	return m.GetWithAlbumsFunc(ctx, id)
}

func (m *MockArtistRepo) IsExists(ctx context.Context, name string) bool {
	return m.IsExistsFunc(ctx, name)
}

func (m *MockArtistRepo) GetArtistAlbumByUserID(ctx context.Context, userID uint, albumID uint) (*model.Album, int, error) {
	return m.GetArtistAlbumByUserIDFunc(ctx, userID, albumID)
}

// Мок IAlbumRepository
type MockAlbumRepo struct {
	// Repository
	CreateFunc func(ctx context.Context, entity *model.Album) error
	UpdateFunc func(ctx context.Context, entity *model.Album) error
	DeleteFunc func(ctx context.Context, id uint) error

	// Searchable
	SearchFunc func(ctx context.Context, query string, limit, offset int) ([]model.Album, int64, error)

	// IAlbumRepository
	GetByIDFunc      func(ctx context.Context, id uint) (*model.Album, error)
	GetWithSongsFunc func(ctx context.Context, id uint) (*model.Album, error)
}

func (m *MockAlbumRepo) Create(ctx context.Context, entity *model.Album) error {
	return m.CreateFunc(ctx, entity)
}

func (m *MockAlbumRepo) Update(ctx context.Context, entity *model.Album) error {
	return m.UpdateFunc(ctx, entity)
}

func (m *MockAlbumRepo) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockAlbumRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.Album, int64, error) {
	return m.SearchFunc(ctx, query, limit, offset)
}

func (m *MockAlbumRepo) GetByID(ctx context.Context, id uint) (*model.Album, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *MockAlbumRepo) GetWithSongs(ctx context.Context, id uint) (*model.Album, error) {
	return m.GetWithSongsFunc(ctx, id)
}

// Мок ISongRepository
type MockSongRepo struct {
	// Repository
	CreateFunc func(ctx context.Context, entity *model.Song) error
	UpdateFunc func(ctx context.Context, entity *model.Song) error
	DeleteFunc func(ctx context.Context, id uint) error

	// Searchable
	SearchFunc func(ctx context.Context, query string, limit, offset int) ([]model.Song, int64, error)

	// ISongRepository
	ExistsInAlbumFunc func(ctx context.Context, albumID uint, songName string)
	GetByIDFunc       func(ctx context.Context, id uint) (*model.Song, error)
	GetByArtistIDFunc func(ctx context.Context, artistID uint, sort string, limit, offset int) ([]model.Song, int64, error)
	GetByAlbumIDFunc  func(ctx context.Context, albumID uint, sort string, limit, offset int) ([]model.Song, int64, error)
	GetFullInfoFunc   func(ctx context.Context, id uint) (*model.Song, *model.Artist, *model.Album, error)
}

func (m *MockSongRepo) Create(ctx context.Context, entity *model.Song) error {
	return m.CreateFunc(ctx, entity)
}

func (m *MockSongRepo) Update(ctx context.Context, entity *model.Song) error {
	return m.UpdateFunc(ctx, entity)
}

func (m *MockSongRepo) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockSongRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.Song, int64, error) {
	return m.SearchFunc(ctx, query, limit, offset)
}

func (m *MockSongRepo) ExistsInAlbum(ctx context.Context, albumID uint, name string) bool {
	return m.ExistsInAlbum(ctx, albumID, name)
}

func (m *MockSongRepo) GetByID(ctx context.Context, id uint) (*model.Song, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *MockSongRepo) GetByArtistID(ctx context.Context, artistID uint, sort string, limit, offset int) ([]model.Song, int64, error) {
	return m.GetByArtistIDFunc(ctx, artistID, sort, limit, offset)
}

func (m *MockSongRepo) GetByAlbumID(ctx context.Context, albumID uint, sort string, limit, offset int) ([]model.Song, int64, error) {
	return m.GetByAlbumIDFunc(ctx, albumID, sort, limit, offset)
}

func (m *MockSongRepo) GetFullInfo(ctx context.Context, id uint) (*model.Song, *model.Artist, *model.Album, error) {
	return m.GetFullInfoFunc(ctx, id)
}

// Мок ILyricsRepository
type MockLyricsRepo struct {
	GetBySongIDFunc    func(ctx context.Context, songID uint) (*model.Lyrics, error)
	UpsertFunc         func(ctx context.Context, lyrics *model.Lyrics) error
	DeleteBySongIDFunc func(ctx context.Context, songID uint) error
}

func (m *MockLyricsRepo) GetBySongID(ctx context.Context, songID uint) (*model.Lyrics, error) {
	return m.GetBySongIDFunc(ctx, songID)
}

func (m *MockLyricsRepo) Upsert(ctx context.Context, lyrics *model.Lyrics) error {
	return m.UpsertFunc(ctx, lyrics)
}

func (m *MockLyricsRepo) DeleteBySongID(ctx context.Context, songID uint) error {
	return m.DeleteBySongIDFunc(ctx, songID)
}

// Мок IUserRepository
type MockUserRepo struct {
	CreateFunc    func(user *model.User) (*model.User, error)
	UpdateFunc    func(user *model.User) (*model.User, error)
	FindByKeyFunc func(key, data string) (*model.User, error)
}

func (m *MockUserRepo) Create(user *model.User) (*model.User, error) {
	return m.CreateFunc(user)
}

func (m *MockUserRepo) Update(user *model.User) (*model.User, error) {
	return m.UpdateFunc(user)
}

func (m *MockUserRepo) FindByKey(key, data string) (*model.User, error) {
	return m.FindByKeyFunc(key, data)
}

// Мок IProfileRepository
type MockProfileRepo struct {
	// Repository
	CreateFunc func(ctx context.Context, entity *model.Profile) error
	UpdateFunc func(ctx context.Context, entity *model.Profile) error
	DeleteFunc func(ctx context.Context, id uint) error

	// IProfileRepository
	GetByUserIDFunc func(ctx context.Context, userID uint) (*model.Profile, error)
}

func (m *MockProfileRepo) Create(ctx context.Context, entity *model.Profile) error {
	return m.CreateFunc(ctx, entity)
}

func (m *MockProfileRepo) Update(ctx context.Context, entity *model.Profile) error {
	return m.UpdateFunc(ctx, entity)
}

func (m *MockProfileRepo) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockProfileRepo) GetByUserID(ctx context.Context, userID uint) (*model.Profile, error) {
	return m.GetByUserIDFunc(ctx, userID)
}
