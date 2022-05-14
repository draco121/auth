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

func (r *mutationResolver) Login(ctx context.Context, input *model.LoginInput) (string, error) {
	if input.Identifier.Phonenumber == nil && input.Identifier.Username == nil {
		return "", fmt.Errorf("please provide username or phonenumber")
	} else {
		if input.Identifier.Phonenumber != nil {
			db := database.Connect()
			defer db.Disconnect()
			user, err := db.FindOneByPhonenumber(*input.Identifier.Phonenumber)
			if err != nil {
				return "", err
			} else {
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
		} else {
			db := database.Connect()
			defer db.Disconnect()
			user, err := db.FindOneByUsername(*input.Identifier.Username)
			if err != nil {
				return "", err
			} else {
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
	}
}

func (r *mutationResolver) Signout(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Createaccount(ctx context.Context, input model.UserInput) (bool, error) {
	db := database.Connect()
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
	sessionid, _ := ctx.Value("sessionid").(string)
	util := utils.Utils{}
	db := database.Connect()
	defer db.Disconnect()
	username, err := util.ValidateJwt(sessionid)
	if err != nil {
		return nil, fmt.Errorf("unauthorized access")
	} else {
		user, err := db.FindOneByUsername(username)
		if err != nil {
			return nil, fmt.Errorf("user not fount")
		} else {
			result := model.User{ID: user.ID, Username: user.Username, Phonenumber: user.Phonenumber}
			return &result, nil
		}
	}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
