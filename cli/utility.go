package cli

import (
	"errors"
	"fmt"
	"instagram-lite/app"
	"instagram-lite/entity"
	"strings"
)

const (
	keyFollow string = "follows"
	keyLike   string = "likes"
	keyUpload string = "uploaded"
	keyPhoto  string = "photo"
)

var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrInvalidKeyword = errors.New("invalid keyword")
	registry          = app.NewAccRegistry()
)

func HandleSetup(relation string) ([]*app.Account, []*app.Account, []*app.Account, error) {
	if isEmpty(relation) {
		return nil, nil, nil, ErrInvalidInput
	}

	relationList := strings.Split(relation, " ")
	subject := make([]*app.Account, 0)

	if len(relationList) != 3 {
		return nil, nil, nil, ErrInvalidInput
	}

	if relationList[1] != keyFollow {
		return nil, nil, nil, ErrInvalidKeyword
	}

	for idx, v := range relationList {
		if idx != 1 {
			a := app.NewAccount(&entity.User{Name: string(v)})
			i, res := registry.IsAccountExist(a)
			if !res {
				registry.Record(a)
				subject = append(subject, a)
				continue
			}
			subject = append(subject, registry.AccountList[i])
		}
	}

	following, follower, err := subject[0].Follow(subject[1])
	return registry.AccountList, following, follower, err
}

func HandleAction(action string) ([]*app.Activity, []*app.Activity, error) {
	if isEmpty(action) {
		return nil, nil, ErrInvalidInput
	}

	if strings.Contains(action, keyLike) && strings.Contains(action, keyUpload) {
		return nil, nil, ErrInvalidKeyword
	}

	if strings.Contains(action, keyLike) {
		return handleLike(action)
	}

	if strings.Contains(action, keyUpload) {
		return handlePost(action)
	}

	return nil, nil, ErrInvalidKeyword
}

func HandleDisplay(display string) (string, error) {
	result := ""
	if isEmpty(display) {
		return "", ErrInvalidInput
	}

	a := app.NewAccount(&entity.User{Name: display})
	i, res := registry.IsAccountExist(a)
	if !res {
		return "", fmt.Errorf("unknown user %s", display)
	}

	result += "\n"
	result += fmt.Sprintf("%s activities:\n", display)
	activity := registry.AccountList[i].GetActivity()
	for _, act := range activity {
		if act.GetAction() == app.Upload {
			if act.GetAccDo() == act.GetAccTo() {
				result += "You uploaded photo\n"
				continue
			}
			result += fmt.Sprintf("%s uploaded photo\n", act.GetAccDo().GetUsername())
			continue
		}

		if act.GetAction() == app.Like {
			if act.GetAccDo().GetUsername() == display && act.GetAccTo().GetUsername() == display {
				result += "You liked your photo\n"
				continue
			}

			if act.GetAccDo().GetUsername() == display {
				result += fmt.Sprintf("You liked %s's photo\n", act.GetAccTo().GetUsername())
				continue
			}

			if act.GetAccTo().GetUsername() == display {
				result += fmt.Sprintf("%s liked your photo\n", act.GetAccDo().GetUsername())
				continue
			}

			result += fmt.Sprintf("%s liked %s's photo\n", act.GetAccDo().GetUsername(), act.GetAccTo().GetUsername())
			continue
		}
	}
	return result, nil
}

func HandleTrending() string {
	result := "Trending photos:\n"
	leaderboard := registry.GetLeaderboard()
	for idx, v := range leaderboard {
		if idx > 2 {
			break
		}

		if like := len(v.GetPhoto().Like); like != 0 {
			result += fmt.Sprintf("%d. %s photo got %d likes\n", idx+1, v.GetUsername(), like)
		}
	}
	return result
}

func handleLike(action string) ([]*app.Activity, []*app.Activity, error) {
	subject := make([]*app.Account, 0)
	arrAction := strings.Split(action, " ")
	if arrAction[1] != keyLike || arrAction[3] != keyPhoto {
		return nil, nil, ErrInvalidKeyword
	}

	for idx, v := range arrAction {
		if idx%2 == 0 {
			a := app.NewAccount(&entity.User{Name: string(v)})
			i, res := registry.IsAccountExist(a)
			if !res {
				return nil, nil, fmt.Errorf("unknown user %s", string(v))
			}
			subject = append(subject, registry.AccountList[i])
		}
	}
	return subject[0].Like(subject[1])
}

func handlePost(action string) ([]*app.Activity, []*app.Activity, error) {
	subject := make([]*app.Account, 0)
	arrAction := strings.Split(action, " ")
	if arrAction[1] != keyUpload || arrAction[2] != keyPhoto {
		return nil, nil, ErrInvalidKeyword
	}

	for idx, v := range arrAction {
		if idx == 0 {
			a := app.NewAccount(&entity.User{Name: string(v)})
			i, res := registry.IsAccountExist(a)
			if !res {
				return nil, nil, fmt.Errorf("unknown user %s", string(v))
			}
			subject = append(subject, registry.AccountList[i])
		}
	}

	act, err := subject[0].Post()
	return nil, act, err
}

func isEmpty(arg string) bool {
	return arg == ""
}
