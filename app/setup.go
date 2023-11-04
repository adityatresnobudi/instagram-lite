package app

import (
	"errors"
	"fmt"
	"instagram-lite/entity"
)

const (
	Upload string = "upload"
	Like   string = "like"
)

var (
	ErrNoPhoto         = errors.New("you don't have a photo")
	ErrNotFollowed     = errors.New("not following the account")
	ErrSameAccount     = errors.New("a user cannot follow themselves")
	ErrAlreadyFollowed = errors.New("you already followed the user")
	ErrUploadTwice     = errors.New("you cannot upload more than once")
	ErrLikedTwice      = errors.New("you already liked the photo")
)

type Account struct {
	username      *entity.User
	photo         *entity.Photo
	followingList []*Account
	followerList  []*Account
	activity      []*Activity
}

func NewAccount(username *entity.User) *Account {
	return &Account{
		username: username,
		photo: &entity.Photo{
			IsExist: false,
			Like:    make([]*entity.User, 0),
		},
		followingList: make([]*Account, 0),
		followerList:  make([]*Account, 0),
		activity:      make([]*Activity, 0),
	}
}

func (a *Account) Follow(acc *Account) ([]*Account, []*Account, error) {
	if a.IsSameAccount(acc) {
		return nil, nil, ErrSameAccount
	}

	if a.HasFollow(acc) {
		return nil, nil, ErrAlreadyFollowed
	}

	a.followingList = append(a.followingList, acc)
	acc.followerList = append(acc.followerList, a)

	return a.followingList, acc.followerList, nil
}

func (a *Account) Post() ([]*Activity, error) {
	if a.photo.IsExist {
		return nil, ErrUploadTwice
	}

	a.photo.IsExist = true
	action := NewActivity(a, Upload, a)
	a.activity = append(a.activity, action)
	a.notifyFollowerUpload(action)

	return a.activity, nil
}

func (a *Account) Like(acc *Account) ([]*Activity, []*Activity, error) {
	switch a.IsSameAccount(acc) {
	case true:
		if !a.HasUploadPhoto() {
			return nil, nil, ErrNoPhoto
		}

		action := NewActivity(a, Like, a)

		if a.HasLikedPhoto(acc, action) {
			return nil, nil, ErrLikedTwice
		}

		a.activity = append(a.activity, action)
		a.notifyFollowerLike(action)
	case false:
		if !a.HasFollow(acc) {
			return nil, nil, fmt.Errorf("unable to like %s's photo", acc.username.Name)
		}

		if !acc.HasUploadPhoto() {
			return nil, nil, fmt.Errorf("%s doesn't have a photo", acc.username.Name)
		}

		action := NewActivity(a, Like, acc)

		if a.HasLikedPhoto(acc, action) {
			return nil, nil, ErrLikedTwice
		}

		a.activity = append(a.activity, action)
		acc.activity = append(acc.activity, action)
		a.notifyFollowerLike(action)
	}

	acc.photo.Like = append(acc.photo.Like, a.username)
	return a.activity, acc.activity, nil
}

func (a *Account) IsSameAccount(acc *Account) bool {
	return a.GetUsername() == acc.GetUsername()
}

func (a *Account) HasFollow(acc *Account) bool {
	for _, account := range acc.followerList {
		if a == account {
			return true
		}
	}
	return false
}

func (a *Account) HasUploadPhoto() bool {
	return a.photo.IsExist
}

func (a *Account) HasLikedPhoto(acc *Account, action *Activity) bool {
	for _, act := range a.activity {
		if act.IsSameActivity(action) {
			return true
		}
	}
	return false
}

func (a *Account) GetActivity() []*Activity {
	return a.activity
}

func (a *Account) GetUsername() string {
	return a.username.Name
}

func (a *Account) GetPhoto() *entity.Photo {
	return a.photo
}

func (a *Account) notifyFollowerUpload(action *Activity) {
	for _, acc := range a.followerList {
		action := NewActivity(a, Upload, acc)
		acc.activity = append(acc.activity, action)
	}
}

func (a *Account) notifyFollowerLike(action *Activity) {
	for _, acc := range a.followerList {
		if len(acc.activity) == 0 {
			if !acc.IsSameAccount(a) {
				acc.activity = append(acc.activity, action)
				continue
			}
		}

		if acc.activity[len(acc.activity)-1].IsSameActivity(action) {
			continue
		}

		acc.activity = append(acc.activity, action)
	}
}

func (a *Account) HasMorePhotoLike(acc *Account) bool {
	return len(a.photo.Like) > len(acc.photo.Like)
}
