package build

import (
	"testing"

	"github.com/MeteorsLiu/llarmvp/pkgs/formula/matrix"
)

func TestBuild(t *testing.T) {
	task, err := NewTask("DaveGamble/cJSON", matrix.Current(), "1.7.18")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(task.Exec())
}
