package stats

import "errors"

var DataValidationError = errors.New("validation failed")
var DataLoadingError = errors.New("loading failed")
var DataParsingError = errors.New("parsing failed")
