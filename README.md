# watchFile

`go-watch-file` is a Go library similar to Python’s`WatchedFileHandler`, for writing files that may berotated, such as log rotation scenarios.  
  
* **only Unix systems are supported**
  
`go-watch-file` 是一个类似 Python `WatchedFileHandler` 的 Go 库，用于写入可能被轮转的文件，如日志轮转场景。  
  
* **仅支持 Unix 系统**




# version

- latest: `v1.0.3`

# Install

```bash
go get -u github.com/gojuukaze/go-watch-file

# use go mod
go get github.com/gojuukaze/go-watch-file@v1.0.3
# or
go mod edit -require=github.com/gojuukaze/go-watch-file@v1.0.3

```



# Example
```go
package main

import (
	"github.com/gojuukaze/go-watch-file"
	"os"
	"fmt"
)

func main() {

	f, _ := watchFile.OpenWatchFile("./a.txt")
	f.WriteString("111")

	os.Remove("./a.txt")

	f.WriteString("222\n")
	f.Write([]byte("44"))

	f.Close()

	//
	f2, _ := watchFile.OpenWatchFile2("./b.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	f2.WriteString("333")
	f2.Close()

	// use other file method
	f3, _ := watchFile.OpenWatchFile2("./b.txt", os.O_RDONLY|os.O_CREATE, 0666)
	f3.Reopen()
	b := make([] byte, 4)
	f3.F.Read(b)
	fmt.Println(b)
	f3.Close()

}

```
# using with logrus

## set output
```go
f,err:=watchFile.OpenWatchFile("m.log")
if err != nil {
	panic(err)
}
Log := logrus.New()
Log.SetOutput(f)
```

## hook
```go
type FileHook struct {
	f         *watchFile.WatchFile
	formatter logrus.Formatter
}

func (hook FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook FileHook) Fire(entry *logrus.Entry) error {
	msg, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}
	hook.f.Write(msg)
	return nil
}

func NewFileHook(name string) FileHook {
	f, err := watchFile.OpenWatchFile(name)
	if err != nil {
		panic(err)
	}
	var hook = FileHook{
		f:         f,
		formatter: &logrus.TextFormatter{DisableColors: true},
	}
	return hook

}

```

```go
Log = logrus.New()
fHook:=NewFileHook("m.log")
Log.AddHook(fHook)
```
