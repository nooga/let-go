# typed: false
# frozen_string_literal: true

class LetGo < Formula
  desc "A Clojure dialect implemented as a bytecode VM in Go"
  homepage "https://github.com/nooga/let-go"
  license "MIT"
  version "1.3.8"

  on_macos do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.8/let-go_1.3.8_darwin_amd64.tar.gz"
      sha256 "ee56fa157b470823a57c2d1a85e3128b41840cc3fd6f3f04560442101b0caefe"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.8/let-go_1.3.8_darwin_arm64.tar.gz"
      sha256 "3facd131a41f211ed08b232bedbb30ce98d18c0cc0f730db5d046d43627bd412"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.8/let-go_1.3.8_linux_amd64.tar.gz"
      sha256 "66694c0d7933091f64f74532f03ea1f93298064855d9d49a28679d4a1f066a66"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.8/let-go_1.3.8_linux_arm64.tar.gz"
      sha256 "818a93eedcedcf03c33fdbd9a5ca91f999bbe77826cd7aad0b17cc11b7f1c985"
    end
  end

  def install
    bin.install "lg"
  end

  test do
    assert_equal "2", shell_output("#{bin}/lg -e '(+ 1 1)'").strip
  end
end
