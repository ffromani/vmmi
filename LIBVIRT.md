# Integrating VMMI in libraries

*VMMI requires minor changes to libvirt and to management applications*

Adding new migration policies should be as simple as possible. Having a library with built-ins policies conflicts with this requirement.
The only real option is to have a runtime, pluggable interface.

We are proposing a [libvirt facade to invoke VMMI helpers using the domain API](https://github.com/fromanirh/vmmi/tree/master/patches/libvirt),
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
