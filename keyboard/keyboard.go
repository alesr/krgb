package keyboard

import (
	"strings"

	"github.com/sstallion/go-hid"

	"github.com/alesr/krgb/xerrors"
)

const (
	KeychronVID = 0x3434
	UsagePage   = 0xFF60
	UsageID     = 0x61
)

var ErrNotFound = xerrors.New(xerrors.DeviceNotFound, "Keychron K-series raw HID interface not found")

type Connection struct {
	Device *hid.Device
	Info   *hid.DeviceInfo
}

func Find() (*Connection, error) {
	var info *hid.DeviceInfo

	hid.Enumerate(KeychronVID, hid.ProductIDAny, func(i *hid.DeviceInfo) error {
		if i.UsagePage == UsagePage && i.Usage == UsageID &&
			strings.HasPrefix(i.ProductStr, "Keychron K") {
			info = i
		}
		return nil
	})

	if info == nil {
		return nil, ErrNotFound
	}

	dev, err := hid.OpenPath(info.Path)
	if err != nil {
		return nil, xerrors.Wrap(xerrors.DeviceOpen, "hid open failed", err)
	}
	return &Connection{Device: dev, Info: info}, nil
}
