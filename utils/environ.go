package utils

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/djmarrerajr/common-lib/errs"
)

var (
	placeholderRegex = regexp.MustCompile(`\$?{(.*?)}`)
)

// Convenience wrapper for our os.Environ that allows us to
// add some useful functionality.
type Environ struct {
	envMap map[string]string
}

func LoadEnv(path string) (Environ, error) {
	envName, exists := os.LookupEnv("ENV")
	if !exists {
		envName = "local"
	}

	err := godotenv.Load(path + "/.env." + envName)
	if err != nil {
		return Environ{}, errs.Wrapf(err, errs.ErrTypeConfiguration, "unable to load .env for %s", envName)
	}

	environ := parseEnviron()

	return environ, nil
}

func (e Environ) Get(key string) (val string, OK bool) {
	val, OK = e.envMap[key]
	return
}

func (e Environ) GetRequired(key string) (val string, err error) {
	return require[string](key, e.Get)
}

func (e Environ) GetInt(key string) (val int, OK bool, err error) {
	value, OK := e.Get(key)
	if OK {
		val, err = strconv.Atoi(value)
		if err != nil {
			err = errs.Wrapf(err, errs.ErrTypeInvalidNumber, "environment variable %s with value of '%s' is not a valid int", key, value)
		}
	}

	return
}

func (e Environ) GetRequiredInt(key string) (int, error) {
	return require[int](key, e.GetInt)
}

func (e Environ) GetBool(key string) (val bool, OK bool, err error) {
	value, OK := e.Get(key)
	if OK {
		switch strings.ToLower(value) {
		case "true", "t", "1", "yes", "y":
			val = true
		case "false", "f", "0", "no", "n":
			val = false
		default:
			err = errs.Errorf(errs.ErrTypeInvalidBoolean, "environment variable %s with value of '%s' is not a valid boolean", key, value)
		}
	}

	return
}

func (e Environ) GetRequiredBool(key string) (bool, error) {
	return require[bool](key, e.GetBool)
}

func GetEnviron() Environ {
	envMap := make(map[string]string)

	for _, item := range os.Environ() {
		split := strings.SplitN(item, "=", 2)
		envMap[split[0]] = split[1]
	}

	return NewEnviron(envMap)
}

func NewEnviron(envMap map[string]string) Environ {
	return Environ{envMap}
}

// require is a type agnostic function that will attempt to retrieve the
// specified key from the environment returning either the value or an error
func require[T any](key string, getterFunc any) (val T, err error) {
	var OK bool

	switch funk := getterFunc.(type) {
	case func(string) (T, bool, error):
		val, OK, err = funk(key)
	case func(string) (T, bool):
		val, OK = funk(key)
	}

	if err == nil && !OK {
		err = errs.Errorf(errs.ErrTypeValidation, "missing required env key: %s", key)
	}

	return
}

// parseEnviron loads the current list of environment variables and their values and
// performs any necessary parameter substitution allowing us to create values that
// reference other env variables
//
// e.g.
//
//		MY_HOME=/users/home/dan
//	 SOMEVAR=${MY_HOME}/some/other/path
//
//	 This function would update the environment variables so they become:
//
//		MY_HOME=/users/home/dan
//	 SOMEVAR=/users/home/dan/some/other/path
//
// This allows us to create more 'portable' .env files
func parseEnviron() Environ {
	osEnv := os.Environ()

	for _, item := range osEnv {
		matches := placeholderRegex.FindAllString(item, -1)
		if len(matches) > 0 {
			for _, match := range matches {
				val := getValueForKey(match)
				sub := regexp.MustCompile(`(?i)\` + match)
				new := sub.ReplaceAllLiteralString(item, val)

				newKeyVal := strings.Split(new, "=")
				os.Setenv(newKeyVal[0], newKeyVal[1])
			}
		}
	}

	return GetEnviron()
}

// getValueForKey will obtain and return the value for the specified
// key as it currently is within the os.Environ()
func getValueForKey(placeholder string) string {
	key := strings.ToUpper(placeholder[1 : len(placeholder)-1])
	if strings.HasPrefix(placeholder, "$") {
		key = strings.ToUpper(placeholder[2 : len(placeholder)-1])
	}

	val, exists := os.LookupEnv(key)
	if !exists {
		val = os.Getenv(strings.ToLower(key))
	}

	return val
}
