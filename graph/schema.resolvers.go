package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"authentication/custom_models"
	"authentication/database"
	"authentication/graph/generated"
	"authentication/graph/model"
	"authentication/grpcclient"
	"authentication/utils"
	"context"
	"fmt"
)

func (r *mutationResolver) Login(ctx context.Context, input *model.LoginInput) (string, error) {
	items, _ := ctx.Value("context_items").(ContextItems)
	db := items.Database
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
		authorization_client, err := grpcclient.GetAuthorizationGrpcClient()
		if err != nil {
			return "", err
		} else {
			defer authorization_client.Close()
			return authorization_client.CreateJWT(*user.ID)
		}
	} else {
		return "", err
	}
}

func (r *mutationResolver) Signout(ctx context.Context) (bool, error) {
	items, _ := ctx.Value("context_items").(ContextItems)
	db := items.Database
	defer db.Disconnect()
	sessionid := items.Sessionid
	authorization_client, err := grpcclient.GetAuthorizationGrpcClient()
	if err != nil {
		return false, err
	} else {
		defer authorization_client.Close()
		return authorization_client.DeleteJWT(*sessionid)
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

func (r *mutationResolver) UpdateUserName(ctx context.Context, newusername string) (bool, error) {
	context, _ := ctx.Value("context_items").(ContextItems)
	db := context.Database
	defer db.Disconnect()
	authorization_client, err := grpcclient.GetAuthorizationGrpcClient()
	if err != nil {
		return false, err
	} else {
		defer authorization_client.Close()
		res, err := authorization_client.ValidateJWT(*context.Sessionid)
		if err != nil {
			return false, err
		} else {
			return db.FindOneAndUpdateUsername(res, newusername)
		}
	}
}

func (r *queryResolver) Getuser(ctx context.Context) (*model.User, error) {
	context, _ := ctx.Value("context_items").(ContextItems)
	db := context.Database
	defer db.Disconnect()
	authorization_client, err := grpcclient.GetAuthorizationGrpcClient()
	if err != nil {
		return nil, err
	} else {
		defer authorization_client.Close()
		res, err := authorization_client.ValidateJWT(*context.Sessionid)
		if err != nil {
			return nil, err
		} else {
			user, err := db.FindOneByUserId(res)
			if err != nil {
				return nil, err
			} else {
				result := model.User{ID: user.ID, Username: user.Username, Phonenumber: user.Phonenumber}
				return &result, nil
			}
		}
	}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

type ContextItems struct {
	Sessionid *string
	Database  *database.DB
}
