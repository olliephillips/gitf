# gitf
A 'Git-like' interface for (also) pushing your Git repository (or any vanilla folder) via FTP.
Written in Go. A work-in-progress.

## Install
You need a Go environment and to download and install the source.
 - go get github.com/olliephillips/gitf
 - go install github.com/olliephillips/gitf

## Documentation
**gitf init**: initialises respository/directory, creates gitf.toml and gitf.log. Adds to .gitignore

Optional flags allow parameters to be written directly to gitf.toml. Without these, defaults will be configured which can be amended by editing the file in a text editor.

-s server -u username -p password -P port -d remote directory -v true submit gitf.toml and log files to version control


**gitf push**: sends files in local directory to FTP server configured in gitf.toml


**gitf pull** (not implemented) : retrieves files to local directory from FTP server configured in gitf.toml


**gitf status**: reports last gitf operation from gitf.log


**gitf log**: reports all gitf operations from gitf.log


**gitf help**: lists available commands

## Roadmap	
gitf push
	- error logging??, channels to do quicker??, single file upload??, single file??, compile and upload the binary??
	
gitf pull
	- all of it