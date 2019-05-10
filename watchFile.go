package watchFile

import (
	"os"
	"syscall"
)

type watchFile struct {
	F    *os.File
	name string
	flag int
	perm os.FileMode
	dev  int32
	ino  uint64
}

func (wf *watchFile) Open() error {
	f, err := os.OpenFile(wf.name, wf.flag, wf.perm)
	if err != nil {
		return err
	}
	wf.F = f

	fileInfo, err := os.Stat(wf.name)
	if err != nil {
		wf.F.Close()
		return err
	}

	stat := fileInfo.Sys().(*syscall.Stat_t)
	wf.dev = stat.Dev
	wf.ino = stat.Ino
	return nil
}

// Reopen file if needed.
func (wf *watchFile) Reopen() error {
	fileInfo, err := os.Stat(wf.name)
	if err != nil {
		wf.F.Close()
		return wf.Open()

	} else {

		stat := fileInfo.Sys().(*syscall.Stat_t)
		if wf.dev != stat.Dev || wf.ino != stat.Ino {
			wf.F.Close()
			return wf.Open()
		}
	}
	return nil
}

func (wf *watchFile) WriteString(s string) (int, error) {
	err := wf.Reopen()
	if err != nil {
		return 0, err
	}
	return wf.F.WriteString(s)
}

func (wf *watchFile) Write(p []byte) (int, error) {
	err := wf.Reopen()
	if err != nil {
		return 0, err
	}
	return wf.F.Write(p)
}

func (wf *watchFile) Close() error {
	return wf.F.Close()
}

func OpenWatchFile(name string) (*watchFile, error) {
	var wf = watchFile{name: name, flag: os.O_WRONLY | os.O_APPEND | os.O_CREATE, perm: 0666}
	err := wf.Open()
	return &wf, err
}

func OpenWatchFile2(name string, flag int, perm os.FileMode) (*watchFile, error) {
	var wf = watchFile{name: name, flag: flag, perm: perm}
	err := wf.Open()
	return &wf, err
}
