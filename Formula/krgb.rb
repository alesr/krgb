class Krgb < Formula
  desc "TUI for controlling Keychron K-series keyboard LED colors via raw HID"
  homepage "https://github.com/alesr/krgb"
  version "0.1.1"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/alesr/krgb/releases/download/v0.1.1/krgb_v0.1.1_darwin_arm64.tar.gz"
      sha256 "fa7b01429c3f75d88d63731281b66d3c8e7c384830e77376d2b2aa1c4bcbeeda"
    else
      url "https://github.com/alesr/krgb/releases/download/v0.1.1/krgb_v0.1.1_darwin_amd64.tar.gz"
      sha256 "e7f74dd2c43052cf3f3c71d9d3f389e5d62c3898716eb6a5731571c9f2e3927a"
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
