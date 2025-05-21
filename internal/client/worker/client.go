package worker

import (
	"fmt"
	"strings"

	"github.com/Alexandrfield/Smaug/internal/common"
	"github.com/Alexandrfield/Smaug/internal/common/security"
)

type TunnelToServer interface {
	Registration(login string, password string) error
	Login(login string, password string) ([]byte, error)
	SaveData(dat security.DataForSave) error
	GetData(name string) (security.DataForSave, error)
	GetAllData() ([]security.DataForSave, error)
}
type clientCLI struct {
	logger        common.Logger
	cred          *security.Credential
	networkClient TunnelToServer
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
	cli.cred = security.GetNewCredential(passwordinp, logininp, temp, cli.logger)
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
	cli.cred = security.GetNewCredential(passwordinp, logininp, temp, cli.logger)
	return nil
}
func SaveData() {

}
func GetData() {

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
