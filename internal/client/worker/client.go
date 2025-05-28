package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	clientCommon "github.com/Alexandrfield/Smaug/internal/client/common"
	clientSecurity "github.com/Alexandrfield/Smaug/internal/client/security"
	"github.com/Alexandrfield/Smaug/internal/common"
)

type TunnelToServer interface {
	Registration(login string, password string) error
	Login(login string, password []byte) error
	SaveData(login string, dat []byte) error
	GetData(login string, dat []byte) ([][]byte, error)
	GetAllData(login string, dat []byte) ([][]byte, error)
}

type clientCLI struct {
	logger          common.Logger
	secyrityManager *clientSecurity.NoteManager
	networkClient   TunnelToServer
}

func (cli *clientCLI) Registration() error {
	var logininp string
	fmt.Printf("         Registration. \n")
	fmt.Printf("input your login:\n > ")
	fmt.Scan(&logininp)
	var passwordinp string
	fmt.Printf("input your password:\n > ")
	fmt.Scan(&passwordinp)
	passwordForServer := clientSecurity.HashInputPassword([]byte(passwordinp))
	err := cli.networkClient.Registration(logininp, string(passwordForServer))
	if err != nil {
		return fmt.Errorf("problem with registration. can't continue. err:%w", err)
	}
	err = cli.login(logininp, passwordForServer)
	if err != nil {
		return fmt.Errorf("problem with login after registration. err:%w", err)
	}

	return nil
}
func (cli *clientCLI) login(login string, password []byte) error {
	err := cli.networkClient.Login(login, password)
	if err != nil {
		return fmt.Errorf("problem with login. err:%w", err)
	}
	cli.secyrityManager.InitParametrs([]byte(login), []byte(password))
	fmt.Printf("you succes logined.\n")
	return nil
}

func (cli *clientCLI) Login() error {
	fmt.Printf("         Login. \n")
	var logininp string
	fmt.Printf("input your login:\n > ")
	fmt.Scan(&logininp)
	var passwordinp string
	fmt.Printf("input your password:\n > ")
	fmt.Scan(&passwordinp)
	passwordForServer := clientSecurity.HashInputPassword([]byte(passwordinp))
	return cli.login(logininp, passwordForServer)
}

func (cli *clientCLI) SaveData() error {
	fmt.Printf("enter a name that you can use to get this information later:\n > ")
	description, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Printf("enter the information to save:\n > ")
	plainText, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	metadata := clientCommon.NewMetadata("text", cli.secyrityManager.GetLogin(), time.Now())
	newNote := cli.secyrityManager.CreateNewNotes(plainText, description, metadata)
	err := cli.networkClient.SaveData(cli.secyrityManager.GetLogin(), newNote.Serialize())
	if err != nil {
		return fmt.Errorf("SaveData err:%w", err)
	}
	return nil
}

func (cli *clientCLI) GetData() {
	fmt.Printf("enter a name information:\n")
	description, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	tempNote := cli.secyrityManager.CreateServiceNotes(description)
	data, err := cli.networkClient.GetData(cli.secyrityManager.GetLogin(), tempNote.Serialize())
	if err != nil {
		fmt.Printf("Problem Get data. err:%s", err)
	}
	for _, val := range data {
		note := cli.secyrityManager.OpenNote(val)
		plaintext, metadata := cli.secyrityManager.GetInfoFromNote(note)
		if len(plaintext) == 0 || metadata == nil {
			continue
		}
		fmt.Printf("name:%s; data:[%s]; \n      metadata: %s\n\n", note.Info, plaintext, metadata.GetInfo())
	}
}
func (cli *clientCLI) GetAllData() {
	cli.logger.Debugf("try GetAllData")
	tempNote := cli.secyrityManager.CreateServiceNotes("get all data")
	data, err := cli.networkClient.GetAllData(cli.secyrityManager.GetLogin(), tempNote.Serialize())
	if err != nil {
		fmt.Printf("Problem Get all data. err:%s", err)
	}
	for _, val := range data {
		note := cli.secyrityManager.OpenNote(val)
		plaintext, metadata := cli.secyrityManager.GetInfoFromNote(note)
		if len(plaintext) == 0 || metadata == nil {
			continue
		}
		fmt.Printf("name:%s; data:[%s]; \n      metadata: %s\n\n", note.Info, plaintext, metadata.GetInfo())
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
	fmt.Printf("Do you have login? yes/no\n")
	var reg string
	fmt.Scan(&reg)
	reg = strings.TrimSpace(reg)
	if reg == "yes" {
		err := cli.Login()
		if err != nil {
			cli.logger.Warnf(" login(). err:%s", err)
			return
		}
	} else {
		err := cli.Registration()
		if err != nil {
			cli.logger.Warnf(" registration(). err:%s", err)
			close(done)
			return
		}
	}
	printCommandInfo()
	var command string
	for {
		fmt.Printf("\n-------------------------------\n")
		fmt.Printf("input your command:\n > ")
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

func ClientAppLoop(done chan struct{}, logger common.Logger, config *Config) {
	secyrityManager := clientSecurity.NewNoteManager(logger)
	cli := clientCLI{logger: logger, secyrityManager: secyrityManager, networkClient: &GRPCClient{serverAddr: config.ServerAddr}}
	cli.ClientAppLoop(done)
}
