package watchFile

import (
	"os"
	"syscall"
)

type WatchFile struct {
	F    *os.File
	Name string
	Flag int
	Perm os.FileMode
	Dev  uint64
	Ino  uint64
}

func (wf *WatchFile) Open() error {
	f, err := os.OpenFile(wf.Name, wf.Flag, wf.Perm)
	if err != nil {
		return err
	}
	wf.F = f

	fileInfo, err := os.Stat(wf.Name)
	if err != nil {
		wf.F.Close()
		return err
	}

	stat := fileInfo.Sys().(*syscall.Stat_t)
	wf.Dev = uint64(stat.Dev)
	wf.Ino = uint64(stat.Ino)
	return nil
}

// Reopen file if needed.
func (wf *WatchFile) Reopen() error {
	fileInfo, err := os.Stat(wf.Name)
	if err != nil {
		wf.F.Close()
		return wf.Open()

	} else {

		stat := fileInfo.Sys().(*syscall.Stat_t)
		if wf.Dev != uint64(stat.Dev) || wf.Ino != uint64(stat.Ino) {
			wf.F.Close()
			return wf.Open()
		}
	}
	return nil
}

func (wf *WatchFile) WriteString(s string) (int, error) {
	err := wf.Reopen()
	if err != nil {
		return 0, err
	}
	return wf.F.WriteString(s)
}

func (wf *WatchFile) Write(p []byte) (int, error) {
	err := wf.Reopen()
	if err != nil {
		return 0, err
	}
	return wf.F.Write(p)
}

func (wf *WatchFile) Close() error {
	return wf.F.Close()
}

func OpenWatchFile(name string) (*WatchFile, error) {
	var wf = WatchFile{Name: name, Flag: os.O_WRONLY | os.O_APPEND | os.O_CREATE, Perm: 0666}
	err := wf.Open()
	return &wf, err
}

func OpenWatchFile2(name string, flag int, perm os.FileMode) (*WatchFile, error) {
	var wf = WatchFile{Name: name, Flag: flag, Perm: perm}
	err := wf.Open()
	return &wf, err
}
