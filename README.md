# le
le自动将超出终端界面的stdout以less优雅的展示  

目前支持系统：linux

## install
```bash
git clone https://github.com/helowd/lesser.git && cd lesser && sudo ./install.sh
```

## uninstall
```bash
rm -f /usr/bin/le
```

## usage
```bash
le <command>
```

### example
```bash
le kubectl get --help  
le cat nginx.log
```
