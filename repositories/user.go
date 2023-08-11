package repositories

import "hexagonal-gotest/models"

// PORT user repository
type UserRepository interface {
	Gets(filter models.RepoGetUserModel) (result []models.RepoUserModel, err error)

	Create(payload models.RepoCreateUserModel) (err error)

	Update(userId string, payload models.RepoUpdateUserModel) (err error)

	Delete(userId string) (err error)
}
