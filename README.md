# copyright-git-checker

Look for HPE standard copyright notices in any files changed in the latest
commit and alert if any of them need the date updating.

Saves these getting caught by the gate checks which is a slower feedback loop.

## Options

| Flag     | Description                                                                               |
| -------- | ----------------------------------------------------------------------------------------- |
| `-fix`   | Automatically update copyright headers to include the current year and re-stage the files |
| `-quiet` | Suppress informational messages (errors are still shown)                                  |

### How `-fix` works

When `-fix` is specified:

- `Copyright 2025 Hewlett Packard` becomes `Copyright 2025-2026 Hewlett Packard`
- `Copyright 2024-2025 Hewlett Packard` becomes `Copyright 2024-2026 Hewlett Packard`
- Files are automatically re-staged with `git add` so the fixes are included in your commit

## Installation

Drop the relevant release into ~/bin then setup:

~/.gitconfig:

```ini
[core]
    excludesfile = /home/duncan/.gitignore
    hooksPath = /home/duncan/.githooks
```

### Check-only hook (default)

Blocks commits with outdated copyrights until you fix them manually:

~/.githooks/pre-commit:

```sh
#!/bin/sh

if [ -e ~/bin/hpe-copyright-checker ]; then
  if [ "$(git config --get guards.skip-hpe-copyright-checker)" != "true" ]; then
    if ! ~/bin/hpe-copyright-checker; then
      echo "HPE copyright checker failed"
      echo 'Use "git config guards.skip-hpe-copyright-checker true" to skip'
      exit 1
    fi
  fi
fi
```

### Auto-fix hook

Automatically fixes copyright headers on every commit:

~/.githooks/pre-commit:

```sh
#!/bin/sh

if [ -e ~/bin/hpe-copyright-checker ]; then
  if [ "$(git config --get guards.skip-hpe-copyright-checker)" != "true" ]; then
    ~/bin/hpe-copyright-checker -fix -quiet
  fi
fi
```
