# krgb

Go TUI for controlling Keychron K-series keyboard LED colors via raw HID.

## Installation

### Homebrew

    brew install alesr/tap/krgb

### Build from source

    go install github.com/alesr/krgb@latest

## Usage

    krgb

Requires a connected Keychron K-series keyboard with the raw HID interface.

## Requirements

- Keychron K-series keyboard with raw HID interface
- Go 1.26.4 (to build from source)

## Platform Support

**macOS only.** Tested on Sequoia with Keychron K13 Pro.

Should work with other K-series models sharing the same raw HID interface, but only K13 Pro has been verified.

## License

MIT
