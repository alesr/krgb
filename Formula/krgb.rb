class Krgb < Formula
  desc "TUI for controlling Keychron K-series keyboard LED colors via raw HID"
  homepage "https://github.com/alesr/krgb"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/alesr/krgb/releases/download/v0.1.0/krgb_v0.1.0_darwin_arm64.tar.gz"
      sha256 "592fc0e4511f360b3dd11dbd784daa7782159ad52519ecb14be32026bd9f7655"
    else
      url "https://github.com/alesr/krgb/releases/download/v0.1.0/krgb_v0.1.0_darwin_amd64.tar.gz"
      sha256 "3f1e08535011012d98ce6b5d65c08c6e3bb071a7ca4fe2396a0868af164a38d8"
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
