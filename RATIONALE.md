# Virtual Machine Migrator Interface Rationale

This document provides the high-level rationale for the decisions that lead to the VMMI specification.
To read the rationale of the VMMI specification details, see SPEC.md

## Policies placement

*VMMI puts the implementation of the policies in a third entity different from both the management application and the libvirt*

Various management application implement their migration policies.
Libvirt - the common infrastructure - should not implement policies. It should, like already does, offer tunables, knobs and protocols, letting the upper layers implementing their policies.
To reduce duplication and improve interoperability, the only option left is isolate the migration policies in a third entity and add it to the picture: the VMMI plugins.
Any management application using VMMI has still the option to have built-in policies. Using VMMI is about adding extensibility and flexibility.

## Integration in libraries

*VMMI does not currently require any change to libvirt*

Adding new migration policies should be as simple as possible. Having a library with built-ins policies conflicts with this requirement.
The only real option is to have a runtime, pluggable interface.
There is no immediate benefit for integrating that aforementioned interface in libvirt. In the future, libvirt may offer a facade to invoke VMMI plugins using the domain API, much like
we do with current migration APIs, but the benefit is minor.

## Shared objects vs processes

*VMMI plugins are implemented with standard UNIX processes*

Using processes to implement the VMMI plugin, compared with shared objects (.so files) maximizes the flexibility, and sets the lowest possible entry barrier.
Furthermore, the damage done by a faulty policy is minimized.

## One process vs multiple processes

*Each VMMI plugin implements one migration policy*

An implementation option could have been to integrate a policy engine in a plugin, or even in libvirt. We believe this approach is unpractical, because it requires the embedding
of a programming language, making the complexitry explode.
Please note that the VMMI specifications doesn't intentionally make impossible to embed such an engine (or embedded language like lua or a lisp dialect) inside a VMMi compliant plugin,
but it also does nothing to encourage this approach.
The simplest possible solution is to implement each policy in a separate plugin process, to maximize isolation and to reduce the API surface.
This is the preferred implementation of VMMI plugins.
