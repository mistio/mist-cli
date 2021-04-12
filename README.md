# Mist CLI

## [Releases and installation instructions](https://github.com/mistio/mist-cli/releases)

## Overview
Wouldn't be nice if you could manage all your clouds and multi-cloud workloads from a single CLI which closely resembles [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) both in functionality and ease of use? This is exactly what Mist CLI allows you to do. Sign up with [mist](https://mist.io/), the opensource multi-cloud management platform, or take a look on mist's other offerings and with this companion CLI tool, you can manage all your clouds and instances efficiently and do even more, within a few minutes! Take a look below to see all the functionality you can gain from the `mist` + `mist CLI` combo.

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

Note: All your configuration settings are saved in the `credentials.json` file in the $HOME/.mist directory.
## Common usage
### Listings
```
$ mist get machines
NAME                                            CLOUD                   STATE           TAGS               
ec2-f                                           EC2 Frankfurt           running         prod,mongodb                  
kvm-instace                                     KVM                     running                           
esxi-v7.0                                       Equinix Metal           running                           
DockerHost - VSphere                            DockerHost - VSphere    unknown         staging                  
gke-machine                                     GCE mist                running                           
xtest-DO-w-volume                               Digital Ocean 3         running                           
test                                            KVM                     terminated
```
### Listings with specific columns
```
 mist get clouds --only name,provider
NAME                     	PROVIDER       	
Aliyun ECS Silicon Valley	aliyun_ecs  	
Azure ARM New            	azure_arm   	
DigitalOcean             	digitalocean	
DockerHost - VSphere     	docker      	
EC2 Frankfurt            	ec2         	         	
EC2 N. California        	ec2         	        	
Equinix Metal            	equinixmetal	
G8                       	gig_g8      	
GCE staging          	    gce         	        	
KVM                      	libvirt     	
KubeVirt                 	kubevirt    	
LXD                      	lxd         	
Linode                   	linode      	      	
Maxihost 2               	maxihost    	
OnApp                    	onapp       	
OpenStack train test     	openstack   	
Rackspace Dallas         	rackspace   	  	
SoftLayer                	softlayer   	
Vultr New                	vultr       	
vSphere 7 on Metal       	vsphere     	 
```
### Listings in different output formats
```
$ mist get machines InfluxDB1 -o yaml
data:
  cloud: Linode
  cores: 1
  cost:
    hourly: 0.013888888888888888
    monthly: 10
  created: "2020-12-11T12:17:59Z"
  created_by: ""
  expiration: null
  ...
  hostname: 1.2.3.4
  ...
  location: Frankfurt, DE
  ...
  state: running
  ...
meta:
  returned: 1
  sort: ""
  start: 0
  total: 2
```
Note: You can also output data in json and csv form by using the `-o <format>` flag where format can be `json`, `csv`, or `yaml`
### Listings with searching by attributes
```
$ mist get machines --search "state:running AND cloud:Linode"
NAME           	CLOUD 	STATE  	TAGS 
InfluxDB1      	Linode	running	staging	
LAMP           	Linode	running	    	
debian-ap-south	Linode	running	
```
Note: instead of: `"state:running AND cloud:Linode"` we could also have written: `"state:running cloud:Linode"`
### Listings with JMESPath query manipulation
Get the total number of your clouds
```
$ mist get clouds -q meta.total
29
```
List the names and private IPs of all the machines
```
$ mist get machines -q "data[:].[name, private_ips]"
[
  [
    "InfluxDB1",
    [
      "192.168.169.117"
    ]
  ],
  [
    "LAMP",
    [
      "192.168.159.9"
    ]
  ],
  [
    "debian-ap-south",
    []
  ]
]
```
Note: `-q "data[:].[name, private_ips]"` is equivalent to this `-q data[:].[name,private_ips]`. The double quotes help you to escape white space on the queries. For more information about JMESPath check this [tutorial](https://jmespath.org/tutorial.html).
### Filter Data for reporting
```
$ mist get machine --only name,cost -o csv
/name,/cost/hourly,/cost/monthly
InfluxDB1,0.013888888888888888,10
debian-ap-south,0.027777777777777776,20
LAMP,0.006944444444444444,5
```
### Find your public key
```
$ mist get key staging -q data.public
ssh-rsa XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX...
```
### SSH
By associating a key to a machine with mist, you can ssh into it without having to fiddle around with ssh keys. When deploying instances with mist you can always have an associated key 
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
Tip: Use `CTRL + D` or type `logout` to the remote terminal to exit.\
Note: The public key needs to be in the user's `~/.ssh/authorized_keys` file in the target machine.\
This can easily be achieved with machines deployed through mist by assigning them a key during creation.
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
    sudo mv mist /usr/local/bin/mist
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
