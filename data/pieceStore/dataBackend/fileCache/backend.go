package fileCacheDataBackend

import (
	"io"
	"os"

	"github.com/anacrolix/missinggo/filecache"

	"github.com/michaeldye/torrent/data/pieceStore/dataBackend"
)

type backend struct {
	c *filecache.Cache
}

func New(fc *filecache.Cache) *backend {
	return &backend{
		c: fc,
	}
}

var _ dataBackend.I = &backend{}

func (me *backend) Delete(path string) (err error) {
	err = me.c.Remove(path)
	return
}

func (me *backend) GetLength(path string) (ret int64, err error) {
	fi, err := me.c.StatFile(path)
	if os.IsNotExist(err) {
		err = dataBackend.ErrNotFound
	}
	if err != nil {
		return
	}
	ret = fi.Size()
	return
}

func (me *backend) Open(path string, flag int) (ret dataBackend.File, err error) {
	ret, err = me.c.OpenFile(path, flag)
	return
}

func (me *backend) OpenSection(path string, off, n int64) (ret io.ReadCloser, err error) {
	f, err := me.c.OpenFile(path, os.O_RDONLY)
	if os.IsNotExist(err) {
		err = dataBackend.ErrNotFound
	}
	if err != nil {
		return
	}
	ret = struct {
		io.Reader
		io.Closer
	}{
		io.NewSectionReader(f, off, n),
		f,
	}
	return
}
