package terminal

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"

	"ydx-goadv-gophkeeper/internal/client/model/consts"
	"ydx-goadv-gophkeeper/internal/client/model/resources"
	"ydx-goadv-gophkeeper/internal/client/services"
	"ydx-goadv-gophkeeper/pkg/model/enum"
)

const (
	maxCapacity   = 1024 * 1024
	successResult = "success"
	helpMsg       = "" +
		"available commands:\n" +
		"	'clear' - to clear terminal\n" +
		"\n" +
		"	'login' - to login\n" +
		"	'register' - to register\n" +
		"\n" +
		"	's [type]' - save resource, where 'type' is: lp - LoginPassword, fl - File, bc - BankCard\n\t\n" +
		"\n" +
		"	'd [id]' - delete resource by id\n" +
		"	'l [type]' - get resources by type, where 'type' is:lp - LoginPassword, fl - File, bc - BankCard\n	or get all if type is empty\n" +
		"	'g [id]' - get loginPassword or BankCard by id\n" +
		"	'gf [id]' - get file by id\n"
)

type CommandParser interface {
	Start()
}

type commandParser struct {
	scanner         *bufio.Scanner
	authService     services.AuthService
	resourceService services.ResourceService
	commands        map[string]func(args []string) (string, error)
}

func NewCommandParser(
	buildVersion string,
	buildDate string,
	authService services.AuthService,
	resourceService services.ResourceService,
) CommandParser {
	fmt.Printf("buildVersion='%s' buildDate='%s'\n%s\n", buildVersion, buildDate, helpMsg)
	cp := &commandParser{
		authService:     authService,
		resourceService: resourceService,
	}
	cp.scanner = cp.initScanner()
	cp.commands = map[string]func(args []string) (string, error){
		"login":    cp.handleLogin,
		"register": cp.handleRegistration,
		"s":        cp.handleSave,
		"d":        cp.handleDelete,
		"l":        cp.handleList,
		"g":        cp.handleGet,
		"gf":       cp.handleGetFile,
		"clear":    cp.handleClear,
		"help":     cp.handleHelp,
	}
	return cp
}

func (cp *commandParser) Start() {
	for {
		cmd := cp.readString("")
		if len(cmd) == 0 {
			continue
		}

		result, err := cp.handle(cmd)
		if err != nil {
			fmt.Printf("error: %v\n", err.Error())
			continue
		}
		fmt.Printf("%s\n", result)
	}
}

func (cp *commandParser) handle(input string) (string, error) {
	arr := strings.Split(input, " ")

	command := arr[0]
	args := arr[1:]

	if f, ok := cp.commands[command]; ok {
		return f(args)
	}
	return "", fmt.Errorf("command '%s' is not supported, type 'help' to display available commands", command)

}

func (cp *commandParser) handleClear(_ []string) (string, error) {
	fmt.Print("\033[H\033[2J")
	return "", nil
}

func (cp *commandParser) handleHelp(_ []string) (string, error) {
	fmt.Print(helpMsg)
	return "", nil
}

func (cp *commandParser) handleLogin(_ []string) (string, error) {
	login := cp.readString("input username")
	password := cp.readPassword()
	_, err := cp.authService.Login(context.Background(), login, password)
	return successResult, err
}

func (cp *commandParser) handleRegistration(_ []string) (string, error) {
	login := cp.readString("input username")
	password := cp.readPassword()
	_, err := cp.authService.Register(context.Background(), login, password)
	return successResult, err
}

func (cp *commandParser) handleGetFile(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("arg '[id]' is empty, type 'help' to display available commands format")
	}
	resId, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		return "", err
	}
	path, err := cp.resourceService.GetFile(context.Background(), int32(resId))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("recieved file saved to: %v", path), nil
}

func (cp *commandParser) handleGet(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("arg '[id]' is empty, type 'help' to display available commands format")
	}
	resId, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		return "", err
	}
	resDescription, err := cp.resourceService.Get(context.Background(), int32(resId))
	if err != nil {
		return "", err
	}

	return resDescription.Format(), nil
}

func (cp *commandParser) handleList(args []string) (string, error) {
	resType := enum.Nan
	if len(args) != 0 {
		if rType, ok := consts.ArgToType[args[0]]; ok {
			resType = rType
		}
	}

	resDescriptions, err := cp.resourceService.GetDescriptions(context.Background(), resType)
	if err != nil {
		return "", err
	}
	var writer strings.Builder
	if len(resDescriptions) == 0 {
		_, err := writer.WriteString("empty")
		if err != nil {
			return "", err
		}
	}
	for _, resDescription := range resDescriptions {
		_, err := writer.WriteString(fmt.Sprintf("id: %d - type: '%s', descr: '%s'\n", resDescription.Id, consts.TypeToArg[resDescription.Type], string(resDescription.Meta)))
		if err != nil {
			return "", err
		}
	}
	return writer.String(), nil
}

func (cp *commandParser) handleDelete(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("arg '[id]' is empty, type 'help' to display available commands format")
	}
	resId, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		return "", err
	}
	err = cp.resourceService.Delete(context.Background(), int32(resId))
	if err != nil {
		return "", err
	}
	return "deleted", nil
}

func (cp *commandParser) handleSave(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("arg '[id]' is empty, type 'help' to display available commands format")
	}
	var resource any
	var meta string
	resType := args[0]
	switch resType {
	case consts.LoginPasswordArg:
		resource, meta = cp.readLoginPassword()
		return cp.saveTextResource(resource, meta, enum.LoginPassword)
	case consts.BankCardArg:
		resource, meta = cp.readBankCard()
		return cp.saveTextResource(resource, meta, enum.BankCard)
	case consts.FileArg:
		return cp.saveFile()
	default:
		return "", fmt.Errorf("resource type argument '%s' is not supported, type 'help' to display available types", resType)
	}
}

func (cp *commandParser) saveTextResource(resource any, meta string, resType enum.ResourceType) (string, error) {
	resourceJson, err := json.Marshal(resource)
	if err != nil {
		return "", err
	}
	id, err := cp.resourceService.Save(context.Background(), resType, resourceJson, []byte(meta))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("saved successfully, id: %v", id), nil
}

func (cp *commandParser) saveFile() (string, error) {
	filePath := cp.readString("input file path")
	meta := cp.readString("input description")
	id, err := cp.resourceService.SaveFile(context.Background(), filePath, []byte(meta))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", id), nil
}

func (cp *commandParser) readLoginPassword() (*resources.LoginPassword, string) {
	login := cp.readString("input login")
	password := cp.readPassword()
	description := cp.readString("input description")

	return resources.NewLoginPassword(login, password), description
}

func (cp *commandParser) readBankCard() (*resources.BankCard, string) {
	number := cp.readString("input number")
	expireAt := cp.readString("input expireAt in format: MM/YY")
	name := cp.readString("input name")
	surname := cp.readString("input surname")
	description := cp.readString("input description")

	return resources.NewBankCard(number, expireAt, name, surname), description
}

func (cp *commandParser) readString(label string) string {
	if len(label) != 0 {
		fmt.Println(label)
	}
	fmt.Print("-> ")
	cp.scanner.Scan()
	return cp.scanner.Text()
}

func (cp *commandParser) readPassword() string {
	fmt.Println("password:")
	fmt.Print("-> ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	fmt.Println()
	return string(bytePassword)
}

func (cp *commandParser) initScanner() *bufio.Scanner {
	buf := make([]byte, maxCapacity)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(buf, maxCapacity)
	return scanner
}
