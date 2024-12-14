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
	// Get service account JSON from environment variable
	firebaseConfig := os.Getenv("FIREBASE_SERVICE_ACCOUNT")
	if firebaseConfig == "" {
		fmt.Println("FIREBASE_SERVICE_ACCOUNT environment variable not set")

	}
	// Convert JSON string to option
	creds := []byte(firebaseConfig)
	opt := option.WithCredentialsJSON(creds)
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
		fmt.Println(err)
		return *data, err
	}
	if data == nil {
		return auth.UserRecord{}, fmt.Errorf("user not found")
	}
	return *data, nil
}
