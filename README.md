# Mist CLI

## [Releases and installation instructions](https://github.com/mistio/mist-cli/releases)

## Overview
The `mist` command line tool lets you control mist instances. For configuration, `mist` CLI looks for a file named `credentials.json` in the $HOME/.mist directory.

### Syntax
Use the following syntax to run `mist` CLI commands from your terminal window:
```
mist [command] [TYPE] [NAME] [flags]
```
where `command`, `TYPE`, `NAME`, and `flags` are:
* `command`: Specifies the operation that you want to perform on one or more resources, for example `get`.

* TYPE: Specifies the resource type. Resource types are case-sensitive and you can specify the singular, plural, or abbreviated forms. For example, the following commands produce the same output:
    ```
    mist get machine machine1
    mist get machines machine1
    mist get ma machine1
    ```
* `NAME`: Specifies the name of the resource. Names are case-sensitive. If the name is omitted, details for all resources are displayed, for example `mist get machines`.

* flags: Specifies optional flags. For example, you can use the `--server` flag to specify the address of the mist installation.

If you need help, just run `mist help` from the terminal window.

### First time usage

In order to use the mist CLI to interact with a mist installation, you should generate first a mist API key. After that you need to create a new context by applying the following command.
```
mist config add-context [flags] <name> <api-key>
```
Note: In case you don't use mist.io's Hosted Service, you should add the `--server` flag in the above command with the URL of your installation

Now you should be able to manage your clouds through mist CLI!
### Common usage
#### Listings
```
$ mist get machines
NAME                                            CLOUD                   STATE           TAGS               
ec2-f                                           EC2 Frankfurt           running                           
kvm-instace                                     KVM                     running                           
esxi-v7.0                                       Equinix Metal           running                           
DockerHost - VSphere                            DockerHost - VSphere    unknown                           
gke-machine                                     GCE mist                running                           
xtest-DO-w-volume                               Digital Ocean 3         running                           
test                                            KVM                     terminated
```
#### SSH
```
$ mist ssh machine-name

The programs included with the Debian GNU/Linux system are free software;
the exact distribution terms for each program are described in the
individual files in /usr/share/doc/*/copyright.

Debian GNU/Linux comes with ABSOLUTELY NO WARRANTY, to the extent
permitted by applicable law.
Last login: Tue Mar 23 18:57:52 2021 from mist-ce_huproxy_1.mist-ce_default
root@3a088b51795b:~#
```
Tip: Use `CTRL + D` or type `logout` to the remote terminal to exit.
#### Filter Data for reporting
```
$ mist get machine --only name, cost -o csv

```
## Installation

### Instructions

#### Linux
1. Download the latest release with the command:
    ```
    curl -LO "https://dl.mist.io/cli/<version>/bin/linux/amd64/mist"
    ```
2. Validate the binary (optional):\
    Download the mist checksum file:
    ```
    curl -LO "https://dl.mist.io/cli/<version>/bin/linux/amd64/mist.sha256"
    ```
    Validate the mist binary against the checksum file:
    ```
    echo "$(<mist.sha256) mist" | sha256sum --check
    ```
    If valid, the output is:
    ```
    mist: OK
    ```
    If the check fails, sha256 exits with nonzero status and prints output similar to:
    ```
    mist: FAILED
    sha256sum: WARNING: 1 computed checksum did NOT match
    ```
3. Make mist binary executable:
    ```
    chmod +x ./mist
    ```
4. Install mist on path (optional):\
    e.g.
    ```
    mv mist /usr/local/bin/mist
    ```
5. Enable mist autocompletion on all your sessions (optional):\
    Bash
    ```
    echo 'source <(mist completion bash)' >>~/.bashrc
    ```
    Zsh
    ```
    echo 'source <(mist completion zsh)' >>~/.zshrc
    ```
#### MacOS
1. Download the latest release with the command:
    ```
    curl -LO "https://dl.mist.io/cli/<version>/bin/darwin/amd64/mist"
    ```
2. Validate the binary (optional):\
    Download the mist checksum file:
    ```
    curl -LO "https://dl.mist.io/cli/<version>/bin/darwin/amd64/mist.sha256"
    ```
    Validate the mist binary against the checksum file:
    ```
    echo "$(<mist.sha256) mist" | sha256sum --check
    ```
    If valid, the output is:
    ```
    mist: OK
    ```
    If the check fails, sha256 exits with nonzero status and prints output similar to:
    ```
    mist: FAILED
    sha256sum: WARNING: 1 computed checksum did NOT match
    ```
3. Make mist binary executable:
    ```
    chmod +x ./mist
    ```
3. Install mist on path (optional):\
    e.g.
    ```
    sudo mv ./mist /usr/local/bin/mist
    sudo chown root: /usr/local/bin/mist
    ```
4. Enable mist autocompletion on all your sessions (optional):\
    Bash
    - Check bash version
        ```
        echo $BASH_VERSION
        ```
    - Install/upgrade bash to v4.0+ if older
        ```
        brew install bash
        ```
    - Source completion script
        ```
        echo 'source <(mist completion bash)' >>~/.bash_profile
        ```
    Zsh
    ```
    echo 'source <(mist completion zsh)' >>~/.zshrc
    ```
#### Windows
1. Download the latest release with the command:
    ```
    curl -LO "https://dl.mist.io/cli/<version>/bin/windows/amd64/mist.exe"
    ```
2. Validate the binary (optional):\
    Download the mist checksum file:
    ```
    curl -LO "https://dl.mist.io/cli/<version>/bin/windows/amd64/mist.exe.sha256"
    ```
    Validate the mist binary against the checksum file:
    - Using Command Prompt to manually compare CertUtil's output to the checksum file downloaded:
        ```
        CertUtil -hashfile mist.exe SHA256
        type mist.exe.sha256
        ```
    - Using PowerShell to automate the verification using the -eq operator to get a True or False result:
        ```
        $($(CertUtil -hashfile .\mist.exe SHA256)[1] -replace " ", "") -eq $(type .\mist.exe.sha256)
        ```
3. Add the binary in to your `PATH`.
4. Install mist on path (optional):\
    e.g.
    ```
    mv mist /usr/local/bin/mist
    ```
5. Enable mist autocompletion on all your sessions (optional):\
    Powershell
    ```
    mist completion powershell
    ```
    Bash
    ```
    echo 'source <(mist completion bash)' >>~/.bashrc
    ```
    Zsh
    ```
    echo 'source <(mist completion zsh)' >>~/.zshrc
    ```
