# gitf #
A 'Git-like' interface for (also) pushing your Git repository (or any vanilla folder) via FTP.
Written in Go. A work-in-progress.

## Install ##
You need a Go environment and to download and install the source.
 - go get github.com/olliephillips/gitf
 - go install github.com/olliephillips/gitf

## Documentation ##
gitf init [opts -m manifest -s server -P port -c connections -u user -p password -d remote directory]
	- creates FTP manifest and logfile with the details provided at commandline optionally, otherwise file contains empty keys
	- adds or creates and adds the FTP manifest to .gitignore by default, -m to not do this
	- gitf.log
	- gitf.toml

gitf push [-b -f filename.txt -v verbose] 
	- Syncs files in directory local files to remote FTP site
	- Option -b to build and upload the binary only	
	
gitf pull 
	- Not implemented

gitf status
	- Returns time date of last FTP operation, it's type, and success or fail
	
gitf log	
	- Returns the entire log history
	
gitf help
	- Lists available commands

## Roadmap ##	
gitf push
	- error logging??, channels to do quicker??, single file upload??
	
gitf pull
	- all of it