package utils_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"

	"github.com/djmarrerajr/common-lib/errs"
	"github.com/djmarrerajr/common-lib/utils"
)

const envFilename = ".env.local"

type EnvironTestSuite struct {
	suite.Suite
}

func (e *EnvironTestSuite) SetupTest() {
	os.Clearenv()
	os.Remove(envFilename)
}

func (l *EnvironTestSuite) TeardownTest() {
	os.Remove(envFilename)
}

func (e *EnvironTestSuite) TestLoadEnv_NoFile_ShouldReturnConfigurationError() {
	_, err := utils.LoadEnv("")

	e.Error(err)
	e.Equal(errs.GetType(err), errs.ErrTypeConfiguration)
}

func (e *EnvironTestSuite) TestGet_NoPlaceholder_ShouldReturnExpectedValue() {
	key, val := "somekey", "somevalue"

	e.createEnvFile([]string{key + "=" + val})

	env, err := utils.LoadEnv("./")

	v, exists := env.Get(key)

	e.NoError(err)
	e.True(exists)
	e.Equal(val, v)
}

func (e *EnvironTestSuite) TestGet_Placeholder_ShouldReturnExpectedValue() {
	format := "/%s/repos"
	envKey, envVal := "repo_root", fmt.Sprintf(format, "{my_home}")
	subKey, subVal := "my_home", "bruno"

	e.createEnvFile([]string{envKey + "=" + envVal, subKey + "=" + subVal})

	env, err := utils.LoadEnv("./")

	v, exists := env.Get(envKey)

	e.NoError(err)
	e.True(exists)
	e.Equal(fmt.Sprintf(format, subVal), v)
}

func (e *EnvironTestSuite) TestGet_PlaceholderWithDollarSign_ShouldReturnExpectedValue() {
	format := "/%s/repos"
	envKey, envVal := "repo_root", fmt.Sprintf(format, "${my_home}")
	subKey, subVal := "my_home", "bruno"

	e.createEnvFile([]string{envKey + "=" + envVal, subKey + "=" + subVal})

	env, err := utils.LoadEnv("./")

	v, exists := env.Get(envKey)

	e.NoError(err)
	e.True(exists)
	e.Equal(fmt.Sprintf(format, subVal), v)
}

// func (e *EnvironTestSuite) TestGetRequired_ValueExists() {
// 	key, val := "somekey", "somevalue"

// 	e.T().Setenv(key, val)

// 	env := utils.GetEnviron()

// 	v, err := env.GetRequired(key)

// 	e.NoError(err)
// 	e.Equal(val, v)
// }

func (e *EnvironTestSuite) TestGetRequired_ValueDoesNotExists() {
	_, err := utils.GetEnviron().GetRequired("missingkey")

	e.Error(err)
	e.Equal(errs.ErrTypeValidation, errs.GetType(err))
}

func (e *EnvironTestSuite) TestGetInt_ValueExists() {
	key, val := "intkey", 123

	e.T().Setenv(key, strconv.Itoa(val))

	v, exists, err := utils.GetEnviron().GetInt(key)

	e.NoError(err)
	e.True(exists)
	e.Equal(val, v)
}

func (e *EnvironTestSuite) TestGetRequiredInt_ValueExists() {
	key, val := "intkey", 123

	e.T().Setenv(key, strconv.Itoa(val))

	v, err := utils.GetEnviron().GetRequiredInt(key)

	e.NoError(err)
	e.Equal(val, v)
}

func (e *EnvironTestSuite) TestGetRequiredInt_ValueIsInvalid() {
	key, val := "intkey", "notanint"

	e.T().Setenv(key, val)

	_, err := utils.GetEnviron().GetRequiredInt(key)

	e.Error(err)
	e.Equal(errs.ErrTypeInvalidNumber, errs.GetType(err))
}

func (e *EnvironTestSuite) TestGetRequiredInt_ValueDoesNotExist() {
	_, err := utils.GetEnviron().GetRequiredInt("missingint")

	e.Error(err)
	e.Equal(errs.ErrTypeValidation, errs.GetType(err))
}

func (e *EnvironTestSuite) TestGetBool_ValuesExist() {
	expectedErr := errs.New(errs.ErrTypeInvalidBoolean, "bad boolean")

	tests := []lo.Tuple3[string, error, bool]{
		{A: "true", B: nil, C: true},
		{A: "t", B: nil, C: true},
		{A: "1", B: nil, C: true},
		{A: "yes", B: nil, C: true},
		{A: "y", B: nil, C: true},
		{A: "false", B: nil, C: false},
		{A: "f", B: nil, C: false},
		{A: "0", B: nil, C: false},
		{A: "no", B: nil, C: false},
		{A: "n", B: nil, C: false},
		{A: "x", B: expectedErr, C: false},
	}

	for _, test := range tests {
		key := "testbool"

		e.Run(test.A, func() {
			e.T().Setenv(key, test.A)

			v, _, err := utils.GetEnviron().GetBool(key)

			if test.B != nil {
				e.Error(err)
				e.Equal(errs.ErrTypeInvalidBoolean, errs.GetType(err))
			} else {
				e.NoError(err)
			}

			e.Equal(test.C, v)
		})
	}
}

func (e *EnvironTestSuite) TestGetRequiredBool_ValuesExist() {
	expectedErr := errs.New(errs.ErrTypeInvalidBoolean, "bad boolean")

	tests := []lo.Tuple3[string, error, bool]{
		{A: "true", B: nil, C: true},
		{A: "t", B: nil, C: true},
		{A: "1", B: nil, C: true},
		{A: "yes", B: nil, C: true},
		{A: "y", B: nil, C: true},
		{A: "false", B: nil, C: false},
		{A: "f", B: nil, C: false},
		{A: "0", B: nil, C: false},
		{A: "no", B: nil, C: false},
		{A: "n", B: nil, C: false},
		{A: "x", B: expectedErr, C: false},
		{A: "", B: expectedErr, C: false},
	}

	for _, test := range tests {
		key := "testbool"

		e.Run(test.A, func() {
			e.T().Setenv(key, test.A)

			v, err := utils.GetEnviron().GetRequiredBool(key)

			if test.B != nil {
				e.Error(err)
				e.Equal(errs.ErrTypeInvalidBoolean, errs.GetType(err))

				fmt.Println(err.Error())
			} else {
				e.NoError(err)
			}

			e.Equal(test.C, v)
		})
	}
}

func (e *EnvironTestSuite) createEnvFile(env []string) {
	file, err := os.Create(envFilename)
	e.NoError(err, "unable to create the file needed to simulate an environment")

	defer file.Close()

	// nolint: errcheck
	for _, item := range env {
		file.WriteString(item + "\n")
	}
}

func TestEnviron(t *testing.T) {
	suite.Run(t, new(EnvironTestSuite))
}
