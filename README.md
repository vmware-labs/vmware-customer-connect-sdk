# vmware-customer-connect-sdk
This SDK builds a layer of abstraction above customerconnect.vmware.com to hide the complexity for the client. This allows for downloads to be requested using the minimum of information.

**WARNING:** This SDK is unofficial and experimental, with no guarantee of API stability.

## Overview

This code is not meant to be compiled directly, but rather be consumed by another compiled program.

### Prerequisites

See `go.mod` for details of the version of Golang used.

Install modules with `go mod download` from the root of the repo.

You must export environmental variables with your credentials to VMware Customer Connect. A generic account with no specific entitlement is required. Make sure that password are enclosed in single quotes to prevent issues with special charactors in passwords.

```
export VMWCC_USER='<user@name>'
export VMWCC_PASS='<password>'
```

## Documentation

Run test with `go test ./...`.

## Contributing

The vmware-customer-connect-sdk project team welcomes contributions from the community. Before you start working with vmware-customer-connect-sdk, please
read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be
signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on
as an open-source patch. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Apache License 


Credit [@hoegaarden] (https://github.com/hoegaarden) for some of the original code.