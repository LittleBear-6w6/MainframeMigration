```bash:usage
xxd -p -g 0 -c [レコード長] target.dat | sed s/^/[レコード長16進数]0000/g | xxd -p -r > conv_target.dat
```
