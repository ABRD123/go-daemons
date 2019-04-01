- [Setup](#setup)
  * [Development Environment](#development-environment)
    + [OSX Go Environment](#osx-go-environment)
    + [OSX godaemons Runtime Setup](#osx-godaemons-runtime-setup)
    + [OSX godaemons Development Setup](#osx-godaemons-development-setup)
  * [Directory Structure](#directory-structure)

# Setup

## Development Environment

### OSX Go Environment

- Install Go
  - **NOTE** the instructions in this section assume a Homebrew installation.
  - If you want the Homebrew installation do (for more info on Hombrew go [here](https://brew.sh/))...
    - brew update
    - brew install go
  - If you want a more manual installation follow the [Getting started with Go steps](https://golang.org/doc/install).
    - If you choose this method pay attention to the GO environment variables as the pathes differ from a Homebrew installation.
- Setup the Golang environment variables.
   - cd
   - We need HISTCONTROL=ignoreboth set in the environment so commands that start with a space will not be logged in bash history.
     - echo $HISTCONTROL
     - If it is not set then add HISTCONTROL=ignoreboth to which ever file you use: .profile or .bashrc
   - vim ~/.golang.sh
     ```
     #!/usr/bin/env bash

     export GOROOT=/usr/local/opt/go/libexec
     export GOPATH=/Users/[yourUserID]/go #TODO update with your user ID
     export GOBIN=$GOPATH/bin
     export PATH=$PATH:$GOBIN:$GOROOT/bin
     ```
   - vim ~/.profile
     ```
     # Go Language
     if [ -f '/Users/[yourUserID]/.golang.sh' ]; then source '/Users/[yourUserID]/.golang.sh'; fi #TODO update with your user ID
     ```
   - Exit and relog in - restart the shell
   - Run "set | grep GO" to verify the go settings
   - Run "set | grep PATH" to verify the PATH is correct
- Create the go folders:
  - mkdir -p $GOPATH/src $GOPATH/pkg $GOPATH/bin
  - For more information on workspaces see [Go workspaces](https://golang.org/doc/code.html#Workspaces).

### OSX godaemons Runtime Setup

- Create a bash script that will setup the required godaemons runtime directories.
  - vim ~/bin/setup-godaemons.sh
    ```
    #!/usr/bin/env bash

    sudo mkdir -p /var/log/godaemons
    sudo chmod 755 /var/log/godaemons
    sudo chown [yourUserID]:wheel /var/log/godaemons #TODO update with your user ID
    sudo mkdir -p /var/run/godaemons
    sudo chmod 755 /var/run/godaemons
    sudo chown [yourUserID]:daemon /var/run/godaemons # TODO update with your user ID
    ```
  - Run ~/bin/setup-godaemons.sh to create the necessary run-time directories.
- As an optional step you can create a launchd file that will run setup-godaemons.sh after every reboot. Otherwise, you'll need to manually re-run ~/bin/setup-godaemons.sh after a reboot.
  - sudo vim /Library/LaunchDaemons/com.godaemons.plist
    ```
    <?xml version="1.0" encoding="UTF-8"?>
    <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
    <plist version="1.0">
      <dict>
        <key>EnvironmentVariables</key>
        <dict>
          <key>PATH</key>
          <string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:</string>
        </dict>
        <key>Label</key>
        <string>com.godaemons</string>
        <key>Program</key>
        <string>/Users/[yourUserID]/bin/setup-godaemons.sh</string> #TODO update with your user ID
        <key>RunAtLoad</key>
        <true/>
        <key>KeepAlive</key>
        <false/>
        <key>LaunchOnlyOnce</key>        
        <true/>
        <key>StandardOutPath</key>
        <string>/tmp/godaemons.stdout</string>
        <key>StandardErrorPath</key>
        <string>/tmp/godaemons.stderr</string>
        <key>UserName</key>
        <string>root</string>
      </dict>
    </plist>
    ```
  - sudo launchctl load -w /Library/LaunchDaemons/com.godaemons.plist

References:
- <https://medium.com/@fahimhossain_16989/adding-startup-scripts-to-launch-daemon-on-mac-os-x-sierra-10-12-6-7e0318c74de1>
- <https://docs.chef.io/resource_launchd.html>

### OSX godaemons Development Setup

- Setup the godaemons environment variables.
   - cd
   - vim ~/.godaemons.sh
     ```
     # godaemons aliases
     alias gored="cd $GOPATH/src/github.com/go-daemons"
     alias goredlogs="cd /var/log/godaemons/"
     
     # godaemons run-time environment
     ```
   - chmod 700 ~/.godaemons.sh
   - vim ~/.profile
     ```
     # Load the environment for the godaemons project.
     if [ -f '/Users/[yourUserID]/.godaemons.sh' ]; then source '/Users/[yourUserID]/.godaemons.sh'; fi #TODO update with your user ID
     ```
   - Exit and relog in - restart the shell
   - Run "set | grep GODAEMONS" to verify the settings
- Create a directory structure under the Go src folder that looks like this:
  - cd $GOPATH/src
  - mkdir -p github.com/go-daemons
- Clone the godaemons project from your fork
  - cd github.com/go-daemons
  - git clone git@github.com:[yourUserId]/go-daemons.git #TODO update with your user ID
- Build the godaemons project and populate the DB
  - godameons
  - make install

## Directory Structure

- <https://github.com/golang-standards/project-layout>
