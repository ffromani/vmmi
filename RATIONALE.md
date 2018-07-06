# Virtual Machine Migrator Interface Rationale

This document provides the high-level rationale for the decisions that lead to the VMMI specification.
To read the rationale of the VMMI specification details, see [the spec](https://github.com/fromanirh/vmmi/blob/master/SPEC.md)

## Policies placement

*VMMI puts the implementation of the policies in a third entity different from both the management application and the libvirt*

Various management application implement their migration policies.
Libvirt - the common infrastructure - should not implement policies. It should, like already does, offer tunables, knobs and protocols, letting the upper layers implementing their policies.
To reduce duplication and improve interoperability, the only option left is isolate the migration policies in a third entity and add it to the picture: the VMMI helpers.
Any management application using VMMI has still the option to have built-in policies.
VMMI just offers more flexibility.

## Language-agnostic

*VMMI is language agnostic, and uses well-known, largely available tools to integrate with the existing stack*

The VMMI specs involve calling external processes using the standard Linux tools, and exchanging JSON messages, which is a widely available and simple format.
Using external processes to implement helpers poses no constraint to the implementation language.

## Shared objects vs processes

*VMMI helpers are implemented with standard UNIX processes*

Using processes to implement the VMMI helper, compared with shared objects (.so files) maximizes the flexibility, and sets the lowest possible entry barrier.
Furthermore, using processes to implement VMMI helper, we make impossible for a bad VMMI helper to make crash either libvirt or the management application.

## One process vs multiple processes

*Each VMMI helper implements one migration policy*

An implementation option could have been to integrate a policy engine in a helper, or even in libvirt.
We believe this approach is unpractical, because it requires the embedding of a programming language, making the complexity skyrocket.
Please note that the VMMI specifications doesn't intentionally make impossible to embed such an engine (or embedded language like LUA
or a lisp dialect) inside a VMMi compliant helper, but it also does nothing to encourage this approach.
The simplest possible solution is to implement each policy in a separate helper process, to maximize isolation and to reduce the API surface.
This is the preferred implementation of VMMI helpers.

## Split of commandline arguments vs configuration message

It can be argued that the configuration message/file can entirely replace the command line parameters, making them redundant.
While this is true, I believe this makes things a little less practical, thus I believe the split has merit.
Command line parameters carry the information which is expected to change on every migration, while the configuration message holds
data which should change rarely, if ever.
