package file_service

import (
	"fmt"

	"github.com/dev-tams/file-upload/internal/models"
)

var planStorage = map[string]int{
	"free":       100 * 1024 * 1024,   // 100 MB
	"pro":        1000 * 1024 * 1024,  // 1 GB
	"enterprise": 10 * 1024 * 1024 * 1024, // 10 GB
}

func (u *Service) AssignPlan(userID, plan string) (*models.User, error) {

	limit, ok := planStorage[plan]
	if !ok {
		return nil, fmt.Errorf("invalid plan: %s", plan)
	}
	if err:= u.repo.AssignPlan(userID, plan); err != nil {
		return nil, err
	}
	if err := u.repo.UpdateUserStorage(userID, int64(limit)); err != nil {
		return nil, err
	}

	user, err := u.repo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	 return user, nil
}
