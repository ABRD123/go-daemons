- [Setup](#setup)
  * [Live/Beta Environment](#livebeta-environment)
    + [Daemon/Server Setup Part I](#daemonserver-setup-part-i)
    + [Daemon/Server Setup Part II](#daemonserver-setup-part-ii)
  * [Development Environment](#development-environment)
    + [OSX Go Environment](#osx-go-environment)
    + [OSX godaemons Runtime Setup](#osx-godaemons-runtime-setup)
    + [OSX godaemons Development Setup](#osx-godaemons-development-setup)
  * [Directory Structure](#directory-structure)

# Setup

## Live/Beta Environment

### Daemon/Server Setup Part I

**NOTE** These steps should only be performed on the daemon/server systems. 

These setps require you to be logged in as the root user. 
   - sudo bash

1. Make sure the Ubuntu system is up-to-date 
   - apt update
   - apt upgrade
1. Install golang 1.11 
   - Production:
     - Ref: [Linux, macOS, and FreeBSD tarballs](https://golang.org/doc/install#install)
     - Download go1.11.4.linux-amd64.tar.gz to your local machine from <https://golang.org/dl/>
     - SCP the file to the prod box through bastion, login to the box, cd to where the files are, and run the below commands.
     - sudo bash
     - tar -C /usr/local -xzf go1.11.4.linux-amd64.tar.gz
   - Beta:
     - add-apt-repository ppa:longsleep/golang-backports
     - apt-get update
     - apt-get install golang-go
1. Update the global environment for golang 
   - cd /etc/profile.d
   - We need HISTCONTROL=ignoreboth set in the environment so commands that start with a space will not be logged in bash history.
     - echo $HISTCONTROL
     - If it is not set then add it in golang.sh as below.
   - vim golang.sh
     - Production:
       - HISTCONTROL=ignoreboth
       - export PATH=$PATH:/usr/local/go/bin
       - export GOROOT=/usr/local/go
     - Beta:
       - HISTCONTROL=ignoreboth
       - export GOROOT=/usr/lib/go
   - relog in
     - Run "go version" to make sure go is really in the path
1. Create a conf file to re-create the /var/run/godaemons directory after a restart.  
   - vim /usr/lib/tmpfiles.d/godaemons.conf
     ```
     d /var/run/godaemons 0775 godaemons godaemons
     ```
1. Create a godaemons user 
   - addgroup --system --gid 876 godaemons
   - adduser --disabled-password --uid 876 --ingroup godaemons godaemons
   - mkdir -p /var/log/godaemons
   - chmod 755 /var/log/godaemons
   - chown -R godaemons:adm /var/log/godaemons
   - mkdir -p /var/run/godaemons
   - chmod 755 /var/run/godaemons
   - chown -R godaemons:godaemons /var/run/godaemons

References:
   - <https://github.com/golang/go/wiki/Ubuntu>
   - <https://help.ubuntu.com/community/EnvironmentVariables#System-wide_environment_variables>
   - <https://www.digitalocean.com/community/tutorials/how-to-restrict-log-in-capabilities-of-users-on-ubuntu>

### Daemon/Server Setup Part II

**NOTE** These steps should only be performed on the daemon/server systems. 

These instructions require you to be logged in as the godaemons user. 
   - sudo bash
   - su - godaemons

1. Add GODAEMONS run-time environment variables (run as godaemons)
   - cd
   - vim ~/.godaemons.sh
   - Double check that you properly configured both *$GOPATH* and *$GOBIN* and that your *$GOBIN* is in your *$PATH*. 
   NOTE: GOROOT should already exist. 
     ```
     # Golang settings for godaemons
     export GOPATH=$HOME/go # path of the workspace
     export GOBIN=$GOPATH/bin # path of the bin folder in the workspace
     export PATH=$GOBIN:$PATH
     alias godaemons="cd $GOPATH/src/github.com/godaemons"
     alias godaemonslogs="cd /var/log/godaemons/"
     ```
   - chmod 700 ~/.godaemons.sh
   - vim ~/.profile
     - \# Load the environment for the godaemons project.
     - if [ -f '/home/godaemons/.godaemons.sh' ]; then source '/home/godaemons/.godaemons.sh'; fi
   * Exit and relog in - restart the shell
   * Run "set | grep GODAMEONS" to verify the settings
   * Run "set | grep GO" to verify the go settings
   * Run "set | grep PATH" to verify the PATH is correct
   - Create the go folders:
     - mkdir -p $GOPATH/src $GOPATH/pkg $GOPATH/bin
1. Create a ssh key for the new godaemons user (run as godaemons)
   - ssh-keygen -t ed25519 -a 100
   - Add the id_ed25519.pub to the godaemons repo settings in this location: 
     - <https://github.com/go-daemons/settings/keys>
1. Clone the godaemons repo (run as godaemons)
   - mkdir -p go/src/github.com/go-daemons
   - cd go/src/github.com/go-daemons
   - git clone git@github.com:aapi123/godaemons.git
1. Update ~/.ssh/known_hosts (run as godaemons)
   - *NOTE this only needs to be run on live systems.*
   - gored
   - cd scripts
   - ./gen_knownhosts.sh
   - If desired view the contents of ~/.ssh/known_hosts to make sure it looks correct. 
1. Build godaemons (run as godaemons)
   - NOTE: First time only and in Beta Only - install dep:
     - Beta:
       - curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
   - cd ~/go/src/github.com/go-daemons
     - make uninstall && make prod

References:
   - <https://security.stackexchange.com/questions/143442/what-are-ssh-keygen-best-practices>

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
