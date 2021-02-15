package cmd_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/cmd"
)

func TestJAID_GetID(t *testing.T) {
	t.Parallel()

	jaid := cmd.JAID{ID: "1"}

	assert.Equal(t, jaid.GetID(), "1")
}

func TestJAID_SetID(t *testing.T) {
	t.Parallel()

	jaid := cmd.JAID{}
	jaid.SetID("1")

	assert.Equal(t, jaid.ID, "1")
}

func TestJobType_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, cmd.DirectRequestJob.String(), "directrequest")
}

func TestJob_GetName(t *testing.T) {
	t.Parallel()

	job := &cmd.Job{}

	assert.Equal(t, job.GetName(), "specDBs")
}

func TestJob_GetTasks(t *testing.T) {
	t.Parallel()

	job := &cmd.Job{
		PipelineSpec: cmd.PipelineSpec{
			DotDAGSource: "    ds1          [type=http method=GET url=\"example.com\" allowunrestrictednetworkaccess=\"true\"];\n    ds1_parse    [type=jsonparse path=\"USD\"];\n    ds1_multiply [type=multiply times=100];\n    ds1 -\u003e ds1_parse -\u003e ds1_multiply;\n",
		},
	}

	tasks, err := job.GetTasks()

	assert.NoError(t, err)
	assert.Equal(t, tasks, []string{
		"ds1 http",
		"ds1_parse jsonparse",
		"ds1_multiply multiply",
	})
}

func TestJob_GetTasks_ParseError(t *testing.T) {
	t.Parallel()

	job := &cmd.Job{
		PipelineSpec: cmd.PipelineSpec{
			DotDAGSource: "invalid dot",
		},
	}

	tasks, err := job.GetTasks()

	assert.Error(t, err)
	assert.Nil(t, tasks)
}

func TestJob_FriendlyTasks(t *testing.T) {
	t.Parallel()

	job := &cmd.Job{
		PipelineSpec: cmd.PipelineSpec{
			DotDAGSource: "    ds1          [type=http method=GET url=\"example.com\" allowunrestrictednetworkaccess=\"true\"];\n    ds1_parse    [type=jsonparse path=\"USD\"];\n    ds1_multiply [type=multiply times=100];\n    ds1 -\u003e ds1_parse -\u003e ds1_multiply;\n",
		},
	}

	assert.Equal(t, job.FriendlyTasks(), []string{
		"ds1 http",
		"ds1_parse jsonparse",
		"ds1_multiply multiply",
	})
}

func TestJob_FriendlyTasks_ParseError(t *testing.T) {
	t.Parallel()

	job := &cmd.Job{
		PipelineSpec: cmd.PipelineSpec{
			DotDAGSource: "invalid dot",
		},
	}

	assert.Equal(t, job.FriendlyTasks(), []string{"error parsing DAG"})
}

func TestJob_FriendlyCreatedAt_DirectRequest(t *testing.T) {
	t.Parallel()

	now := time.Now()

	job := &cmd.Job{
		Type: cmd.DirectRequestJob,
		DirectRequestSpec: &cmd.DirectRequestSpec{
			CreatedAt: now,
		},
	}

	assert.Equal(t, job.FriendlyCreatedAt(), now.Format(time.RFC3339))
}

func TestJob_FriendlyCreatedAt_FluxMonitor(t *testing.T) {
	t.Parallel()

	now := time.Now()

	job := &cmd.Job{
		Type: cmd.FluxMonitorJob,
		FluxMonitorSpec: &cmd.FluxMonitorSpec{
			CreatedAt: now,
		},
	}

	assert.Equal(t, job.FriendlyCreatedAt(), now.Format(time.RFC3339))
}

func TestJob_FriendlyCreatedAt_OffChainReporting(t *testing.T) {
	t.Parallel()

	now := time.Now()

	job := &cmd.Job{
		Type: cmd.OffChainReportingJob,
		OffChainReportingSpec: &cmd.OffChainReportingSpec{
			CreatedAt: now,
		},
	}

	assert.Equal(t, job.FriendlyCreatedAt(), now.Format(time.RFC3339))
}

func TestJob_FriendlyCreatedAt_InvalidType(t *testing.T) {
	t.Parallel()

	job := &cmd.Job{
		Type: "invalid type",
	}

	assert.Equal(t, job.FriendlyCreatedAt(), "unknown")
}

func TestJob_FriendlyCreatedAt_NoSpecExists(t *testing.T) {
	t.Parallel()

	job := &cmd.Job{
		Type: cmd.DirectRequestJob,
	}

	assert.Equal(t, job.FriendlyCreatedAt(), "N/A")
}

func TestJob_ToRow(t *testing.T) {
	t.Parallel()

	now := time.Now()

	job := &cmd.Job{
		JAID: cmd.JAID{ID: "1"},
		Name: "Test Job",
		Type: cmd.DirectRequestJob,
		DirectRequestSpec: &cmd.DirectRequestSpec{
			CreatedAt: now,
		},
		PipelineSpec: cmd.PipelineSpec{
			DotDAGSource: "    ds1          [type=http method=GET url=\"example.com\" allowunrestrictednetworkaccess=\"true\"];\n    ds1_parse    [type=jsonparse path=\"USD\"];\n    ds1_multiply [type=multiply times=100];\n    ds1 -\u003e ds1_parse -\u003e ds1_multiply;\n",
		},
	}

	assert.Equal(t, job.ToRow(), [][]string{
		{"1", "Test Job", "directrequest", "ds1 http", now.Format(time.RFC3339)},
		{"1", "Test Job", "directrequest", "ds1_parse jsonparse", now.Format(time.RFC3339)},
		{"1", "Test Job", "directrequest", "ds1_multiply multiply", now.Format(time.RFC3339)},
	})
}
