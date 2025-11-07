package repositories

import (
	"github.com/dev-tams/file-upload/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)



// WithUserLock safely executes a DB operation with a locked user row.
func (r *DbRepository) WithUserLock(userID string, fn func(tx *gorm.DB, user *models.User) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lock the user row
	var user models.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&user, "id = ?", userID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Run the provided function inside the locked context
	if err := fn(tx, &user); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
