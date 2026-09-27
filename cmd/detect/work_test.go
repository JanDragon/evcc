package detect

import (
	"testing"

	"github.com/evcc-io/evcc/cmd/detect/tasks"
	"github.com/stretchr/testify/assert"
)

func TestPostProcessSortsByIPThenType(t *testing.T) {
	res := []tasks.Result{
		{Task: tasks.Task{ID: "b", Type: tasks.TaskType("modbus")}, ResultDetails: tasks.ResultDetails{IP: "192.168.0.11"}},
		{Task: tasks.Task{ID: "c", Type: tasks.TaskType("keba")}, ResultDetails: tasks.ResultDetails{IP: "192.168.0.10"}},
		{Task: tasks.Task{ID: "a", Type: tasks.TaskType("http")}, ResultDetails: tasks.ResultDetails{IP: "192.168.0.10"}},
	}

	got := postProcess(res)

	assert.Equal(t, []string{"192.168.0.10/http", "192.168.0.10/keba", "192.168.0.11/modbus"}, []string{
		got[0].IP + "/" + string(got[0].Type),
		got[1].IP + "/" + string(got[1].Type),
		got[2].IP + "/" + string(got[2].Type),
	})
}
