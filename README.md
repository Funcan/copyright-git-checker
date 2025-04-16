Look for HPE standard copyright notices in any files changed in the latest
commit and alert if any of them need the date updating.

Saves these getting caught by the gate checks which is a slower feedback loop.

Drop the relevent release into ~/bin then setup:

~/.gitconfig :
```
[core]
    excludesfile = /home/duncan/.gitignore
    hooksPath = /home/duncan/.githooks
```

~/.githooks/pre-commit :
```
if [[ -e ~/bin/hpe-copyright-checker ]]; then
  if [[ "$(git config --get guards.skip-hpe-copyright-checker)" != "true" ]]; then
    ~/bin/hpe-copyright-checker
    if [[ $? -ne 0 ]]; then
      echo "HPE copyright checker failed"
      echo 'Use "git config guards.skip-hpe-copyright-checker true" to skip'
      exit 1
    fi
  fi
fi
```
