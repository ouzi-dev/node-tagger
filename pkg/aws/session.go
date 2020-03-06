package aws

import "github.com/aws/aws-sdk-go/aws/session"

// Gets the aws session to use for looking up credstash secrets falling back to the environment config
func GetAwsSessionFromEnv() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return nil, err
	}

	return sess, nil
}
