package logerr

import (
	"errors"
	"reflect"
	"testing"

	juju "github.com/juju/errors"
	pkg "github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
)

type ErrorTestSuite struct {
	suite.Suite
	m map[string]interface{}
}

func (s *ErrorTestSuite) SetupSuite() {
	s.m = map[string]interface{}{"testID": "0"}
}

func (s *ErrorTestSuite) TestWithField() {
	err := WithField(nil, "testID", "0")
	s.Nilf(err, "error expexted nil, was %#v", err)

	wrappedErr := errors.New("test error")
	err = WithField(wrappedErr, "testID", "0")

	s.NotNil(err, "error expected not nil, was %#v\n", err)
	lErr, ok := err.(*Error)
	if !ok {
		s.FailNowf("error type expexted *Error, was %s", reflect.TypeOf(err).Name())
	}
	s.Equal(s.m, lErr.Fields, "log error fields are not matching")
	s.Equal(wrappedErr, lErr.prev, "wrong wrapped error")
}

func (s *ErrorTestSuite) TestWithFields() {
	err := WithFields(nil, s.m)
	s.Nilf(err, "error expexted nil, was %#v", err)

	wrappedErr := errors.New("test error")
	err = WithFields(wrappedErr, s.m)

	s.NotNil(err, "error expected not nil, was %#v\n", err)
	lErr, ok := err.(*Error)
	if !ok {
		s.FailNowf("error type expexted *Error, was %s", reflect.TypeOf(err).Name())
	}
	s.Equal(s.m, lErr.Fields, "log error fields are not matching")
	s.Equal(wrappedErr, lErr.prev, "wrong wrapped error")
}

func (s *ErrorTestSuite) TestWrapper() {
	cause := errors.New("cause")
	lerr := WithFields(cause, s.m)
	wrapped := juju.Annotate(lerr, "wrapped")
	s.Equal("wrapped: cause", wrapped.Error(), "error strings not matching")

	fields := GetFields(wrapped)
	s.Equal(s.m, fields, "field map not matching")
}

func (s *ErrorTestSuite) TestCauser() {
	cause := errors.New("cause")
	lerr := WithFields(cause, s.m)
	wrapped := pkg.WithMessage(lerr, "wrapped")
	s.Equal("wrapped: cause", wrapped.Error(), "error strings not matching")

	fields := GetFields(wrapped)
	s.Equal(s.m, fields, "field map not matching")
}

func (s *ErrorTestSuite) TestDeferWtithFields() {
	var err error
	tFunc := func() {
		defer DeferWithFields(&err, s.m)
	}
	tFunc()
	s.Nilf(err, "error expexted nil, was %#v", err)

	first := errors.New("test error")
	err = first
	tFunc()

	s.NotNil(err, "error expected not nil, was %#v\n", err)
	lErr, ok := err.(*Error)
	if !ok {
		s.FailNowf("error type expexted *Error, was %s", reflect.TypeOf(err).Name())
	}
	s.Equal(s.m, lErr.Fields, "log error fields are not matching")
	s.Equal(first, lErr.prev, "wrong wrapped error")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestErrorTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}
