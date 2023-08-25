package supabasestorageuploader

import "errors"

var (
	errFileNotFound     = errors.New("fileHeader is null")
	errFileNotInStorage = errors.New("file not found, check your storage name, file path, and file name")
	errLinkNotFound     = errors.New("file not found, check your storage name, file path, file name, and policy")
)
