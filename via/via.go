package via

import (
	"github.com/sstallion/go-hid"

	"github.com/alesr/krgb/xerrors"
)

const (
	CmdSetValue   = 0x07
	CmdSaveEEPROM = 0x09

	ChanRGBMatrix = 0x03

	ValBrightness = 1
	ValEffect     = 2
	ValSpeed      = 3
	ValColor      = 4
)

func writeCmd(dev *hid.Device, cmd, valueID byte, payload []byte) error {
	buf := make([]byte, 32)
	buf[1] = cmd
	buf[2] = ChanRGBMatrix
	buf[3] = valueID
	copy(buf[4:], payload)

	if _, err := dev.Write(buf); err != nil {
		return xerrors.Wrap(xerrors.DeviceWrite, "via write failed", err)
	}
	return nil
}

func SetValue(dev *hid.Device, valueID byte, payload []byte) error {
	return writeCmd(dev, CmdSetValue, valueID, payload)
}

func SaveValue(dev *hid.Device, valueID byte, payload []byte) error {
	return writeCmd(dev, CmdSaveEEPROM, valueID, payload)
}
