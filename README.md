# gitf
A 'Git-like' interface for (also) pushing your Git repository (or any vanilla folder) via FTP.
Written in Go. A work-in-progress.

## Install
You need a Go environment and to download and install the source.
 - go get github.com/olliephillips/gitf
 - go install github.com/olliephillips/gitf

## Documentation
**init**: initialises respository/directory, creates gitf.toml and gitf.log. Adds to .gitignore

-s server -u username -p password -P port


**push**: sends files in local directory to FTP server configured in gitf.toml

-s server -u username -p password -P port


**pull** (not implemented) : retrieves files to local directory from FTP server configured in gitf.toml

-s server -u username -p password -P port


**status**: reports last gitf operation from gitf.log


**log**: reports all gitf operations from gitf.log


**help**: lists available commands

## Roadmap	
gitf push
	- error logging??, channels to do quicker??, single file upload??
	
gitf pull
	- all of it