---
sidebar_position: 5
---

# uplift completion

Generates an uplift autocompletion script for your target shell.

```sh
uplift completion [COMMAND]
```

## Commands

### bash

To load the completions in your current shell session:

```sh
source <(uplift completion bash)
```

To Load the completions for every new session:

#### Linux

```sh
uplift completion bash > /etc/bash_completion.d/uplift
```

#### MacOS

```sh
uplift completion bash > /usr/local/etc/bash_completion.d/uplift
```

### zsh

To load the completions in your current shell session:

```sh
source <(uplift completion zsh)
```

To load the completions for every new session:

```sh
uplift completion zsh > "${fpath[1]}/_uplift"
```

### fish

To load the completions in your current shell session:

```sh
uplift completion fish | source
```

To load the completions for every new session:

```sh
uplift completion fish > ~/.config/fish/completions/uplift.fish
```
