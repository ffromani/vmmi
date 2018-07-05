# Virtual Machine Migrator Interface Rationale

This document provides the high-level rationale for the decisions that lead to the VMMI specification.
To read the rationale of the VMMI specification details, see [the spec](https://github.com/fromanirh/vmmi/blob/master/SPEC.md)

## Policies placement

*VMMI puts the implementation of the policies in a third entity different from both the management application and the libvirt*

Various management application implement their migration policies.
Libvirt - the common infrastructure - should not implement policies. It should, like already does, offer tunables, knobs and protocols, letting the upper layers implementing their policies.
To reduce duplication and improve interoperability, the only option left is isolate the migration policies in a third entity and add it to the picture: the VMMI plugins.
Any management application using VMMI has still the option to have built-in policies.
VMMI just offers more flexibility.

## Language-agnostic

*VMMI is language agnostic, and uses well-known, largely available tools to integrate with the existing stack*

The VMMI specs involve calling external processes using the standard Linux tools, and exchanging JSON messages, which is a widely available and simple format.
Using external processes to implement plugins poses no constraint to the implementation language.

## Integration in libraries

*VMMI requires minor changes to libvirt and to management applications*

Adding new migration policies should be as simple as possible. Having a library with built-ins policies conflicts with this requirement.
The only real option is to have a runtime, pluggable interface.
In the future, libvirt may offer a [facade to invoke VMMI plugins using the domain API](https://github.com/fromanirh/vmmi/tree/master/patches/libvirt),
much like we do with current migration APIs. This change is minor and not invasive.
The change required to management application is very simple: just call a different, but similar, API (be it integrated in libvirt or in an ancillary library).

### Integration in libvirt

*Advantages in integrating VMMI support in libvirt*

On one hand, VMMI support doesn't *require* integration in libvirt; on the other hand, libvirt seems a natural fit for such integration.
Proper VMMI integration in a virtual machine management stack requires a process supervisor. That process supervisor must:

1. spawn and wait for the VMMI helper termination
2. handle the VMMI helper I/O as per spec
3. keep an handle for the VMMI helper to facilitate the recovery should the management application crashes: the management application
   needs a way to know if a VMMI-assisted migration is in progress.

Integrating VMMI support in libvirt nicely addresses all the requirements above. Otherwise, proper VMMI integration requires another
(thin) helper library to handle the above tasks. The [libvirt patch](https://github.com/fromanirh/vmmi/blob/master/patches/libvirt/0001-POC-WIP-domain-introduce-virDomainMigrateWithHelper.patch)
is simple and little intrusive.

## Shared objects vs processes

*VMMI plugins are implemented with standard UNIX processes*

Using processes to implement the VMMI plugin, compared with shared objects (.so files) maximizes the flexibility, and sets the lowest possible entry barrier.
Furthermore, using processes to implement VMMI plugin, we make impossible for a bad VMMI plugin to make crash either libvirt or the management application.

## One process vs multiple processes

*Each VMMI plugin implements one migration policy*

An implementation option could have been to integrate a policy engine in a plugin, or even in libvirt.
We believe this approach is unpractical, because it requires the embedding of a programming language, making the complexity skyrocket.
Please note that the VMMI specifications doesn't intentionally make impossible to embed such an engine (or embedded language like LUA
or a lisp dialect) inside a VMMi compliant plugin, but it also does nothing to encourage this approach.
The simplest possible solution is to implement each policy in a separate plugin process, to maximize isolation and to reduce the API surface.
This is the preferred implementation of VMMI plugins.

## Split of commandline arguments vs configuration message

It can be argued that the configuration message/file can entirely replace the command line parameters, making them redundant.
While this is true, I believe this makes things a little less practical, thus I believe the split has merit.
Command line parameters carry the information which is expected to change on every migration, while the configuration message holds
data which should change rarely, if ever.
