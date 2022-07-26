# Mist CLI

Mist CLI is a command line tool for managing multicloud infrastructure. It closely resembles [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) both in terms of functionality and ease of use.

The CLI requires a connection to at least one instance of the Mist Cloud Management Platform. The available options include [Mist CE](https://github.com/mistio/mist-ce), [Mist EE](https://github.com/mistio/mist-ee) and [Mist HS](https://mist.io). You can connect it to multiple ones.

## Quickstart

First, install the CLI following the instructions [here](https://github.com/mistio/mist-cli/releases). Mist CLI supports Windows, Linux and MacOS.

If you don't have an account in a Mist instance already, the fastest way to get started is to sign up for a [free trial of Mist HS](https://mist.io/signup).

Sign in Mist and generate an API key from your account section, e.g. https://mist.io/my-account/tokens. Copy the API token generated and, in your local machine, create a new context with the following command:

```
mist config add-context <name> <api-key>
```

If you don't use Mist HS, you should add the `--server` flag with the URL of your installation, e.g.

```
mist config add-context --server <URL> <name> <api-key>
```

All your configuration settings are saved in the `credentials.json` file in the `$HOME/.mist` directory.

You are now ready to manage your clouds from the command line!


## Syntax

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

## Examples

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

You can output data in JSON, YAML and CSV format by using the `-o <format>` flag. The supported `format` options are `json`, `csv`, and `yaml`.

Here is an example with YAML:

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

### Listings with searching by attributes

```
$ mist get machines --search "state:running AND cloud:Linode"
NAME           	CLOUD 	STATE  	TAGS
InfluxDB1      	Linode	running	staging
LAMP           	Linode	running	    	
debian-ap-south	Linode	running
```

`"state:running AND cloud:Linode"` is also equivalent to `"state:running cloud:Linode"`.

### Listings with JMESPath query manipulation

Get the total number of your clouds:

```
$ mist get clouds -q meta.total
29
```

List the names and private IPs of all the machines:

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

`-q "data[:].[name, private_ips]"` is also equivalent to `-q data[:].[name,private_ips]`. The double quotes help you escape white space on the queries. For more information about JMESPath check this [tutorial](https://jmespath.org/tutorial.html).

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

If a machine is associated with an SSH key in Mist, you can connect to it without access to the private key.

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

You can use `CTRL + D` or type `logout` to the remote terminal to exit.

Please note, that the public key needs to be in the user's `~/.ssh/authorized_keys` file in the target machine. This is done automatically when you create a machine through Mist.

### Kubeconfig

With the `mist kubeconfig` command you can get auto-renewing kubeconfig credentials for kubectl.

```
~ $ mist get cluster
NAME    	CLOUD                           	TOTAL NODES	TAGS 
dijkstra	64e0e46aac45456eab6825e0f9757007	8          	    	
turing  	64e0e46aac45456eab6825e0f9757007	6          	    	

~ $ mist kubeconfig update turing --yes
Clusters "turing" added to the local kubeconfig

~ $ kubectl cluster-info
Kubernetes control plane is running at https://23.115.105.54:443
GLBCDefaultBackend is running at https://23.115.105.54:443/api/v1/namespaces/kube-system/services/default-http-backend:http/proxy
KubeDNS is running at https://23.115.105.54:443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
Metrics-server is running at https://23.115.105.54:443/api/v1/namespaces/kube-system/services/https:metrics-server:/proxy
```