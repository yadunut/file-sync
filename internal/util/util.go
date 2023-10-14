package util

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

// returns true if either of dirs is a child of the other. This is to ensure that we don't have issues like
// ~/Documents and ~/Documents/folder both being watched
// This function also expects absolute paths.
func IsDirChild(dirA string, dirB string) bool {
	aLen := len(strings.Split(dirA, string(os.PathSeparator)))
	bLen := len(strings.Split(dirB, string(os.PathSeparator)))
	if aLen > bLen {
		dirA, dirB = dirB, dirA
	}
	return strings.Contains(dirB, dirA)
}

func CreateLogger() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}
