common projecttemplate

## 1 install

```bash
cd projecttemplate/main
go mod tidy
```

## 2 build

```bash
make clean && make
```

> + 会生成可执行文件 `projecttemplate`

## 3 start

```bash
./projecttemplate -c output/projecttemplate.toml
```
