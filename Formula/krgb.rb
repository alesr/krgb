class Krgb < Formula
  desc "TUI for controlling Keychron K-series keyboard LED colors via raw HID"
  homepage "https://github.com/alesr/krgb"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/alesr/krgb/releases/download/v0.1.0/krgb_v0.1.0_darwin_arm64.tar.gz"
      sha256 ""
    else
      url "https://github.com/alesr/krgb/releases/download/v0.1.0/krgb_v0.1.0_darwin_amd64.tar.gz"
      sha256 ""
    end
  end

  def install
    bin.install "krgb_darwin_#{Hardware::CPU.arch}" => "krgb"
  end

  test do
    assert_predicate bin/"krgb", :exist?
  end
end
