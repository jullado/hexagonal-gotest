package repositories_test

import (
	"fmt"
	"hexagonal-gotest/models"
	"hexagonal-gotest/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGetsUser(t *testing.T) {
	type args struct {
		filter models.RepoGetUserModel
	}
	tests := []struct {
		name       string
		args       args
		wantResult []models.RepoUserModel
		wantErr    bool
	}{
		{
			name: "success1",
			args: args{
				filter: models.RepoGetUserModel{
					Username: "admin",
				},
			},
			wantResult: []models.RepoUserModel{
				{Username: "admin"},
			},
			wantErr: false,
		},
		{
			name: "error1",
			args: args{
				filter: models.RepoGetUserModel{
					Username: "admin",
				},
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "error2",
			args: args{
				filter: models.RepoGetUserModel{
					Username: "admin",
				},
			},
			wantResult: nil,
			wantErr:    true,
		},
	}

	collection := "users"
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mt.Run(tt.name, func(mt *mtest.T) {
				// ------------------- Arrange (เตรียมของ) --------------------

				// mock Find
				if tt.wantErr {
					mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
						Index:   1,
						Code:    123,
						Message: "some error",
					}))
				} else {
					data, _ := bson.Marshal(tt.args.filter)
					doc := bson.D{}
					bson.Unmarshal(data, &doc)

					first := mtest.CreateCursorResponse(1, fmt.Sprintf("%v.%v", "DBtest", collection), mtest.FirstBatch, doc)
					killCursors := mtest.CreateCursorResponse(0, fmt.Sprintf("%v.%v", "DBtest", collection), mtest.NextBatch)
					mt.AddMockResponses(first, killCursors)
				}

				userRepo := repositories.NewUserRepository(mt.DB, collection)

				// -------------------- Act (กระทำ)--------------------
				gotResult, err := userRepo.Gets(models.RepoGetUserModel{UserId: tt.args.filter.UserId, Username: tt.args.filter.Username})

				// -------------------- Assert (ยืนยัน) --------------------
				assert.Equal(mt, tt.wantErr, err != nil)
				assert.Equal(mt, tt.wantResult, gotResult)
			})
		})
	}
}

func TestCreateUser(t *testing.T) {
	type args struct {
		payload models.RepoCreateUserModel
	}
	tests := []struct {
		name       string
		args       args
		wantResult bson.D
		wantErr    bool
	}{
		{
			name: "error1",
			args: args{
				payload: models.RepoCreateUserModel{
					UserId:   "123",
					Username: "admin",
					Password: "admin01",
				},
			},
			wantResult: mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   1,
				Code:    123,
				Message: "some error",
			}),
			wantErr: true,
		},
		{
			name: "success1",
			args: args{
				payload: models.RepoCreateUserModel{
					UserId:   "123",
					Username: "admin",
					Password: "admin01",
				},
			},
			wantResult: mtest.CreateSuccessResponse(),
			wantErr:    false,
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mt.Run(tt.name, func(mt *mtest.T) {
				// ------------------- Arrange (เตรียมของ) --------------------

				// mock InsertOne
				mt.AddMockResponses(tt.wantResult)

				userRepo := repositories.NewUserRepository(mt.DB, "users")

				// -------------------- Act (กระทำ)--------------------
				err := userRepo.Create(models.RepoCreateUserModel{UserId: tt.args.payload.UserId, Username: tt.args.payload.Username, Password: tt.args.payload.Password})

				// -------------------- Assert (ยืนยัน) --------------------
				assert.Equal(mt, tt.wantErr, err != nil)
			})
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type args struct {
		userId  string
		payload models.RepoUpdateUserModel
	}
	tests := []struct {
		name       string
		args       args
		wantResult bson.D
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "error1",
			args: args{
				userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191",
				payload: models.RepoUpdateUserModel{
					Username: "admin",
					Password: "admin01",
				},
			},
			wantResult: mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   1,
				Code:    123,
				Message: "some error",
			}),
			wantErr: true,
		},
		{
			name: "error2",
			args: args{
				userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191",
				payload: models.RepoUpdateUserModel{
					Username: "admin",
					Password: "admin01",
				},
			},
			wantResult: mtest.CreateSuccessResponse(bson.E{Key: "ok", Value: "0"}, bson.E{Key: "nModified", Value: 0}, bson.E{Key: "n", Value: 0}),
			wantErr:    true,
		},
		{
			name: "success1",
			args: args{
				userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191",
				payload: models.RepoUpdateUserModel{
					Username: "admin",
					Password: "admin01",
				},
			},
			wantResult: mtest.CreateSuccessResponse(bson.E{Key: "ok", Value: "1"}, bson.E{Key: "nModified", Value: 1}, bson.E{Key: "n", Value: 1}),
			wantErr:    false,
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt.Run(tt.name, func(mt *mtest.T) {
				// ------------------- Arrange (เตรียมของ) --------------------

				// mock UpdateOne
				mt.AddMockResponses(tt.wantResult)

				userRepo := repositories.NewUserRepository(mt.DB, "users")

				// -------------------- Act (กระทำ)--------------------
				err := userRepo.Update(tt.args.userId, models.RepoUpdateUserModel{Username: tt.args.payload.Username, Password: tt.args.payload.Password})

				// -------------------- Assert (ยืนยัน) --------------------
				assert.Equal(mt, tt.wantErr, err != nil)
			})
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type args struct {
		userId string
	}
	tests := []struct {
		name       string
		args       args
		wantResult bson.D
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "error1",
			args: args{
				userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191",
			},
			wantResult: mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   1,
				Code:    123,
				Message: "some error",
			}),
			wantErr: true,
		},
		{
			name: "error2",
			args: args{
				userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191",
			},
			wantResult: mtest.CreateSuccessResponse(bson.E{Key: "ok", Value: "0"}, bson.E{Key: "acknowledged", Value: true}, bson.E{Key: "n", Value: 0}),
			wantErr:    true,
		},
		{
			name: "success1",
			args: args{
				userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191",
			},
			wantResult: mtest.CreateSuccessResponse(bson.E{Key: "ok", Value: "1"}, bson.E{Key: "acknowledged", Value: true}, bson.E{Key: "n", Value: 1}),
			wantErr:    false,
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt.Run(tt.name, func(mt *mtest.T) {
				// ------------------- Arrange (เตรียมของ) --------------------

				// mock DeleteOne
				mt.AddMockResponses(tt.wantResult)

				userRepo := repositories.NewUserRepository(mt.DB, "users")

				// -------------------- Act (กระทำ)--------------------
				err := userRepo.Delete(tt.args.userId)

				// -------------------- Assert (ยืนยัน) --------------------
				assert.Equal(mt, tt.wantErr, err != nil)
			})
		})
	}
}
