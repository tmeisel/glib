package writer

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ctxPkg "github.com/tmeisel/glib/ctx"
	logPkg "github.com/tmeisel/glib/log"
	"github.com/tmeisel/glib/log/fields"
)

var (
	log *Writer
)

func TestNew(t *testing.T) {
	writer := new(bytes.Buffer)
	log = New(writer, false, logPkg.LevelDebug)

	assert.Equal(t, writer, log.writer)
	assert.Equal(t, logPkg.LevelDebug, log.level)
	assert.Equal(t, false, log.production)
}

func TestNewStdWriter(t *testing.T) {
	log = NewStdWriter(true, logPkg.LevelDebug)
	assert.Equal(t, os.Stdout, log.writer)
}

func TestWriter_Debug(t *testing.T) {
	buf := new(bytes.Buffer)
	log = New(buf, false, logPkg.LevelDebug)

	ctx := ctxPkg.WithLogFields(context.Background(), fields.Bool("context", true))

	msg := "this is a debug message"
	log.Debug(ctx, msg)

	output := buf.String()

	assert.Contains(t, output, "[debug]")
	assert.Contains(t, output, msg)

	// assert the timestamp is correct. to avoid a flaky test, when it's run close to
	// the end or beginning of a minute (day,month,...), it's parsing the timestamp
	// and comparing it to the current time
	dateTimeRegexp := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`)
	match := dateTimeRegexp.FindString(output)

	require.NotEmpty(t, match, fmt.Sprintf("output '%s' does not contain ''", output))

	matchedTime, err := time.ParseInLocation("2006-01-02T15:04:05", match, time.Local)
	if err != nil {
		t.Error(err)
		return
	}

	assert.WithinDuration(t, time.Now(), matchedTime, time.Second)
}

func TestWriter_Info(t *testing.T) {
	buf := new(bytes.Buffer)
	log = New(buf, false, logPkg.LevelInfo)
	ctx := ctxPkg.WithLogFields(context.Background(), fields.Bool("context", true))
	msg := "this is a info message"

	log.Info(ctx, msg)
	log.Debug(ctx, msg)

	output := buf.String()

	assert.Contains(t, output, msg)
	assert.Contains(t, output, "[info]")
	assert.NotContains(t, output, "[debug]")
}

func TestWriter_SetLevel(t *testing.T) {
	buf := new(bytes.Buffer)
	log = New(buf, false, logPkg.LevelDebug)
	log.Debug(context.Background(), "first")

	require.NoError(t, log.SetLevel(logPkg.LevelInfo))

	log.Debug(context.Background(), "second")

	assert.Contains(t, buf.String(), "first")
	assert.NotContains(t, buf.String(), "second")
}
