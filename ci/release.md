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
### Binaries

