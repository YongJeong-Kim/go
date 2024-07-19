package service

type Service struct {
	Friend Friender
	User   UserManager
}

func NewService(friend Friender, user UserManager) *Service {
	return &Service{
		Friend: friend,
		User:   user,
	}
}
