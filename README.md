# watchFile
watches the file to see if it has changed since the last write.  
If the file has changed, the old file stream is closed, and the file opened to get a new file.  

This handler is based on a suggestion by python WatchedFileHandler  

**NOTE: only use under Unix**

用于打开会被删除的文件，如果被删除则自动打开新文件。  
**注意：只支持Unix**

它可以用来写会被切割的日志文件

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
