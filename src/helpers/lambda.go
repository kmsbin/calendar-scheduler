package helpers

import "os"

func IsRunningInLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}
