package worker

import (
	"fmt"
	"strings"
	"time"

	clientCommon "github.com/Alexandrfield/Smaug/internal/client/common"
	clientSecurity "github.com/Alexandrfield/Smaug/internal/client/security"
	"github.com/Alexandrfield/Smaug/internal/common"
)

type TunnelToServer interface {
	Registration(login string, password string) error
	Login(login string, password string) error
	SaveData(description string, dat []byte) error
	GetData(name string) ([][]byte, error)
	GetAllData() ([][]byte, error)
}

type clientCLI struct {
	logger          common.Logger
	secyrityManager *clientSecurity.NoteManager
	networkClient   TunnelToServer
}

func (cli *clientCLI) Registration() error {
	var logininp string
	fmt.Printf("         Registratuin. \n")
	fmt.Printf("input your login:\n")
	fmt.Scan(&logininp)
	var passwordinp string
	fmt.Printf("input your password:\n")
	fmt.Scan(&passwordinp)
	passwordForServer := clientSecurity.HashInputPassword(passwordinp)
	err := cli.networkClient.Registration(logininp, passwordForServer)
	if err != nil {
		return fmt.Errorf("Problem with registration. can't continue. err:%w", err)
	}
	err = cli.login(logininp, passwordinp)
	if err != nil {
		return fmt.Errorf("Problem with login after registration. err:%w", err)
	}

	return nil
}
func (cli *clientCLI) login(login string, password string) error {
	passwordForServer := clientSecurity.HashInputPassword(password)
	err := cli.networkClient.Login(login, passwordForServer)
	if err != nil {
		return fmt.Errorf("Problem with login. err:%w", err)
	}
	cli.secyrityManager.InitParametrs([]byte(login), []byte(password))
	return nil
}

func (cli *clientCLI) Login() error {
	fmt.Printf("         Login. \n")
	var logininp string
	fmt.Printf("input your login:\n")
	fmt.Scan(&logininp)
	var passwordinp string
	fmt.Printf("input your password:\n")
	fmt.Scan(&passwordinp)
	return cli.login(logininp, passwordinp)
}

func (cli *clientCLI) SaveData() error {
	fmt.Printf("enter a name that you can use to get this information later:\n")
	var description string
	fmt.Scan(&description)
	fmt.Printf("enter the information to save:\n")
	var plainText string
	fmt.Scan(&plainText)
	metadata := clientCommon.NewMetadata("text", cli.secyrityManager.GetLogin(), time.Now())
	newNote := cli.secyrityManager.CreateNewNotes(plainText, description, metadata)
	err := cli.networkClient.SaveData(description, newNote.Serialize())
	if err != nil {
		return fmt.Errorf("SaveData err:%w", err)
	}
	return nil
}

func (cli *clientCLI) GetData() {
	fmt.Printf("enter a name information:\n")
	var description string
	fmt.Scan(&description)

	data, err := cli.networkClient.GetData(description)
	if err != nil {
		fmt.Printf("Problem Get data. err:%s", err)
		return
	}
	for _, val := range data {
		note := cli.secyrityManager.OpenNote(val)
		plaintext, metadata := cli.secyrityManager.GetInfoFromNote(note)
		if len(plaintext) == 0 || metadata == nil {
			continue
		}
		fmt.Printf("data:%s\n", plaintext)
	}
}
func (cli *clientCLI) GetAllData() {
	var description string
	fmt.Scan(&description)

	data, err := cli.networkClient.GetAllData()
	if err != nil {
		fmt.Printf("Problem Get all data. err:%s", err)
		return
	}
	for _, val := range data {
		note := cli.secyrityManager.OpenNote(val)
		plaintext, metadata := cli.secyrityManager.GetInfoFromNote(note)
		if len(plaintext) == 0 || metadata == nil {
			continue
		}
		fmt.Printf("data:%s\n", plaintext)
	}
}

func printCommandInfo() {
	fmt.Printf("commands:\n")
	fmt.Printf("           list - get all saved info\n")
	fmt.Printf("           get - gey info\n")
	fmt.Printf("           save - add new info\n")
	fmt.Printf("           exit - stop programm\n")
}
func (cli *clientCLI) ClientAppLoop(done chan struct{}) {
	fmt.Printf("-------------------------------\n")
	fmt.Printf("------ Smaug application ------\n")
	fmt.Printf("-------------------------------\n\n\n")
	printCommandInfo()
	var command string
	for {
		fmt.Printf("\n-------------------------------\n")
		fmt.Printf("input your command:\n")
		fmt.Scan(&command)
		command = strings.TrimSpace(command)

		switch command {
		case "exit":
			cli.logger.Infof("stop client loop")
			close(done)
			return
		case "list":
			cli.GetAllData()
		case "get":
			cli.GetData()
		case "save":
			cli.SaveData()
		default:
			printCommandInfo()
		}
	}
}

func ClientAppLoop(done chan struct{}, logger common.Logger) {
	secyrityManager := clientSecurity.NewNoteManager(logger)
	cli := clientCLI{logger: logger, secyrityManager: secyrityManager}
	cli.ClientAppLoop(done)
}
