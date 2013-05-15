package fetion

import (
    "encoding/gob"
    "io"
)

func (f *Fetion) SaveCache(w io.Writer) (err error) {
    encoder := gob.NewEncoder(w)
    err = encoder.Encode(f.friends)
    return err
}

func (f *Fetion) LoadCache(r io.Reader) (err error) {
    decoder := gob.NewDecoder(r)
    err = decoder.Decode(&f.friends)
    return err
}
