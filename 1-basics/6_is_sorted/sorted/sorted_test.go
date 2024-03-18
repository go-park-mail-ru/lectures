package sorted

import "testing"

func TestCheckSuccess(t *testing.T) {
	var sortedInput = []string{"1", "2", "3", "4", "5"}

	err := Check(sortedInput)
	if err != nil {
		t.Fatalf("Check failed on sorted input: %s", err)
	}
}

func TestBlabla(t *testing.T) {
	t.Error(123)
}

func TestCheckFail(t *testing.T) {
	var mixedInput = []string{"1", "5", "3", "2", "4"}

	err := Check(mixedInput)
	if err == nil {
		t.Fatalf("Check not failed on mixed input: %s", err)
	}
}
