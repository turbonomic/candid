# candid
A rest client for passing information from CNAE to Turbo for impacted lead nodes.

### Build it
```bash
go build
```

### Set the host, login credentials of the turbonomic.server
```terminal
Usage of ./candid:
  -candidhost string
    	the address of candid.server (default "")
  -candiduser string
    	username to login to candid.server (default "administrator")
  -candidpass string
    	password to login to candid.server (default "")
  -turbohost string
    	the address of turbo.server (default "")
  -turbouser string
    	username to login to turbo.server (default "administrator")
  -turbopass string
    	password to login to turbo.server (default "")
  -logtostderr
    	log to standard error instead of files

```

### Run it
```bash
./candid --candidhost=https://10.10.174.134 --candiduser=admin --candidpasspass=password \
  --turbohost=https://10.10.174.134 --turbouser=administrator --turbopassword=password
```

