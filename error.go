package supabasestorageuploader

import "errors"

var (
	ErrFileNotFound     = errors.New("fileHeader is null")
	ErrFileNotInStorage = errors.New("file not found, check your storage name, file path, and file name")
	ErrLinkNotFound     = errors.New("file not found, check your storage name, file path, file name, and policy")
	ErrBadRequest       = errors.New("received bad request, check your request body")
)
