package mp4

import "io"

// Media Box (mdia - mandatory)
//
// Contained in : Track Box (trak)
//
// Status: decoded
//
// Contains all information about the media data.
type MdiaBox struct {
	Mdhd *MdhdBox
	Hdlr *HdlrBox
	Minf *MinfBox
}

func DecodeMdia(r io.Reader, size uint64) (Box, error) {
	l, err := DecodeContainer(r, size)
	if err != nil {
		return nil, err
	}
	m := &MdiaBox{}
	for _, b := range l {
		switch b.Type() {
		case "mdhd":
			m.Mdhd = b.(*MdhdBox)
		case "hdlr":
			m.Hdlr = b.(*HdlrBox)
		case "minf":
			m.Minf = b.(*MinfBox)
		default:
			return nil, &BadFormatErr{
				enclosingBox:  "mdia",
				unexpectedBox: b.Type(),
			}
		}
	}
	return m, nil
}

func (b *MdiaBox) Type() string {
	return "mdia"
}

func (b *MdiaBox) Size() uint64 {
	sz := AddHeaderSize(b.Mdhd.Size())
	if b.Hdlr != nil {
		sz += AddHeaderSize(b.Hdlr.Size())
	}
	if b.Minf != nil {
		sz += AddHeaderSize(b.Minf.Size())
	}
	return sz
}

func (b *MdiaBox) Dump() {
	b.Mdhd.Dump()
	if b.Minf != nil {
		b.Minf.Dump()
	}
}

func (b *MdiaBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	err = b.Mdhd.Encode(w)
	if err != nil {
		return err
	}
	if b.Hdlr != nil {
		err = b.Hdlr.Encode(w)
		if err != nil {
			return err
		}
	}
	return b.Minf.Encode(w)
}
