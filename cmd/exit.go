package cmd

import (
	"fmt"
	"os"
	"strings"
)

const (
	// lifecycle errors not specific to any phase: 1-19
	CodeFailed = 1 // CodeFailed indicates generic lifecycle error
	// 2: reserved
	CodeInvalidArgs = 3
	// 4: CodeInvalidEnv
	// 5: CodeNotFound
	// 9: CodeFailedUpdate

	// API errors
	CodeIncompatiblePlatformAPI  = 11
	CodeIncompatibleBuildpackAPI = 12

	// detect phase errors: 20-29
	CodeFailedDetect = 20 // CodeFailedDetect indicates that no buildpacks detected
	// CodeFailedDetectWithErrors indicated that no buildpacks detected and at least one errored
	CodeFailedDetectWithErrors = 21
	CodeDetectError            = 22 // CodeDetectError indicates generic detect error

	// analyze phase errors: 30-39
	CodeAnalyzeError = 32 // CodeAnalyzeError indicates generic analyze error

	// restore phase errors: 40-49
	CodeRestoreError = 42 // CodeRestoreError indicates generic restore error

	// build phase errors: 50-59
	CodeFailedBuildWithErrors = 51 // CodeFailedBuildWithErrors indicates buildpack error during /bin/build
	CodeBuildError            = 52 // CodeBuildError indicates generic build error

	// export phase errors: 60-69
	CodeExportError = 62 // CodeExportError indicates generic export error

	// rebase phase errors: 70-79
	CodeRebaseError = 72 // CodeRebaseError indicates generic rebase error

	// launch phase errors: 80-89
	CodeLaunchError = 82 // CodeLaunchError indicates generic launch error
)

type ErrorFail struct {
	Err    error
	Code   int
	Action []string
}

func (e *ErrorFail) Error() string {
	message := "failed to " + strings.Join(e.Action, " ")
	if e.Err == nil {
		return message
	}
	return fmt.Sprintf("%s: %s", message, e.Err)
}

func FailCode(code int, action ...string) *ErrorFail {
	return FailErrCode(nil, code, action...)
}

func FailErr(err error, action ...string) *ErrorFail {
	code := CodeFailed
	if err, ok := err.(*ErrorFail); ok {
		code = err.Code
	}
	return FailErrCode(err, code, action...)
}

func FailErrCode(err error, code int, action ...string) *ErrorFail {
	return &ErrorFail{Err: err, Code: code, Action: action}
}

func Exit(err error) {
	if err == nil {
		os.Exit(0)
	}
	DefaultLogger.Errorf("%s\n", err)
	if err, ok := err.(*ErrorFail); ok {
		os.Exit(err.Code)
	}
	os.Exit(CodeFailed)
}

func ExitWithVersion() {
	DefaultLogger.Infof(buildVersion())
	os.Exit(0)
}
