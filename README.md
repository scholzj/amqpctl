[![CircleCI](https://circleci.com/gh/scholzj/amqpctl.svg?style=shield)](https://circleci.com/gh/scholzj/amqpctl)

# amqpctl

`amqpctl` is a command line based client for the AMQP Management protocol. It is written in Go language and tested against [Apache Qpid Dispatch](http://qpid.apache.org/components/dispatch-router/index.html) AMQP router. It should be compatible with all implementations of AMQP Management Working Draft 9 and Working Draft 11.

<!-- TOC -->

- [amqpctl](#amqpctl)
    - [Dependencies](#dependencies)
    - [Installation](#installation)
    - [Configuration](#configuration)
        - [From command line](#from-command-line)
        - [Using config file](#using-config-file)
    - [Usage](#usage)
        - [QUERY](#query)
        - [READ](#read)
        - [UPDATE](#update)
        - [CREATE](#create)
        - [DELETE](#delete)
        - [GET TYPES](#get-types)
        - [GET ATTRIBUTES](#get-attributes)
        - [GET ANNOTATIONS](#get-annotations)
        - [GET MANAGEMENT NODES](#get-management-nodes)
        - [HELP](#help)
    - [TODO](#todo)

<!-- /TOC -->

## Dependencies

`amqpctl` is using [Apache Qpid Proton](http://qpid.apache.org/proton/index.html) as its underlying AMQP client. Apache Qpid Proton has to be installed to use `amqpctl`. For some features such as SSL support and SASL authentication you will also need OpenSSL or CyrusSASL installed.

The tool is written using [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) libraries for creating CLI application in Go language.

## Installation

`amqpctl` is currently dependent on dynamicaly loaded Qpid Proton libraries. To install `amqpctl` you have to:

1. Install Apache Qpid Proton from your OS repositories or build it from source codes
2. `go get github.com/scholzj/amqpctl`

## Configuration

### From command line

Connection to the AMQP Management node can be configured form the commend line. There are several flags to configure the connection details such as authentication or encryption. For more details please check the help.

```
$ amqpctl --help
...
  -h, --help                    help for amqpctl
  -b, --hostname string         AMQP hostname (default localhost) (default "localhost")
      --password string         AMQP password
  -p, --port int                AMQP port (default 5672) (default 5672)
      --sasl-mechanism string   AMQP SASL mechanism
      --ssl-ca string           SSL certification authority certificate(s)
      --ssl-cert string         SSL certificate for client authentication
      --ssl-key string          SSL private key for client authentication
      --ssl-skip-verify         Skip hostname verification
      --username string         AMQP username
...
```

### Using config file

`amqpctl` configuration can be also conviniently passed from config file. You can specify the file from command line as `amqpctl --config /path/to/my/config.yaml ...`. `amqpctl` will also automatically load the config file names .amqpctl.yaml from the current directory or from your home directory.

The config file should be a simple YAML list of keys and values. The keys correspond to the connection configuratin flags from command line. for example:
```yaml
hostname: localhost
port: 5672
username: admin
password: 123456
```

## Usage

`amqpctl` supports all operations from the AMQP Management specification. You can try it for example against the [Apache Qpid Dispatch](http://qpid.apache.org/components/dispatch-router/index.html) AMQP router.

### QUERY

QUERY operation retrieves selected attributes of Manageable Entities that can be read at this Management Node.

* Query all entities
```
$ amqpctl query
TYPE                                  NAME                                             IDENTITY
org.apache.qpid.dispatch.management   self                                             self
org.apache.qpid.dispatch.log          log/DEFAULT                                      log/DEFAULT
org.apache.qpid.dispatch.log          log/HTTP                                         log/HTTP
...
```

* Query all entities with given type
```
$ amqpctl query listener
TYPE                                NAME                    IDENTITY
org.apache.qpid.dispatch.listener   listener/0.0.0.0:amqp   listener/0.0.0.0:amqp
```

* List only specific attributes
```
$ amqpctl query listener name identity host port
HOST      NAME                    IDENTITY                PORT
0.0.0.0   listener/0.0.0.0:amqp   listener/0.0.0.0:amqp   amqp
```

* List all available attributes
```
$ amqpctl query logStats --all-attributes
DEBUGCOUNT   CRITICALCOUNT   TRACECOUNT   WARNINGCOUNT   INFOCOUNT   NOTICECOUNT   TYPE                                ERRORCOUNT   IDENTITY                NAME
286          0               0            0              1           0             org.apache.qpid.dispatch.logStats   25           logStats/AGENT          AGENT
110          0               112          0              2           0             org.apache.qpid.dispatch.logStats   0            logStats/POLICY         POLICY
0            0               2978         1              219         4             org.apache.qpid.dispatch.logStats   0            logStats/SERVER         SERVER
0            0               3            0              0           50            org.apache.qpid.dispatch.logStats   0            logStats/CONTAINER      CONTAINER
0            0               1457         0              5           0             org.apache.qpid.dispatch.logStats   0            logStats/ROUTER_CORE    ROUTER_CORE
...
```

### READ

READ operation retrieves attributes of a Manageable Entity. The entity can be identified by name, identity or by other attribute (other attributes supported only in WD11).

* Select enetity based on identity
```
$ amqpctl read identity/listener/0.0.0.0:amqp
ATTRIBUTE            VALUE
idleTimeoutSeconds   16
maxFrameSize         16384
maxSessions          32768
multiTenant          false
host                 0.0.0.0
role                 normal
logMessage           none
port                 amqp
identity             listener/0.0.0.0:amqp
name                 listener/0.0.0.0:amqp
stripAnnotations     both
addr                 
http                 false
linkCapacity         1000
requireEncryption    false
type                 org.apache.qpid.dispatch.listener
authenticatePeer     true
requireSsl           false
saslMechanisms       PLAIN DIGEST-MD5 CRAM-MD5
cost                 1
```

* Select entity based on name
```
$ amqpctl read name/self
ATTRIBUTE   VALUE
identity    self
type        org.apache.qpid.dispatch.management
name        self
```

### UPDATE

UPDATE operation update a Manageable Entity. The entity can be identified by name, identity or by other attribute (other attributes supported only in WD11).

```
$ amqpctl update identity/log/SERVER enable=trace+
ATTRIBUTE   VALUE
enable      trace+
module      SERVER
type        org.apache.qpid.dispatch.log
identity    log/SERVER
name        log/SERVER
```

### CREATE

CREATE operation creates a new Manageable Entity. The type attribute is mandatory.

```
$ amqpctl create type=listener name=myListener port=9999 host=0.0.0.0 saslMechanisms=ANONYMOUS
ATTRIBUTE            VALUE
authenticatePeer     false
name                 myListener
idleTimeoutSeconds   16
saslMechanisms       ANONYMOUS
maxFrameSize         16384
cost                 1
logMessage           none
maxSessions          32768
port                 9999
addr                 
multiTenant          false
requireEncryption    false
host                 0.0.0.0
requireSsl           false
http                 false
type                 org.apache.qpid.dispatch.listener
stripAnnotations     both
role                 normal
identity             listener/0.0.0.0:9999:myListener
```

### DELETE

DELETE operation deletes a Manageable Entity. The entity can be identified by name, identity or by other attribute (other attributes supported only in WD11).

```
$ amqpctl delete name/myListener
Manageable Entity successfully deleted.
```

### GET TYPES

GET TYPES operation reutns all implemented entity types.

* Get all entity types
```
$ amqpctl get types                  
TYPE                                               PARENT
org.apache.qpid.dispatch.policy                    org.apache.qpid.dispatch.configurationEntity, org.apache.qpid.dispatch.entity
org.apache.qpid.dispatch.fixedAddress              org.apache.qpid.dispatch.configurationEntity, org.apache.qpid.dispatch.entity
org.apache.qpid.dispatch.router.config.linkRoute   org.apache.qpid.dispatch.configurationEntity, org.apache.qpid.dispatch.entity
...
```

* Get only entity types of selected type
```
$ amqpctl get types operationalEntity
TYPE                                           PARENT
org.apache.qpid.dispatch.logStats              org.apache.qpid.dispatch.operationalEntity, org.apache.qpid.dispatch.entity
org.apache.qpid.dispatch.allocator             org.apache.qpid.dispatch.operationalEntity, org.apache.qpid.dispatch.entity
org.apache.qpid.dispatch.vhostStats            org.apache.qpid.dispatch.operationalEntity, org.apache.qpid.dispatch.entity
org.apache.qpid.dispatch.operationalEntity     org.apache.qpid.dispatch.entity
org.apache.qpid.dispatch.connection            org.apache.qpid.dispatch.operationalEntity, org.apache.qpid.dispatch.entity
org.apache.qpid.dispatch.router.node           org.apache.qpid.dispatch.operationalEntity, org.apache.qpid.dispatch.entity
org.apache.qpid.dispatch.org.amqp.management   org.apache.qpid.dispatch.operationalEntity, org.apache.qpid.dispatch.entity
org.apache.qpid.dispatch.router.address        org.apache.qpid.dispatch.operationalEntity, org.apache.qpid.dispatch.entity
org.apache.qpid.dispatch.management            org.apache.qpid.dispatch.org.amqp.management, org.apache.qpid.dispatch.operationalEntity, org.apache.qpid.dispatch.entity
org.apache.qpid.dispatch.router.link           org.apache.qpid.dispatch.operationalEntity, org.apache.qpid.dispatch.entity
```

### GET ATTRIBUTES

GET ATTRIBUTES operation returns attributes for different entity types.

* List attributes for all entities
```
$ amqpctl get attributes
TYPE                                               ATTRIBUTES
org.apache.qpid.dispatch.configurationEntity       name, identity, type
org.apache.qpid.dispatch.sslProfile                certDb, certFile, keyFile, passwordFile, password, uidFormat, displayNameFile, name, identity, type
org.apache.qpid.dispatch.linkRoutePattern          prefix, dir, connector, name, identity, type
...
```

* List attributes only for particular entity
```
$ amqpctl get attributes logStats
TYPE                                ATTRIBUTES
org.apache.qpid.dispatch.logStats   traceCount, debugCount, infoCount, noticeCount, warningCount, errorCount, criticalCount, name, identity, type
```

### GET ANNOTATIONS

GET ANNOTATIONS operations returns all annotations implemented by different entity types.

```
$ amqpctl get annotations
```

### GET MANAGEMENT NODES

GET MGMTNODES operation lists all management ndoes known to the current management node.

```
$ amqpctl get mgmtnodes
```

### HELP

See `amqpctl help` for more details about the usage and `amqpctl help <operation>` supported operations.

## TODO

* Unit tests
* Integration tests
* Create a static build so that not everyone has to install Proton
* Add possibility to call custom operations
