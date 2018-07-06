# Virtual Machine Migrator Interface Specification 

## Version
This is VMMI spec version 0.4.1

## Overview

This document proposes a generic migration policy approach for virtual machine managers on Linux, the _Virtual Machine Migrator Interface_, or _VMMI_.
In this document we always use the term 'migration' as shortcut for 'live migration' - the process of migrating a VM while the guest is running, with minimal or no disruption
to the Guest workload, and without the guest noticing.

Migrating a VM is not a simple task, with many tradeoffs to consider and knobs to tune to make the operation succeed.
Each different management applications implements their own logic and policies to adapt to use cases.

This solutions aims to provide

1. a common interface to make the migration policies pluggable.
2. a mean to make the migration policies thus interchangeable across management applications.

The key words "must", "must not", "required", "shall", "shall not", "should", "should not", "recommended", "may" and "optional" are used as specified in [RFC 2119](https://www.ietf.org/rfc/rfc2119.txt).

## General Considerations

- the management application is built on top of libvirt, so it uses the [libvirt APIs](https://libvirt.org/html/index.html) to manage VMs, including migrations.
- the preferred VMMI migration model is the [managed peer to peer model](https://libvirt.org/migration.html#flowpeer2peer), but there is no constraint enforced on the migration model.
- the management application must set up anything the VM requires to run (e.g. shared storage) before to engage the VMMI implementations.
- the management application is in charge to clean up the resources required by the VM to run.
- the usage of a VMMI implementation replaces any call to the [virDomainMigrateToURI API family](https://libvirt.org/html/libvirt-libvirt-domain.html#virDomainMigrateToURI3).

## Design principles

The purpose of a VMMI implementation is to fully encapsulate a migration policy and its implementation in a self contained unit of code (the callee), opaque to the management application (the caller).

The management application offloads completely the task to the implementation, with no dynamic tuning.

In other words, the management application is expected to run the VMMI implementation and wait for its completion, reading back only the operation status once the operation completed - whatever the result.
For this reason, the VMMI implementation must stay silent with respect to its caller: any message it may sent is not-actionable, thus useless.
The VMMI implementation and the management application will use independent connection to libvirt to consume the migration stats.


### Messages: exchanging data with the outside world

Any exchange of data between a VMMI implementation and anything else.
Each JSON message must have at least those keys:

- "vmmiVersion" (string): version of the VMMI SPEC this message complies to. Current value must be "0.4.1"
- "contentType" (string): describes the content of the message, e.g the not-specified, message-variable fields.
  The value of the key is not specified. Some values, however, are reserved for well-known messages. The reserved fields are
  - "configuration": for configuration info
  - "completion": for migration completion report

The following fields are optional for "contentType": "configuration", and mandatory for each other message

- "timestamp" (uint): seconds since UNIX EPOCH (1970/01/01 00:00:00.00) - aka UNIX time

Example message:
```
  {
    "vmmiVersion": "0.4.1",
    "timestamp": 1528117452,
    "contentType": "foobar",
    "foobar": {
      "fizz": 42,
      "buzz": "baz"
    }
  }
```

#### Input/Output principles

A VMMI compliant implementation:

- must read the configuration message from either its stdin or a a file path given as parameter (see "Parameters");
the implementation must not expect any other inbound message.

- must report the final status of its execution to stderr. This may be either success or error. This is the
only message a VMMI compliant implementation must send.

- should use stdout to send any other message. The implementation should assume that writing to stdout never blocks.
It is a responsability to the management app to ensure this is true. Please note that the management app is *not* required
to consume those messages, it can just discard them - even just binding the stdout to /dev/null.

- must send messages to stdout only when requested by the management application, and must never send unsolicited messages.
The management application must send a specific signal (see "Signal handling") to request a new status message to be sent.
The VMMI implementation must send a status message as soon as possible once the signal is received.
The VMMI implementation must send *at most* one status message for each signal received.
Both the VMMI implementation and the management application must treat the delivery of status messages as best effort.

- may log other data to other channels (private log file, system log) using any other means, but it must not assume
the client application reads those messages.

- must enforce the sequence of the messages and should not assume any implicit ordering provided by the channel.
In other words, a VMMI compliant implementation must ensure the messages it emits are ordered from the source.
This is done using the "timestamp" field of the messages.

- must not abort if no configuration file is available.

- must update its configuration file using this resolution order:
defaults, specified configuration. The most recent data must always overwrite the old

- if the configuration is read from standard input, then it should not expect the standard input
to be closed after it received the configuration data. However, it must ignore any additional data which
may be sent through the standard input

- must wait for the configuration message to be read before to perform any change to the system.

#### Message ordering

In the simplest case, a VMMI implementation just sends the Completion message to the management app (see specification below).

The implementation can also send status messages not yet specified. In any case, the implementation must signal the ordering
using the "timestamp" field of the messages.

This applies
- between status messages sent to stdout
- between messages sent to stdout and the completion message sent to stderr

The implementation must always take in account that the management app is free to multiplex stdout and stderr
in the same channel (e.g. a pipe, or a socket)

#### Message specification: Configuration

The configuration data of each implementation must support at least the following keys:

- "connection" (string): specifies how to connect to libvirt. Default is "qemu:///system"

- "verbose" (int): sets the implementation verbosiness. A implementation sends output using stdout and stderr (see below).
  The following values are defined:
  * 0: the implementation is completely silent except for fatal error messages

Example configuration message
```
  {
    "vmmiVersion": "0.4.1",
    "contentType": "configuration",
    "configuration": {
      "connection": "qemu:///system",
      "verbose": 10,
    }
  }
```

A VMMI compliant implementation is expected to honour the above keys. It cannot read and silently discard them.


#### Message specification: Completion

A VMMI compliant implementation must always report its termination -except for crashes- sending a Completion message to its client.

The "Completion" message has a different layout depending on the migration terminated successfully or with error.
Should the migration succeed, a VMMI compliant implementation must signal this state sending an Completion message with a "success" payload.

The completion payload has one mandatory key, "result".

The value of the "result" key may be either "success" for migration completed correctly, or "error" otherwise.
Like the top-level "contentType" key, the "completion" payload will have another key depending on the value of "result".

- if "result" equals to "success", the "completion" payload must have another key "success", whose value must be an object.
The content of that object is not specified. A empty object is a valid value.
Example of succesfull termination
```
  {
    "vmmiVersion": "0.4.1",
    "timestamp": 1528117329,
    "contentType": "completion",
    "completion": {
      "result": "success",
      "success": {}
    }
  }
```

- if "result" equals to "error", the "completion" payload must have another key "error", whose value must be an object which
must hold three more keys:
  - "code" (int): error code
  - "message" (string): succint message describing the error, like a log entry
  - "details" (string): more user-friendly and verbose error description.

Example of failed migration report message:
```
  {
    "vmmiVersion": "0.4.1",
    "timestamp": 1528117329,
    "contentType": "completion",
    "completion": {
      "result": "error",
      "error": {
        "code": 42,
        "message": "generic error",
        "details": "generic error explained in a user-friendly way"
      }
    }
  }
```

## VMMI Implementations

### Overview

A VMMI compliant implementation is a single executable which will be placed under the VMMI canonical directory.
The VMMI canonical directory is `/usr/libexec/vmmi`. (*TODO: this may change before the spec is finalized*)

A VMMI compliant implementation may have any name as long as it is both a valid UTF8 name and a valid filesystem entry.
The only reserved name is `migrate`, which must not be used and it is reserved to the implementation

### Lifecycle

A VMMI compliant implementation:

- must be implemented using an operating system process which starts when the migration begins, and exits when "termination"
conditions are met (see "Termination" section).
- must never exceed the lifetime of a migration, except for the necessary termination
and cleanup duties.
- is executed when the migration starts, but it must not perform any change to the system,
including actually starting the migration using the libvirt APIs, until it got the configuration data (see Configuration)

### Termination

A VMMI compliant implementation is expected to terminate in the following cases:
- crash (bug)
- signal received - see "signal handling" below
- libvirt reports the migration completed - either successfully, or aborted

If the implementation detects the operation is aborted, because of a bug, a signal received, or because notified by libvirt,
it must do everything possible to properly clean up, freeing resources and leaving the system in a consistent state.

### Failure model

The introduction of VMMI implementations adds another entity between libvirt and the management application.

If libvirt crashes, or if the VMMI implementation loses the connection to libvirt for any reason, it must exit with error.

If the VM under migration disappears for any reason outside the control of a VMMI implementation, the VMMI implementation must exit with error.
Please note that a VMMI implementation must expect and handle only those termination conditions: migration completed, migration aborted.
Please note that the management application is free to abort the migration (e.g. calling virDomainAbortJob) anytime outside the control of the VMMI implementation.

A VMMI compliant implementation must treat the "migration aborted by the management application" state like the "migration aborted by the hypervisor"
state and exit as soon as possible.

If the VMMI implementation crashes, it must NOT try to abort the current migration process. See also the section "Recovering state".
If the management application crashes, or terminates while the VMMI implementation is still running, the VMMI implementation must continue running as usual.

### Recovering state

When a VMMI compliant implementation starts, it must check if a migration is already in progress. If so, it must take control of it and start enforcing
the policy. A VMMI compliant implementation must not assume that a migration currently in progress was started by a previous instance of itself.

Should a management application using VMMI detect the crash of a VMMI implementation, it is free to take any action, including running again the
same VMMI implementation, running a different one or do nothing.

### Signal handling

A VMMI compliant implementation must react to the following signals


- SIGUSR1: send a new status message to stdout as soon as possible.
- SIGTERM: exit early as possible, but MUST free any resources and clean up the system, must NOT abort the current migration.
- SIGSTOP: abort the current migration, perform any other operation like SIGTERM was received.
- SIGINT: like SIGSTOP.


### Parameters

A VMMI compliant implementation must support the following parameters:


- VM uuid (required) - uuid of the VM to migrate.
- destination URI (required) - URI of the migration destination.
- migration URI (required) - QEMU-specific migration destination.
- configuration file path (optional) - the path to the configuration file, or '-' for stdin.


The path may be a simple hyphen ('-'), which signals the implementation to read the configuration from standard input. The implementation must not perform any change to the system,
including outputting any data to stdout or stderr, until it got either the configuration or error.
If no path or no hyphen is specified, the implementation must either use its defaults and attempt to read the default configuration file: `/etc/vmmi/conf.d/<implementation>.json`.

A VMMI compliant implementation:

- must abort with error if the configuration file is given, and the data is malformed or partial.
- must abort if it attempts to read the default configuration file and fails because the data is malformed or partial.

### Plugin runtime data

Each VMMI compliant implementation is responsible for storing, cleaning up and recovering any runtime data it may need.
