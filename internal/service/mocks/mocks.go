package mocks

import (
	"context"
	"music-lib/internal/model"
)

// MockArtistRepo для IArtistRepository
type MockArtistRepo struct {
	CreateFunc                func(ctx context.Context, entity *model.Artist) (*model.Artist, error)
	UpdateFunc                func(ctx context.Context, entity *model.Artist) (*model.Artist, error)
	DeleteFunc                func(ctx context.Context, id uint) error
	SearchFunc                func(ctx context.Context, query string, limit, offset int) ([]model.Artist, int64, error)
	GetByIDFunc               func(ctx context.Context, id uint) (*model.Artist, error)
	GetByUserIDFunc           func(ctx context.Context, userID uint) (*model.Artist, error)
	GetWithAlbumsFunc         func(ctx context.Context, id uint) (*model.Artist, error)
	IsExistsFunc              func(ctx context.Context, name string) bool
	GetArtistAlbumByUserIDFunc func(ctx context.Context, userID uint, albumID uint) (*model.Album, int, error)
}

func (m *MockArtistRepo) Create(ctx context.Context, entity *model.Artist) (*model.Artist, error) {
	return m.CreateFunc(ctx, entity)
}

func (m *MockArtistRepo) Update(ctx context.Context, entity *model.Artist) (*model.Artist, error) {
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

// MockAlbumRepo для IAlbumRepository
type MockAlbumRepo struct {
	CreateFunc        func(ctx context.Context, entity *model.Album) (*model.Album, error)
	UpdateFunc        func(ctx context.Context, entity *model.Album) (*model.Album, error)
	DeleteFunc        func(ctx context.Context, id uint) error
	SearchFunc        func(ctx context.Context, query string, limit, offset int) ([]model.Album, int64, error)
	GetByIDFunc       func(ctx context.Context, id uint) (*model.Album, error)
	GetWithSongsFunc func(ctx context.Context, id uint) (*model.Album, error)
}

func (m *MockAlbumRepo) Create(ctx context.Context, entity *model.Album) (*model.Album, error) {
	return m.CreateFunc(ctx, entity)
}

func (m *MockAlbumRepo) Update(ctx context.Context, entity *model.Album) (*model.Album, error) {
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

// MockSongRepo для ISongRepository
type MockSongRepo struct {
	CreateFunc         func(ctx context.Context, entity *model.Song) (*model.Song, error)
	UpdateFunc         func(ctx context.Context, entity *model.Song) (*model.Song, error)
	DeleteFunc         func(ctx context.Context, id uint) error
	SearchFunc         func(ctx context.Context, query string, limit, offset int) ([]model.Song, int64, error)
	ExistsInAlbumFunc  func(ctx context.Context, albumID uint, songName string) bool
	GetByIDFunc        func(ctx context.Context, id uint) (*model.Song, error)
	GetByArtistIDFunc  func(ctx context.Context, artistID uint, sort string, limit, offset int) ([]model.Song, int64, error)
	GetByAlbumIDFunc   func(ctx context.Context, albumID uint, sort string, limit, offset int) ([]model.Song, int64, error)
	GetFullInfoFunc    func(ctx context.Context, id uint) (*model.Song, *model.Artist, *model.Album, error)
}

func (m *MockSongRepo) Create(ctx context.Context, entity *model.Song) (*model.Song, error) {
	return m.CreateFunc(ctx, entity)
}

func (m *MockSongRepo) Update(ctx context.Context, entity *model.Song) (*model.Song, error) {
	return m.UpdateFunc(ctx, entity)
}

func (m *MockSongRepo) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockSongRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.Song, int64, error) {
	return m.SearchFunc(ctx, query, limit, offset)
}

func (m *MockSongRepo) ExistsInAlbum(ctx context.Context, albumID uint, songName string) bool {
	return m.ExistsInAlbumFunc(ctx, albumID, songName)
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

// MockLyricsRepo для ILyricsRepository
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

// MockUserRepo для IUserRepository
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

// MockProfileRepo для IProfileRepository
type MockProfileRepo struct {
	CreateFunc     func(ctx context.Context, entity *model.Profile) (*model.Profile, error)
	UpdateFunc     func(ctx context.Context, entity *model.Profile) (*model.Profile, error)
	DeleteFunc     func(ctx context.Context, id uint) error
	GetByUserIDFunc func(ctx context.Context, userID uint) (*model.Profile, error)
}

func (m *MockProfileRepo) Create(ctx context.Context, entity *model.Profile) (*model.Profile, error) {
	return m.CreateFunc(ctx, entity)
}

func (m *MockProfileRepo) Update(ctx context.Context, entity *model.Profile) (*model.Profile, error) {
	return m.UpdateFunc(ctx, entity)
}

func (m *MockProfileRepo) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockProfileRepo) GetByUserID(ctx context.Context, userID uint) (*model.Profile, error) {
	return m.GetByUserIDFunc(ctx, userID)
}

// MockGenreRepo для IGenreRepository
type MockGenreRepo struct {
	CreateFunc   func(ctx context.Context, entity *model.Genre) (*model.Genre, error)
	UpdateFunc   func(ctx context.Context, entity *model.Genre) (*model.Genre, error)
	DeleteFunc   func(ctx context.Context, id uint) error
	GetByIdFunc  func(ctx context.Context, id uint) (*model.Genre, error)
	GetByIdsFunc func(ctx context.Context, ids []uint) ([]model.Genre, error)
	IsExistsFunc func(ctx context.Context, name string) bool
}

func (m *MockGenreRepo) Create(ctx context.Context, entity *model.Genre) (*model.Genre, error) {
	return m.CreateFunc(ctx, entity)
}

func (m *MockGenreRepo) Update(ctx context.Context, entity *model.Genre) (*model.Genre, error) {
	return m.UpdateFunc(ctx, entity)
}

func (m *MockGenreRepo) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockGenreRepo) GetById(ctx context.Context, id uint) (*model.Genre, error) {
	return m.GetByIdFunc(ctx, id)
}

func (m *MockGenreRepo) GetByIds(ctx context.Context, ids []uint) ([]model.Genre, error) {
	return m.GetByIdsFunc(ctx, ids)
}

func (m *MockGenreRepo) IsExists(ctx context.Context, name string) bool {
	return m.IsExistsFunc(ctx, name)
}

// MockSongGenreRepo для ISongGenreRepository
type MockSongGenreRepo struct {
	CreateFunc func(ctx context.Context, entity *model.SongGenre) (*model.SongGenre, error)
	UpdateFunc func(ctx context.Context, entity *model.SongGenre) (*model.SongGenre, error)
	DeleteFunc func(ctx context.Context, id uint) error
}

func (m *MockSongGenreRepo) Create(ctx context.Context, entity *model.SongGenre) (*model.SongGenre, error) {
	return m.CreateFunc(ctx, entity)
}

func (m *MockSongGenreRepo) Update(ctx context.Context, entity *model.SongGenre) (*model.SongGenre, error) {
	return m.UpdateFunc(ctx, entity)
}

func (m *MockSongGenreRepo) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

// MockPermissionRepo для IPermissionRepository
type MockPermissionRepo struct {
	CreateFunc       func(ctx context.Context, entity *model.ResourcePermission) (*model.ResourcePermission, error)
	UpdateFunc       func(ctx context.Context, entity *model.ResourcePermission) (*model.ResourcePermission, error)
	DeleteFunc       func(ctx context.Context, id uint) error
	HasPermissionFunc func(userID, resourceID uint, resourceType model.Resource, permission model.Permission) bool
}

func (m *MockPermissionRepo) Create(ctx context.Context, entity *model.ResourcePermission) (*model.ResourcePermission, error) {
	return m.CreateFunc(ctx, entity)
}

func (m *MockPermissionRepo) Update(ctx context.Context, entity *model.ResourcePermission) (*model.ResourcePermission, error) {
	return m.UpdateFunc(ctx, entity)
}

func (m *MockPermissionRepo) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockPermissionRepo) HasPermission(userID, resourceID uint, resourceType model.Resource, permission model.Permission) bool {
	return m.HasPermissionFunc(userID, resourceID, resourceType, permission)
}