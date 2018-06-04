# Virtual Machine Migrator Interface Specification 

## Version
This is VMMI spec version 0.1.0 (201806)

## Overview

This document proposes a generic plugin-based migration policy solution for virtual machine managers on Linux, the _Virtual Machine Migrator Interface_, or _VMMI_.
In this document we always use the term 'migration' as shortcut for 'live migration' - the process of migrating a VM while the guest is running, with minimal or no disruption
to the Guest workload, and without the guest noticing.

Migrating a VM is not a simple task, with many tradeoffs and knobs to consider. Each different management applications implements their own logic and policies to adapt to use cases.
This solutions aims to provide
1. a common interface to make the migration policies pluggable.
2. a mean to make the migration policies thus interchangeable across management applications.

The key words "must", "must not", "required", "shall", "shall not", "should", "should not", "recommended", "may" and "optional" are used as specified in [RFC 2119][rfc-2119].

## General Considerations

- the management application is built on top of libvirt
- the VMMI migration model is limited to the [managed peer to peer model](https://libvirt.org/migration.html#flowpeer2peer)
- the management application must set up anything the VM requires to run (e.g. shared storage) before to engage the VMMI plugins
- the management application is in charge to clean up the resources required by the VM to run
- the usage of a VMMI plugin replaces any call to the [virDomainMigrateToURI API family](https://libvirt.org/html/libvirt-libvirt-domain.html#virDomainMigrateToURI3)

### Messages: exchanging data with the client application

Any exchange of data between a VMMI plugin and the client application must happen through JSON messages.
Each JSON message must have at least those keys:

- "vmmiVersion" (string): version of the VMMI SPEC this message complies to. Current value must be "0.1.0"
- "contentType" (string): describes the content of the message, e.g the not-specified, message-variable fields.
  The value of the key is not specified. Some values, however, are reserved for well-known messages. The reserved fields are
  - "configuration": for plugin configuration
  - "migrationprogress": for migration progress updates
  - "migrationcomplete": for migration completion report, either succesfull or not
  - "error": for error message(s).

Each message except the ones with "contentType": "configuration" must have one additional mandatory field

- "timestamp" (uint): seconds since UNIX EPOCH (1970/01/01 00:00:00.00) - aka UNIX time


## VMMI Plugin

### Overview

A VMMI compliant plugin is a single executable which will be placed under the VMMI canonical directory.
The VMMI canonical directory is `/usr/libexec/vmmi`.
A VMMI compliant plugin may have any name as long as it is both a valid UTF8 name and a valid filesystem entry.
The only reserved name is `migrate`, which must not be used and it is reserved to the implementation

### Lifecycle

A VMMI plugin is an operating system process which starts when the migration begins, and exits when "termination"
conditions are met (see "Termination" section).
A VMMI compliant plugin must never exceed the lifetime of a migration, except for the necessary termination
and cleanup duties.
A VMMI plugin is executed when the migration starts, but it must not perform any change to the system,
including actually starting the migration using the libvirt APIs, until it got the configuration data (see Configuration)

### Termination

A VMMI compliant plugin is expected to terminate in the following cases:
- crash (bug)
- signal received - see "signal handling" below
- libvirt reports the migration completed - either succesfully, or aborted

If the plugin detects the operation is aborted, because of a bug, a signal received, or because notified by libvirt,
it must do everything possible to properly clean up, freeing resources and leaving the system in a consistent state.

### Failure model

The introduction of VMMI plugins adds another entity between libvirt and the management application.
If libvirt crashes, or if the VMMI plugin loses the connection to libvirt for any reason, it must exit with error.
If the VM under migration disappears for any reason outside the control of a VMMI plugin, the VMMI plugin must exit with error.
Please note that a VMMI plugin must expect and handle only those teo termination conditions: migration completed, migration aborted.
Please note that the management application is free to abort the migraton (virDomainAbortJob) anytime outside the control of the VMMI plugin.
A VMMI compliant plugin must treat the "migration aborted by the management application" state like the "migration aborted by the hypervisor"
state and exit as soon as possible.
If the VMMI plugin crashes, it must NOT try to abort the current migration process. See also the section "Recovering state".
If the management application crashes, or terminates while the VMMI plugin is still running, the VMMI plugin must continue running as usual.

### Signal handling

A VMMI compliant plugin must react to the following signal


- SIGKILL: exit early as possible, must NOT attempt to cleanup, must NOT abort the current migration
- SIGTERM: exit early as possible, but MUST free any resources and clean up the system, must NOT abort the current migration
- SIGSTOP: abort the current migration, perform any other operation like SIGTERM was received.


### Parameters

A VMMI compliant plugin must support the following parameters:


- VM uuid (required) - uuid of the VM to migrate.
- destination URI (required) - URI of the migration destination.
- the path to the configuration file (optional).


The path may be a simple hyphen ('-'), which signals the plugin to read the configuration from standard input. The plugin must not perform any change to the system, including outputting any data to stdout or stderr, until it got either the configuration or error.
If no path or no hyphen is specified, the plugin must either use its defaults and attempt to read the default configuration file: `/etc/vmmi/conf.d/<plugin>.conf`.
A VMMI compliant plugin must abort with error if the configuration file is given, but the data is malformed or partial.
A VMMI compliant plugin must abort if it attempts to read the default configuration file and fails because the data is malformed or partial.

### Configuration

Each plugin has *two* configuration files.
The first configuration file is `/etc/vmmi/migrate.conf`. It holds the parameters *each* plugin must support.
The second configuration file is either given as plugin parameter, or the default file 

A VMMI compliant plugin must not abort if no configuration file is available, including the default files `/etc/vmmi/conf.d/<plugin>.conf`. The given path always overrides the default.
A VMMI compliant plugin must process either file, and shall not, under any circumstance, process both files.
A VMMI compliant plugin must update its configuration file using this resolution order: defaults, common configuration file, specified configuration. The most recent data must always overwrite the old
A VMMI compliant plugin which reads its configuration from standard input should not expect the standard input to be closed after it received the configuration data. However, it must ignore any additional data which
may be sent through the standard input

#### Common configuration file

The common configuration file defines the following keys:

- "connection" (string): specifies how to connect to libvirt. Default is "qemu:///system"

- "verbose" (int): sets the plugin verbosiness. A plugin sends output using stdout and stderr (see below). The following values are defined:
  * 0: the plugin is completely silent except for fatal error messages
  * 1: in addition to level 0, the plugin also reports the result of the migration
  * 10: in addition to level 1, the plugin also sends progress reports when progress data is available from libvirt
  Default is 1

- "progressReportRate" (string): specifies a time interval to report progress.
  * if set to zero, report progress when data is available from libvirt, if the `verbose` flag allows so
  * if set to any time, sends progress report each _time_ interval, even if not changed since last update.
  Default is 0

Example configuration file
```
  {
    "vmmiVersion": "0.1.0",
    "contentType": "configuration",
    "configuration": {
      "connection": "qemu:///system",
      "verbose": 10,
      "progressReportRate": "30s"
    }
  }
```

#### Plugin-specific configuration file

TBD

### Plugin runtime data

TBD

### Plugin input

Currently not supported. A compliant plugin should expect not to receive input from any source other than the configuration data.

### Reporting output

A VMMI compliant plugin may use both standard output (stdout) and standard error (stderr) to communicate with the client application.
It must not assume any other communication channel.
A VMMI compliant plugin may log other data to other channels (private log file, system log) using any other means, but it must not assume
the client application access those messages.
A VMMI compliant plugin must not assume stdout and stderr are separate channel. A client application is free to multiplex them both
in another channel (e.g. a socket, a pipe)
A VMMI compliant plugin must enforce the sequence of the messages and should not assume any implicit ordering provided by the channel.
In other words, a VMMI compliant plugin must ensure the messages it emits are ordered from the source.

A VMMI compliant plugin may send progress report messages if its configuration allows so.
If a VMMI plugin sends progress report it must send at least two messages signaling zero progress and complete progress.

A VMMI compliant plugin may send a migration completion message. A migration completion explicitely signal the end of migration.
Please note that this is a different message with respect to `complete progress` message.
If a VMMI plugin sends a migration status message, this message must be sent after `progress complete` message.
If a VMMI plugin sends a migration status message, this message must be the last message sent to stdout.

Example messages:

```
  {
    "vmmiVersion": "0.1.0",
    "timestamp": 1528117283,
    "contentType": "migrationProgress",
    "migrationProgress": {
      "percentage": "48",
      "iteration": 2,
    }
  }
```
```
  {
    "vmmiVersion": "0.1.0",
    "timestamp": 1528117312,
    "contentType": "migrationComplete",
    "migrationComplete": {
      "downtime": "10ms"
    }
  }
```


Example valid message sequences:

  1. migration status


  1. migration progress: 0%
  2. migration progress: 100%
  3. migration status


  1. migration progress: 0%
  2. migration progress: 25%
  3. migration progress: 50%
  4. migration progress: 75%
  5. migration progress: 100%
  6. migration status
  

### Reporting errors


Example messages:
```
  {
    "vmmiVersion": "0.1.0",
    "timestamp": 1528117329,
    "contentType": "error",
    "error": {
      "code": 42,
      "message": "generic error",
      "details": "generic error explained in a user-friendly way"
    }
  }
```
