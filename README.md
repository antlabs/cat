## cat
cat 是golang 实现的cat命令，完成了cat命令的所有功能，对posix cat命令实现细节感兴趣的童鞋可以看下，同时也是https://github.com/guonaihong/clop 库的使用示例

## install
```
go build -o cat github.com/antlabs/cat/_cmd
```
## usage
```console
Usage:
    ./cat [Flags] <files> 

Flags:
    -A,--show-all            equivalent to -vET
    -E,--show-ends           display $ at end of each line
    -T,--show-tabs           display TAB characters as ^I
    -b,--number-nonblank     number nonempty output lines, overrides
    -e                       
    -n,--number              number all output lines
    -s,--squeeze-blank       suppress repeated empty output lines
    -t                       equivalent to -vT
    -v,--show-nonprinting    use ^ and M- notation, except for LFD and TAB

Args:
    <files>  
```