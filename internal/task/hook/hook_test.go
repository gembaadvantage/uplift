package hook

import (
	"context"
	"os"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec_DryRun(t *testing.T) {
	gittest.InitRepository(t)

	err := Exec(context.Background(), []string{"touch out.txt"}, ExecOptions{DryRun: true})
	require.NoError(t, err)

	assert.NoFileExists(t, "out.txt")
}

func TestExec_InjectEnvVars(t *testing.T) {
	gittest.InitRepository(t)

	env := []string{"VARIABLE=VALUE", "ANOTHER_VARIABLE=ANOTHER VALUE"}

	cmds := []string{
		"echo -n VALUE1=$VARIABLE VALUE2=$ANOTHER_VARIABLE > out.txt",
	}

	err := Exec(context.Background(), cmds, ExecOptions{Env: env})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "VALUE1=VALUE VALUE2=ANOTHER VALUE", string(data))
}

func TestExec_MergesEnvVars(t *testing.T) {
	gittest.InitRepository(t)

	env := []string{"TESTING=123"}

	cmds := []string{
		"printenv > out.txt",
	}

	err := Exec(context.Background(), cmds, ExecOptions{Env: env})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Contains(t, string(data), "PATH=")
	assert.Contains(t, string(data), "PWD=")
	assert.Contains(t, string(data), "HOME=")
	assert.Contains(t, string(data), "TESTING=123")
}

func TestExec_VarsWithWhitespace(t *testing.T) {
	gittest.InitRepository(t)

	env := []string{"ONE = 1", "TWO= 2", "THREE    =    3"}

	cmds := []string{
		"echo -n $ONE $TWO $THREE > out.txt",
	}

	err := Exec(context.Background(), cmds, ExecOptions{Env: env})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "1 2 3", string(data))
}

func TestExec_LoadDotEnvFiles(t *testing.T) {
	gittest.InitRepository(t)

	dotenv1 := `ONE=1
TWO   =   2`
	os.WriteFile(".env", []byte(dotenv1), 0o600)

	os.Mkdir("custom", 0o755)
	dotenv2 := "THREE=    3"
	os.WriteFile("custom/another.env", []byte(dotenv2), 0o600)

	cmds := []string{
		"echo -n $ONE $TWO $THREE > out.txt",
	}

	err := Exec(context.Background(), cmds, ExecOptions{Env: []string{".env", "custom/another.env"}})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "1 2 3", string(data))
}

func TestExec_FailsOnInvalidDotEnvFile(t *testing.T) {
	gittest.InitRepository(t)
	gittest.TempFile(t, ".env", "INVALID")

	err := Exec(context.Background(), []string{}, ExecOptions{Env: []string{".env"}})
	require.Error(t, err)
}

func TestExec_FailsIfDotEnvFileNotFound(t *testing.T) {
	gittest.InitRepository(t)

	err := Exec(context.Background(), []string{}, ExecOptions{Env: []string{"does-not-exist.env"}})
	require.Error(t, err)
}
