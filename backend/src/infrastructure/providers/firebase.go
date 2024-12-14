package providers

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"google.golang.org/api/option"
)

type FirebaseProvider struct {
	App     *firebase.App
	Auth    *auth.Client
	Context context.Context
}

type IFirebaseProvider interface {
	GetUserFirebaseAuth(userID string) (auth.UserRecord, error)
}

func NewFirebaseProvider() IFirebaseProvider {
	context := context.Background()
	opt := option.WithAPIKey(os.Getenv("FIREBASE_API_KEY"))
	app, err := firebase.NewApp(context, nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
	}
	auth, err := app.Auth(context)
	if err != nil {
		fmt.Printf("error initializing auth: %v\n", err)
	}
	return &FirebaseProvider{
		App:     app,
		Auth:    auth,
		Context: context,
	}
}

func (f *FirebaseProvider) GetUserFirebaseAuth(userID string) (auth.UserRecord, error) {
	data, err := f.Auth.GetUser(f.Context, userID)
	if err != nil {
		return *data, err
	}
	if data == nil {
		return auth.UserRecord{}, fmt.Errorf("user not found")
	}
	return *data, nil
}
