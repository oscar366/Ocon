package helpers

import (
	"errors"
	"fmt"
	"time"
)
//THIS errors.go is for printing errors
func NewError(msg string) error {
	now := time.Now().Format(time.UnixDate)//get time
	fmt.Printf("\033[31m%s at: %s\033[0m\n", msg, now)//printerr
	return errors.New(msg)//return error
}