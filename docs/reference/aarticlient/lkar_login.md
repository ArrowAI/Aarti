## aarticlient login

Login to an Artifact Registry repository

```
aarticlient login [registry] [flags]
```

### Examples

```

Log in with username and password from command line flags:
  aarticlient login -u username -p password localhost:5000

Log in with username and password from stdin:
  aarticlient login -u username --password-stdin localhost:5000

Log in with username and password in an interactive terminal and no TLS check:
  aarticlient login --insecure localhost:5000

```

### Options

```
  -h, --help             help for login
      --password-stdin   Take the password from stdin
```

### Options inherited from parent commands

```
      --ca-file string   CA certificate file
  -d, --debug            Enable debug logging
  -k, --insecure         Do not verify tls certificates
  -p, --pass string      Password
  -H, --plain-http       Use http instead of https
  -u, --user string      Username
```

### SEE ALSO

* [aarticlient](aarticlient.md)	 - An OCI based Artifact Registry

