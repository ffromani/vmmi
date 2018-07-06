# VMMI - the Virtual Machine Migrator Interface

[![GoDoc](https://godoc.org/github.com/fromanirh/vmmi/pkg/vmmi?status.svg)](https://godoc.org/github.com/fromanirh/vmmi/pkg/vmmi)

## What is VMMI?

VMMI (_Virtual Machine Migrator Interface_) consists of a specification and libraries for writing helpers to implement
live migration of virtual machines across hosts, along with a number of supported helpers.
VMMI concerns itself only with the migration of the VMs, using [libvirt](http://libvirt.org).
The preferred migration mode is the [managed peer to peer](https://libvirt.org/migration.html#flowpeer2peer).

## VMMI highlights:

1. The migration policies are implemented as helper processes.
2. A "migration policy" is any entity that supervises the migration process and that increases the chance of the migration succesfully complete.
3. The helper process own the migration - they start them, they tune its parameters at runtime.
4. The management application can still abort or get the migration status using standard libvirt APIs.
5. The helper process is used as "fire and forget" - the client just need to wait for completion and collect the exit status.
6. The helper process will receive parameters and configuration at startup, no further interaction needed once it started.

## Why Develop VMMI?

Live migrating a Virtual Machine is the process of moving the virtual machine process from one hypervisor to another, usually running on a different host, with minimal to none
disruption to the service(s) provided by the guest, and without the guest noticing.

Performing the live migration is a complex task, which has many possible solutions. Despite being implemented by many management application, is not yet a solved problem.
Assuming that the source and destination sides are fully compatible, thus the live migration process can start and has a chance to successfully complete, still the
The live migration process can fail for many reasons, including the workload of the guest, and the state of the live migration medium (e.g. the network link).

Most management applications implements live migration policies, e.g. monitor the state of the migration and tune the knobs exposed by the hypervisor (in our case, by the libvirt interface)
to help the migration finish (converge) successfully, or abort if a timeout expired.
A migration policy can be thought as a process monitoring the migration state and changing the migration settings according to some rules to produced a desired outcome.

The purpose of VMMI is to encapsulate the migration policies in external entities -the helpers- and make them agnostic with respect to the management application, to make them interchangeable.
VMMI interacts with the management application using well defined interface leveraging JSON messages, and uses libvirt to actually interact with the hypervisor(s).

![VMMI architecture](https://github.com/fromanirh/vmmi/blob/master/Documentation/arch.png "VMMI architecture")

### Requirements

The VMMI spec is language agnostic.
The VMMI spec and the reference helpers all assume libvirt manages the VMs being migrated, and the migration occurs in managed peer to peer mode.

### Reference Plugins

The VMMI project maintains a set of [reference helpers](https://github.com/fromanirh/vmmi/tree/master/helpers) that implement the VMMI specification.

### Running the helpers

Patches to integrate VMMI in popular Virtual Machine Management applications are available for:

* [ovirt](https://gerrit.ovirt.org/#/q/status:open+project:vdsm+branch:master+topic:assisted-migration)

## What might VMMI do in the future?

The first purpose of VMMI is to test if the concept of moving the migration policy in a separate helper out of the management application is viable or not.
If the concept proves itself worthy, we aim to propose integration of VMMI inside libvirt

## Contact

For any questions about VMMI, please reach out using email (fromani at redhat)

