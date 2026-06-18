class Krgb < Formula
  desc "TUI for controlling Keychron K-series keyboard LED colors via raw HID"
  homepage "https://github.com/alesr/krgb"
  version "0.1.0-beta.3"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/alesr/krgb/releases/download/v0.1.0-beta.3/krgb_v0.1.0-beta.3_darwin_arm64.tar.gz"
      sha256 "c46ab40571c4ff2855927d12d37d44a38365e6e2082b3f0b0de2b5066cff3999"
    else
      url "https://github.com/alesr/krgb/releases/download/v0.1.0-beta.3/krgb_v0.1.0-beta.3_darwin_amd64.tar.gz"
      sha256 "4d2f4064318578f234e7e0030fd58227c6eebe57a70a2d8b0454d0259f61deda"
    end
  end

  license "MIT"

  def install
    bin.install "krgb_darwin_#{Hardware::CPU.arch}" => "krgb"
  end

  test do
    assert_predicate bin/"krgb", :exist?
  end
end
