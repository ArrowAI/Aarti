## aarticlient completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	aarticlient completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
aarticlient completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
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

* [aarticlient completion](aarticlient_completion.md)	 - Generate the autocompletion script for the specified shell

