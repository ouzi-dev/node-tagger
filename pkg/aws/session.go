package aws

import "github.com/aws/aws-sdk-go/aws/session"

// Gets the aws session
func GetAwsSessionFromEnv() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return nil, err
	}

	return sess, nil
}
