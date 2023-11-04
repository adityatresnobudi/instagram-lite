package app

type Activity struct {
	accDo  *Account
	action string
	accTo  *Account
}

func NewActivity(acc1 *Account, action string, acc2 *Account) *Activity {
	return &Activity{
		accDo:  acc1,
		action: action,
		accTo:  acc2,
	}
}

func (ac *Activity) IsSameActivity(act *Activity) bool {
	return ac.accDo.IsSameAccount(act.accDo) && ac.action == act.action && ac.accTo.IsSameAccount(act.accTo)
}

func (ac *Activity) GetAccDo() *Account {
	return ac.accDo
}

func (ac *Activity) GetAction() string {
	return ac.action
}

func (ac *Activity) GetAccTo() *Account {
	return ac.accTo
}