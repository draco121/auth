package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auth/custom_models"
	"auth/database"
	"auth/graph/generated"
	"auth/graph/model"
	"auth/utils"
	"context"
	"fmt"
	"time"
)

type ContextItems struct {
	Sessionid *string
	Database  *database.DB
}

func (r *mutationResolver) Login(ctx context.Context, input *model.LoginInput) (string, error) {
	items, _ := ctx.Value("context_items").(ContextItems)
	db := items.Database
	if input.Identifier.Phonenumber == nil && input.Identifier.Username == nil {
		return "", fmt.Errorf("please provide username or phonenumber")
	} else {
		var user *custom_models.User
		var err error
		if input.Identifier.Phonenumber != nil {
			user, err = db.FindOneByPhonenumber(*input.Identifier.Phonenumber)
			if err != nil {
				return "", err
			}
		} else {
			user, err = db.FindOneByUsername(*input.Identifier.Username)
			if err != nil {
				return "", err
			}
		}
		util := utils.Utils{}
		if err := util.ComparePassword(*user.Password, input.Password); err == nil {
			token, err := util.CreateJwt(*user.Username)
			if err != nil {
				return "", err
			} else {
				tokenmodel := custom_models.Token{Username: *user.Username, Token: token, Timestamp: time.Now().Unix()}
				db.InsertToken(&tokenmodel)
				return token, nil
			}
		} else {
			return "", err
		}
	}
}

func (r *mutationResolver) Signout(ctx context.Context) (bool, error) {
	items, _ := ctx.Value("context_items").(ContextItems)
	util := utils.Utils{}
	db := items.Database
	defer db.Disconnect()
	sessionid := items.Sessionid
	_, err := util.ValidateJwt(*sessionid)
	if err != nil {
		return false, err
	} else {
		_, err := db.IsTokenExists(*sessionid)
		if err != nil {
			return false, err
		} else {
			_, err := db.FindOneAndDeleteToken(*sessionid)
			return true, err
		}
	}
}

func (r *mutationResolver) Createaccount(ctx context.Context, input model.UserInput) (bool, error) {
	items, _ := ctx.Value("context_items").(ContextItems)
	db := items.Database
	defer db.Disconnect()
	_, err := db.FindOneByPhonenumber(input.Phonenumber)
	if err != nil {
		_, err := db.FindOneByUsername(input.Username)
		if err != nil {
			util := utils.Utils{}
			var e error
			input.Password, e = util.GenerateHashedPassword(input.Password)
			if e != nil {
				return false, e
			}
			db.Save(&input)
			return true, nil
		} else {
			return false, fmt.Errorf("username is already in use")
		}
	} else {
		return false, fmt.Errorf("account with this phone number already exists")
	}
}

func (r *queryResolver) Getuser(ctx context.Context) (*model.User, error) {
	context, _ := ctx.Value("context_items").(ContextItems)
	util := utils.Utils{}
	db := context.Database
	defer db.Disconnect()
	username, err := util.ValidateJwt(*context.Sessionid)
	if err != nil {
		return nil, fmt.Errorf("unauthorized access")
	} else {
		exists, err := db.IsTokenExists(*context.Sessionid)
		if exists {
			user, err := db.FindOneByUsername(username)
			if err != nil {
				return nil, fmt.Errorf("user not fount")
			} else {
				result := model.User{ID: user.ID, Username: user.Username, Phonenumber: user.Phonenumber}
				return &result, nil
			}
		} else {
			return nil, err
		}
	}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
