package app

import (
	"errors"
	"sort"
)

type AccRegistry struct {
	AccountList []*Account
}

var (
	ErrUserExist = errors.New("username already exist")
)

func NewAccRegistry() *AccRegistry {
	return &AccRegistry{
		AccountList: make([]*Account, 0),
	}
}

func (ar *AccRegistry) Record(acc *Account) ([]*Account, error) {
	if _, res := ar.IsAccountExist(acc); res {
		return nil, ErrUserExist
	}
	ar.AccountList = append(ar.AccountList, acc)
	return ar.AccountList, nil
}

func (ar *AccRegistry) IsAccountExist(acc *Account) (int, bool) {
	for idx, account := range ar.AccountList {
		if acc.IsSameAccount(account) {
			return idx, true
		}
	}
	return -1, false
}

func (ar *AccRegistry) GetLeaderboard() []*Account {
	accountList := ar.AccountList
	sort.Slice(accountList, func(i int, j int) bool {
		return accountList[i].HasMorePhotoLike(accountList[j])
	})
	return accountList
}
