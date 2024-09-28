package authentication

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/NUS-EVCHARGE/ev-user-service/dto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/sirupsen/logrus"
	"log"
)

type AuthenticationController interface {
	RegisterUser(signUpRequest dto.SignUpRequest) error
	LoginUser(credentials dto.Credentials) (dto.LoginResponse, error)
	ConfirmUser(userInfo dto.ConfirmUser) error
	ResendChallengeCode(resendRequest dto.SignUpResendRequest) error
	LogoutUser(accessToken string) error
	GetUserInfo(accessToken string) (*cognitoidentityprovider.GetUserOutput, error)
	AssociateMFADevice(accessToken string) (*cognitoidentityprovider.AssociateSoftwareTokenOutput, error)
	VerifyMFADevice(accessToken, session, userCode string) error
}

type AuthenticationControllerImpl struct {

}

var (
	awsRegion = "ap-southeast-1" // Your AWS Region
	clientID = "5tbesavn2ita7r8anua8d5k901"
	clientSecret  = "1nc806v9oi3rgocov4p4uqthq1qs5ficstn8e39876q9rqrhk84o"
	AuthenticationControllerObj AuthenticationController
	cognitoClient = setupCognitoClient()
)

func NewAuthenticationController() {
	AuthenticationControllerObj = &AuthenticationControllerImpl{}
}


func (a AuthenticationControllerImpl) RegisterUser(signUpRequest dto.SignUpRequest) error {
	secretHash := generateSecretHash(clientSecret, signUpRequest.Email, clientID)

	//preferredUsername := signUpRequest.Email
	//if atIndex := strings.Index(preferredUsername, "@"); atIndex != -1 {
	//	preferredUsername = preferredUsername[:atIndex]
	//}

	input := &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(clientID),
		Username:   aws.String(signUpRequest.Email),
		Password:   aws.String(signUpRequest.Password),
		SecretHash: aws.String(secretHash),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(signUpRequest.Email),
			},
			{
				Name:  aws.String("preferred_username"),
				Value: aws.String(signUpRequest.Email),
			},
		},
	}

	result, err := cognitoClient.SignUp(input)
	if err != nil {
		return err
	}

	fmt.Printf("User %s registered successfully\n", *result.UserSub)
	return nil
}

func (a AuthenticationControllerImpl) LoginUser(credentials dto.Credentials) (dto.LoginResponse, error) {
	secretHash := generateSecretHash(clientSecret, credentials.Email, clientID)

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(credentials.Email),
			"PASSWORD":    aws.String(credentials.Password),
			"SECRET_HASH": aws.String(secretHash),
		},
		ClientId: aws.String(clientID),
	}
	login := dto.LoginResponse{}

	authResp, err := cognitoClient.InitiateAuth(input)
	if err != nil {
		login.Status = "failed"
		return login , err
	}

	// Handle MFA setup and challenges
	if authResp.ChallengeName != nil && *authResp.ChallengeName != "" {
		fmt.Println("Challenge required:", authResp.ChallengeName)
		logrus.WithField("login", login).Info("User authentication failed \n", credentials.Email)
		login.Status = *authResp.ChallengeName
		//handleChallenge(cognitoClient, clientID, clientSecret, loginCredential.Email, *authResp.Session, *authResp.ChallengeName)
	} else {
		login.AccessToken = *authResp.AuthenticationResult.AccessToken
		login.RefreshToken = *authResp.AuthenticationResult.RefreshToken
		login.IdToken = *authResp.AuthenticationResult.IdToken
		login.ExpiresIn = int(*authResp.AuthenticationResult.ExpiresIn)
		login.Status = "success"
	}

	logrus.WithField("login", login).Info("User authenticated successfully \n", credentials.Email)
	return login, nil

}

func (a AuthenticationControllerImpl) ConfirmUser(userInfo dto.ConfirmUser) error {

	secretHash := generateSecretHash(clientSecret, userInfo.Email, clientID)

	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(clientID),
		Username:         aws.String(userInfo.Email),
		ConfirmationCode: aws.String(userInfo.ConfirmationCode),
		SecretHash:       aws.String(secretHash),
	}

	_, err := cognitoClient.ConfirmSignUp(input)
	if err != nil {
		return err
	}

	fmt.Printf("User %s confirmed successfully\n", userInfo.Email)

	return nil
}

func (a AuthenticationControllerImpl) ResendChallengeCode(resendRequest dto.SignUpResendRequest) error {
	secretHash := generateSecretHash(clientSecret, resendRequest.Email, clientID)

	input := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId:   aws.String(clientID),
		Username:   aws.String(resendRequest.Email),
		SecretHash: aws.String(secretHash),
	}

	_, err := cognitoClient.ResendConfirmationCode(input)
	if err != nil {
		return err
	}

	fmt.Printf("Confirmation code resent to %s\n", resendRequest.Email)

	return nil
}

func (a AuthenticationControllerImpl) LogoutUser(accessToken string) error {
	input := &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(accessToken),
	}

	_, err := cognitoClient.GlobalSignOut(input)
	if err != nil {
		return err
	}

	fmt.Printf("User %s logged out successfully\n", accessToken)
	return nil
}

func (a AuthenticationControllerImpl) GetUserInfo(accessToken string) (*cognitoidentityprovider.GetUserOutput, error) {
	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	}

	userInfo, err := cognitoClient.GetUser(input)
	if err != nil {
		return nil, err
	}

	logrus.WithField("accessToken", accessToken).Info("User authenticated successfully")
	return userInfo, nil
}

func (a AuthenticationControllerImpl) AssociateMFADevice(accessToken string) (*cognitoidentityprovider.AssociateSoftwareTokenOutput, error) {
input := &cognitoidentityprovider.AssociateSoftwareTokenInput{
AccessToken: aws.String(accessToken),
}

result, err := cognitoClient.AssociateSoftwareToken(input)
if err != nil {
return nil, err
}

return result, nil
}

func (a AuthenticationControllerImpl) VerifyMFADevice(accessToken, session, userCode string) error {
	input := &cognitoidentityprovider.VerifySoftwareTokenInput{
		AccessToken: aws.String(accessToken),
		Session:     aws.String(session),
		UserCode:    aws.String(userCode),
	}

	_, err := cognitoClient.VerifySoftwareToken(input)
	if err != nil {
		return err
	}

	return nil
}

func generateSecretHash(clientSecret, userName, clientID string) string {
	key := []byte(clientSecret)
	message := userName + clientID
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	hash := mac.Sum(nil)
	secretHash := base64.StdEncoding.EncodeToString(hash)
	return secretHash
}

func setupCognitoClient() *cognitoidentityprovider.CognitoIdentityProvider {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	if err != nil {
		log.Fatalf("failed to create session, %v", err)
	}

	// Create a Cognito Identity Provider client
	cognitoClient := cognitoidentityprovider.New(sess)
	return cognitoClient
}

func getCognitoClient() *cognitoidentityprovider.CognitoIdentityProvider {
	return cognitoClient
}