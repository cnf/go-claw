package core

import "testing"

func testValidate(t *testing.T, str string, expecterr bool) {
	err := ValidateName(str)

	if (expecterr) && (err == nil) {
		t.Errorf("Name '%s' should not validate but does!", str)
	} else if (!expecterr) && (err != nil) {
		t.Errorf("Name '%s' should validate but does not: %s!", str, err.Error())
	}
}

func Test_Validate(t *testing.T) {
	testValidate(t, "", true)
	testValidate(t, "a@bc", true)
	testValidate(t, "A#bc", true)
	testValidate(t, "A$bc", true)
	testValidate(t, "A!bc", true)
	testValidate(t, "A()bc", true)
	testValidate(t, "A&bc", true)
	testValidate(t, "A%bc", true)
	testValidate(t, "A[]bc", true)

	testValidate(t, "abcABC", false)
	testValidate(t, "abcABC_123", false)
	testValidate(t, "abcABC-123", false)
	testValidate(t, "abcABC-123_test", false)
}
