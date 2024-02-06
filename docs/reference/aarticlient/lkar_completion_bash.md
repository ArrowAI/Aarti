## aarticlient completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(aarticlient completion bash)

To load completions for every new session, execute once:

#### Linux:

	aarticlient completion bash > /etc/bash_completion.d/aarticlient

#### macOS:

	aarticlient completion bash > $(brew --prefix)/etc/bash_completion.d/aarticlient

You will need to start a new shell for this setup to take effect.


```
aarticlient completion bash
```

### Options

```
  -h, --help              help for bash
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

