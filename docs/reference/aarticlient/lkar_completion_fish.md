## aarticlient completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	aarticlient completion fish | source

To load completions for every new session, execute once:

	aarticlient completion fish > ~/.config/fish/completions/aarticlient.fish

You will need to start a new shell for this setup to take effect.


```
aarticlient completion fish [flags]
```

### Options

```
  -h, --help              help for fish
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

