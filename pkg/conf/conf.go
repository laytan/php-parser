package conf

import (
	"github.com/laytan/php-parser/pkg/errors"
	"github.com/laytan/php-parser/pkg/version"
)

type Config struct {
	Version          *version.Version
	ErrorHandlerFunc func(e *errors.Error)
}
