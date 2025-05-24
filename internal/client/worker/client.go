package worker

import (
	"fmt"
	"strings"
	"time"

	clientCommon "github.com/Alexandrfield/Smaug/internal/client/common"
	clientSecurity "github.com/Alexandrfield/Smaug/internal/client/security"
	"github.com/Alexandrfield/Smaug/internal/common"
	"github.com/Alexandrfield/Smaug/internal/common/security"
)

type TunnelToServer interface {
	Registration(login string, password string) error
	Login(login string, password string) ([]byte, error)
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
	err := cli.networkClient.Registration(logininp, passwordinp)
	if err != nil {
		return fmt.Errorf("Problem with registration. can't continue. err:%w", err)
	}
	temp, err := cli.networkClient.Login(logininp, passwordinp)
	if err != nil {
		return fmt.Errorf("Problem with login after registration. err:%w", err)
	}
	cred := security.GetNewCredential(passwordinp, logininp, temp, cli.logger)
	cli.secyrityManager = clientSecurity.NewNoteManager(cli.logger, cred)
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
	temp, err := cli.networkClient.Login(logininp, passwordinp)
	if err != nil {
		return fmt.Errorf("Problem with login. err:%w", err)
	}
	cred := security.GetNewCredential(passwordinp, logininp, temp, cli.logger)
	cli.secyrityManager = clientSecurity.NewNoteManager(cli.logger, cred)
	return nil
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
func GetAllData() {

}
func (cli *clientCLI) ClientAppLoop(done chan struct{}) {
	fmt.Printf("-------------------------------\n")
	fmt.Printf("------ Smaug application ------\n")
	fmt.Printf("-------------------------------\n\n\n")

	var command string
	for {
		fmt.Printf("input your command:\n")
		fmt.Scan(&command)
		command = strings.TrimSpace(command)
		if command == "stop" || command == "exit" {
			cli.logger.Infof("stop client loop")
			close(done)
			return
		}

	}
}

func ClientAppLoop(done chan struct{}, logger common.Logger) {
	cli := clientCLI{logger: logger}
	cli.ClientAppLoop(done)
}
