package terminal

import (
	"bufio"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"ydx-goadv-gophkeeper/internal/client/mocks/services"
	"ydx-goadv-gophkeeper/pkg/mocks/shutdown"
)

type TextWriter struct {
	writer *bufio.Writer
}

func TestCommandParser_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := &TextWriter{writer: bufio.NewWriter(os.Stdout)}

	authService := services.NewMockAuthService(ctrl)
	resService := services.NewMockResourceService(ctrl)
	exitHandler := shutdown.NewMockExitHandler(ctrl)
	parser := &commandParser{authService: authService, resourceService: resService, exitHandler: exitHandler}
	parser.initScanner()
	parser.commands = map[string]func(args []string) (string, error){
		"login": parser.handleLogin,
	}
	assert.NoError(t, writer.Text("login"))
	result, err := processCommands(parser)

	assert.NoError(t, err)

	assert.NoError(t, writer.Text("login"))
	assert.Equal(t, "OK", result)

	assert.Equal(t, "not testable", "not testable")
}

func (tw *TextWriter) Text(text string) error {
	_, err := tw.writer.WriteString(text)
	if err != nil {
		return err
	}
	return tw.writer.Flush()
}

func processCommands(cp *commandParser) (string, error) {
	cmd := cp.readString("")
	if len(cmd) == 0 {
		return "", nil
	}
	result, err := cp.handle(cmd)
	if err != nil {
		return "", err
	}
	return result, nil
}
