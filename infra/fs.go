package infra

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"xiyu.com/common"
)

const file_max_size = (200 << 20)

const Aavator = "avator"
const Mind = "mind"
const Img = "img"

type Fs struct {
	prefix  string
	blockNo int
	offset  int64
	file    *os.File
	mu      sync.Mutex
	common.TAutoCloseable
}

var fsMapping = make(map[string]*Fs)

func InitFs() {
	fs := NewFs(Aavator)
	fsMapping[Aavator] = fs
	fs = NewFs(Mind)
	fsMapping[Mind] = fs
	fs = NewFs(Img)
	fsMapping[Img] = fs
}

func GetFs(key string) (*Fs, error) {
	fs, err := fsMapping[key]
	if !err {
		return nil, errors.New("Not exist")
	}
	return fs, nil
}

func (tfs *Fs) Close() {
	if tfs.file != nil {
		tfs.file.Close()
	}
}

func NewFs(prefix string) *Fs {
	fs := &Fs{prefix: prefix, blockNo: 0, offset: 0}
	fs.load()
	common.TaddItem(fs)
	return fs
}

func (tfs *Fs) Read(ref string) ([]byte, error) {
	parts := strings.SplitN(ref, "-", 3)
	if len(parts) != 3 {
		return nil, errors.New("parts error")
	}

	offset, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		common.Logger.Warnf("%s offset error:%s", ref, parts[1])
		return nil, errors.New("parts error")
	}

	len, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		common.Logger.Warnf("%s len error:%s", ref, parts[1])
		return nil, errors.New("parts error")
	}

	name := fmt.Sprintf("%s/%s/%s", common.GlbBaInfa.Conf.Infra.FsDir, tfs.prefix, parts[0])
	file, err := os.OpenFile(name, os.O_RDONLY, 0644)
	if err != nil {
		common.Logger.Warnf("open %s failed:%s", name, err.Error())
		return nil, err
	}

	dat := make([]byte, len)
	defer file.Close()
	n, err := file.ReadAt(dat, offset)
	if err != nil {
		common.Logger.Warnf("read %s failed:%s", name, err.Error())
		return nil, err
	}

	if n != int(len) {
		common.Logger.Warnf("read %s len error: %d != %d", name, n, len)
		return nil, errors.New("read len error")
	}

	return dat, nil
}

func (tfs *Fs) Write(dat []byte) (string, error) {
	tfs.mu.Lock()
	defer tfs.mu.Unlock()
	if (tfs.offset + (int64(len(dat)))) >= file_max_size {
		tfs.file.Close()
		tfs.file = nil
		tfs.blockNo = tfs.blockNo + 1
		err := tfs.openFile()
		if err != nil {
			return "", err
		}
	}
	n, err := tfs.file.Write(dat)
	if err != nil {
		common.Logger.Errorf("write fs:%s-%d failed:%s", tfs.prefix, tfs.blockNo, err.Error())
		return "", err
	}
	if n != len(dat) {
		common.Logger.Errorf("write fs:%s-%d len error: %d != %d", tfs.prefix, tfs.blockNo, n, len(dat))
		return "", errors.New("len error")
	}

	blockRef := fmt.Sprintf("%d-%d-%d", tfs.blockNo, tfs.offset, n)
	tfs.offset += int64(n)
	return blockRef, nil
}

func (tfs *Fs) openFile() error {
	subDir := fmt.Sprintf("%s/%s", common.GlbBaInfa.Conf.Infra.FsDir, tfs.prefix)
	name := fmt.Sprintf("%s/%d", subDir, tfs.blockNo)
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		common.Logger.Errorf("open fs:%s failed:%s", file, err.Error())
		return err
	}

	info, err := file.Stat()
	if os.IsNotExist(err) {
		common.Logger.Errorf("Stat fs:%s failed:%s", file, err.Error())
		return err
	}

	//赋值
	tfs.offset = info.Size()
	tfs.file = file
	if info.Size() >= file_max_size {
		tfs.blockNo = tfs.blockNo + 1
		tfs.offset = 0
		file := fmt.Sprintf("%s/%d", subDir, tfs.blockNo)
		//关闭老的
		tfs.file.Close()
		tfs.file, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			common.Logger.Errorf("open fs:%s failed:%s", file, err.Error())
			return err
		}
	}

	return nil
}

func (tfs *Fs) load() {
	subDir := fmt.Sprintf("%s/%s", common.GlbBaInfa.Conf.Infra.FsDir, tfs.prefix)
	os.MkdirAll(subDir, 0755)
	curBlock := -1
	err := filepath.WalkDir(subDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		block, err := strconv.Atoi(d.Name())
		if err != nil {
			common.Logger.Warnf("Atoi %s failed:%s", d.Name(), err.Error())
			return nil
		}
		if block > curBlock {
			curBlock = block
		}
		return nil
	})

	if err != nil {
		common.Logger.Errorf("WalkDir fs:%s failed:%s", tfs.prefix, err.Error())
		panic(err)
	}

	if curBlock < 0 {
		common.Logger.Infof("Load fs:%s is new", subDir)
		curBlock = 0
	}

	tfs.blockNo = curBlock
	err = tfs.openFile()
	if err != nil {
		panic(err)
	}
}
