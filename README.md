# Monetasa
IoT data marketplace based on blockchain.

## Install
```
cd $GOPATH/src
git clone https://gitlab.com/drasko/monetasa
```
## Set Dev Env

### Add SSH Key
Add your `ssh` key to GitLab, then use ssh git remote:
```
git remote set-url origin git@gitlab.com:drasko/monetasa.git
```

### Configure git To Use SSH
Configure `git` to use `ssh` for GitLab, in order to enable
`dep` functionality, as explained [here](https://gist.github.com/shurcooL/6927554) -
otherwise it will break beacuse it can not handle username/password prompt.

```
git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"
```
