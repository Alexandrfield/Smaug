package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Alexandrfield/Smaug/internal/common"
)

type clientSaveData interface {
	Registration(login string, password []byte) error
	Login(login string, password []byte) error
	SaveText(description string, data []byte) error
	SaveFile(description string, filename string) error
	GetData(description string) []*OpenData
	GetAllData() []*OpenData
}

type clientCLI struct {
	client clientSaveData
}

func (cli *clientCLI) Registration() error {
	var login string
	fmt.Printf("         Registration. \n")
	fmt.Printf("input your login:\n > ")
	fmt.Scan(&login)
	var password string
	fmt.Printf("input your password:\n > ")
	fmt.Scan(&password)
	err := cli.client.Registration(login, []byte(password))
	if err != nil {
		fmt.Printf("error Registration. err:%s\n", err)
	}
	return nil
}
func (cli *clientCLI) Login() error {
	fmt.Printf("         Login. \n")
	var login string
	fmt.Printf("input your login:\n > ")
	fmt.Scan(&login)
	var password string
	fmt.Printf("input your password:\n > ")
	fmt.Scan(&password)
	err := cli.client.Login(login, []byte(password))
	if err != nil {
		fmt.Printf("error Login. err:%s\n", err)
	}
	return nil
}

func (cli *clientCLI) SaveText() {
	fmt.Printf("enter a name that you can use to get this information later:\n > ")
	description, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Printf("enter the information to save:\n > ")
	data, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	err := cli.client.SaveText(description, []byte(data))
	if err != nil {
		fmt.Printf("save data error. %s \n", err)
	}
}
func (cli *clientCLI) SaveFile() {
	fmt.Printf("enter a name that you can use to get this information later:\n > ")
	description, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Printf("enter filename to save:\n > ")
	filename, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	filename = strings.TrimSpace(filename)
	err := cli.client.SaveFile(description, filename)
	if err != nil {
		fmt.Printf("save data error. %s \n", err)
	}
}

func (cli *clientCLI) GetData() {
	fmt.Printf("enter a name information:\n")
	description, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	openData := cli.client.GetData(description)
	for i, v := range openData {
		fmt.Printf("%d) %s", i, v.Show())
		fmt.Printf(" -----  *** ----- \n")
	}
}
func (cli *clientCLI) GetAllData() {
	openData := cli.client.GetAllData()
	for i, v := range openData {
		fmt.Printf("%d) %s", i, v.Show())
		fmt.Printf("\n -----  *** ----- \n")
	}
}

func printCommandInfo() {
	fmt.Printf("commands:\n")
	fmt.Printf("           list - get all saved info\n")
	fmt.Printf("           get - gey info\n")
	fmt.Printf("           saveText - add new info\n")
	fmt.Printf("           saveFile - add new info\n")
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
			fmt.Printf("Login has been failed!")
			return
		}
	} else {
		err := cli.Registration()
		if err != nil {
			fmt.Printf("Registration has been failed!")
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
			fmt.Printf("stop program")
			close(done)
			return
		case "list":
			cli.GetAllData()
		case "get":
			cli.GetData()
		case "saveText":
			cli.SaveText()
		case "saveFile":
			cli.SaveFile()
		default:
			printCommandInfo()
		}
	}
}

func ClientAppLoop(done chan struct{}, logger common.Logger, config *Config) {
	client := NewSmaugClient(logger, config)
	cli := clientCLI{client: client}
	cli.ClientAppLoop(done)
}
