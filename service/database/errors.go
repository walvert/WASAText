package database

import "errors"

var ErrAlreadyExists = errors.New("user already exists in chat")
var ErrEmptyImageURL = errors.New("image URL is empty")
