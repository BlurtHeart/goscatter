package test

type Company struct {
	User Talker
}

func NewCompany(t Talker) *Company {
	return &Company{
		User: t,
	}
}

func (c *Company) Meeting(guestName string) string {
	return c.User.Talk(guestName)
}
