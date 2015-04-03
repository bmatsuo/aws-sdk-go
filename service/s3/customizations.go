package s3

import "github.com/awslabs/aws-sdk-go/aws"

func init() {
	initService = func(s *aws.Service) {
		// Some S3 operations require Content-MD5 to be set
		s.Handlers.Build.PushBack(contentMD5)

		// S3 uses custom error unmarshaling logic
		s.Handlers.UnmarshalError.Init()
		s.Handlers.UnmarshalError.PushBack(unmarshalError)
	}
}
