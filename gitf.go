package main

/*
push - error logging, channels to do quicker, single file upload
pull
*/

import (
	"bytes"
	"log"
	"os"
	"flag"
	"fmt"
	"sync"
	"time"
	"strings"
	"io/ioutil"
	"path/filepath"
	"github.com/jlaffaye/ftp"	
	"github.com/BurntSushi/toml"
)

var configFile = "gitf.toml"
var logFile = "gitf.log"
var writeConf = gitfConfig {
	FTP: Config {},
}

type gitfConfig struct {
    FTP Config
}

type Config struct {
    Server string `toml:"server"`
    Port int `toml:"port"`
	User string `toml:"user"`
	Pwd string `toml:"pwd"`
	RemoteDir string `toml:"remote_dir"`
    MaxConnections int `toml:"max_connections"`
}	

// Writes the TOML config file, using flags if passed, or defaults
func writeConfig(wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := os.Stat(configFile)
	if err != nil {
		f, err := os.Create(configFile)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		
		var firstBuffer bytes.Buffer
		e:= toml.NewEncoder(&firstBuffer)
		err = e.Encode(writeConf)
		if err != nil {
			log.Fatal(err)
		}

		f.WriteString(firstBuffer.String())	
		return
	}
}

// Reads the TOML config file, or displays error message if not exist
func readConfig()(gitfConfig) {
	var config gitfConfig

	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatal(configFile, " could not be found, please use gitf init")
	}
	
	_, err = toml.DecodeFile(configFile, &config)
	if err != nil {
		log.Fatal(err)
	}	
	return config	
}

// Creates the log file if it does not exist
func createLog(wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := os.Stat(logFile)
	if err != nil {
		f, err := os.Create(logFile)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	return	
}

func addLog(msg string, status string) {
	_, err := os.Stat(logFile)
	if err == nil {
		tS := time.Now()
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND, 0660)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		log := fmt.Sprintf("%s (%s) : %s\n", msg, status, tS)
		f.WriteString(log)
	}
	return
}	

func addGitignore(wg *sync.WaitGroup) {
	defer wg.Done()
	ignoreFile := ".gitignore"
	_, err := os.Stat(ignoreFile)
	if err != nil {
		// Does not exist create
		f, err := os.Create(ignoreFile)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		f.WriteString("#Ignore gitf files\n")
		f.WriteString("gitf.toml\n")
		f.WriteString("gitf.log\n")
	} else {
		// Exists, append
		f, err := os.OpenFile(ignoreFile, os.O_RDWR|os.O_APPEND, 0660)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		f.WriteString("#Ignore gitf files\n")
		f.WriteString("gitf.toml\n")
		f.WriteString("gitf.log\n")
	}	
	return	
}	

func initCommand() {
	/*
	gitf init [opts -v versioncontrol for gitf files -a activemode ]
	*/
	
	// Parse flags
	initFlag := flag.NewFlagSet("", flag.ExitOnError)
	ServerArg := initFlag.String("s", "localhost", "IP address or name of FTP server")
	PortArg:= initFlag.Int("P", 21, "Port to connect to FTP server on")
	UserArg:= initFlag.String("u", "Username", "Username for login")
	PwdArg:= initFlag.String("p", "Password", "Password for login")
	RemoteDirArg := initFlag.String("d", "/", "Remote directory to upload files to")
	MaxConnectionsArg := initFlag.Int("c", 3, "Maximum concurrent connections to be opened")
	VersionControlArg := initFlag.String("v", "false", "Should gitf.toml and gitf.log file be commited to Git")
	
	initFlag.Parse(os.Args[2:])
	
	// Pass to configuration structs
	writeConf.FTP.Server = *ServerArg
	writeConf.FTP.Port = *PortArg
	writeConf.FTP.User = *UserArg
	writeConf.FTP.Pwd = *PwdArg
	writeConf.FTP.RemoteDir = *RemoteDirArg
	writeConf.FTP.MaxConnections = *MaxConnectionsArg
	versionControl := *VersionControlArg		
	
	// Test for all files
	_, err1 := os.Stat(configFile)
	_, err2 := os.Stat(logFile)	
	
	if err1 != nil && err2 != nil {
		
		var wg sync.WaitGroup
		if versionControl == "false" {
			wg.Add(3)
			go addGitignore(&wg)
		} else {
			wg.Add(2)
		}	

		go writeConfig(&wg)
		go createLog(&wg)

		wg.Wait()
		
		// Test for all files
		_, err3 := os.Stat(configFile)
		_, err4 := os.Stat(logFile)	

		if versionControl == "false" {
			_, err5 := os.Stat(".gitignore")
			if err3 == nil && err4 == nil && err5 == nil {
				addLog("Init", "OK")
			} else {
				addLog("Init", "FAIL - could not write log and/or config files")
			}
		} else {
			if err3 == nil && err4 == nil {
				addLog("Init", "OK")
			} else {
				addLog("Init", "FAIL - could not write log and/or config files")
			}
		}
	} else {
		// Already initialised, print message
		fmt.Println("!gitf already initialised for this repository/directory")
		addLog("Init", "FAIL - already initialised")
	}	
}	

func pushCommand() {
	var ignoreItems = []string{configFile, logFile, ".gitignore", ".git", ".DS_Store"}

	// Scan directory and get list for files and subdirecoties
	curDir, _ := os.Getwd()
	fileArray := []string{}
	dirArray := []string{}	
	
	err := filepath.Walk(curDir, func(path string, f os.FileInfo, _ error) error {
		if (path != curDir) {
			trimPath := strings.Replace(path, curDir, "", 1)
			trimPath = strings.Replace(trimPath, "/", "", 1)
			walkIgnore := false
			for _ , ignore := range ignoreItems {
				if ignore == f.Name() {
					walkIgnore = true
				}	
			}
			
			if walkIgnore == false {
				if f.IsDir() {
					// Directory
					dirArray = append(dirArray, trimPath)
				} else {
					// File
					fileArray = append(fileArray, trimPath)
				}
				
				/*	
				if strings.Index(trimPath, ".") != -1 {
					// File
					fileArray = append(fileArray, trimPath)
				} else {
					// Directory
					dirArray = append(dirArray, trimPath)
				}
				*/
			}		
		}
		return nil
	})	
	
	if err != nil {
		
	}	

	// FTP testing 
	// Connect
	config := readConfig()
	server := fmt.Sprintf("%s:%d", config.FTP.Server, config.FTP.Port)
	ftp, err := ftp.Connect(server)
	if err != nil {
		log.Fatal(err)
	}
	
	// Login
	err = ftp.Login(config.FTP.User, config.FTP.Pwd)
	if err != nil {
		log.Fatal(err)
	}

	// Check directories make if not exist
	var fullPath string
	for _, dir := range dirArray {
		fullPath = config.FTP.RemoteDir + dir
		err = ftp.ChangeDir(fullPath)
		if err != nil {
			err = ftp.MakeDir(fullPath)
			if err != nil {
				log.Fatal(err)
			}	
		}	
	}
	
	// Upload the files - Channels
	for _, file := range fileArray {
		f, err := os.Open(file)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		fullPath = config.FTP.RemoteDir + file
		err = ftp.Stor(fullPath, f)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Push (OK) : " + fullPath)
		}		
	}	

	ftp.Logout()
	ftp.Quit()
	
	addLog("Push", "OK")	
}

func pullCommand() {
	//var config = ReadConfig()
	addLog("Pull", "OK")
}

func statusCommand() {	
	logs, err := ioutil.ReadFile(logFile)
	if err != nil {
	    log.Fatal(logFile, " could not be found, please use gitf init")
	}
	logLines := strings.Split(string(logs), "\n")
	logLineCount := len(logLines)
	fmt.Println("# gitf status")
	fmt.Println(logLines[logLineCount - 2])
}

func logCommand() {
	// Read log file in display to std.out
	logs, err := ioutil.ReadFile(logFile)
	if err != nil {
	    log.Fatal(logFile, " could not be found, please use gitf init")
	}
	logLines := strings.Split(string(logs), "\n")
	logLineCount := len(logLines)
	fmt.Println("# gitf log")
	for key, value := range logLines {
		// Don't want the blank line
		if key !=  logLineCount - 1 { 
	    	fmt.Println(value)
		}
	}
}	

func help(){
	// Some help on the gitf commands and arguments
	fmt.Println("#gitf help: commands and arguments (optional)")
	fmt.Println(" init: initialises respository/directory, creates gitf.toml and gitf.log. Adds to .gitignore")
	fmt.Println("  -s server -u username -p password -P port")
	fmt.Println(" push: sends files in local directory to FTP server configured in gitf.toml")
	fmt.Println("  -s server -u username -p password -P port")
	fmt.Println(" pull: retrieves files to local directory from FTP server configuted in gitf.toml")
	fmt.Println("  -s server -u username -p password -P port")
	fmt.Println(" status: reports last gitf operation from gitf.log")
	fmt.Println(" log: reports all gitf operations from gitf.log")
}
					
func main() {
	if len(os.Args) > 1 {
		command := os.Args[1];
		switch command {
			case "init":
				initCommand()
			case "push":
				pushCommand()
			case "pull":
				pullCommand()
			case "status":
				statusCommand()
			case "log":
				logCommand()
			case "help":
				help()		
			default:
				help()
		}
	} else {
		help()
	}		
}	