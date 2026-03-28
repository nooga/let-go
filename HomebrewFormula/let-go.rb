# typed: false
# frozen_string_literal: true

class LetGo < Formula
  desc "A Clojure dialect implemented as a bytecode VM in Go"
  homepage "https://github.com/nooga/let-go"
  license "MIT"
  version "1.0.0"

  on_macos do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.0.0/let-go_1.0.0_darwin_amd64.tar.gz"
      sha256 "41a00e145b2beb4a17953bc2c78b58dd870006b0785fd54105a80a08a63a4e17"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.0.0/let-go_1.0.0_darwin_arm64.tar.gz"
      sha256 "5e75ab343b879749caa9689ae7078fa5d65ad2f790a1424d6100e25379b1018c"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.0.0/let-go_1.0.0_linux_amd64.tar.gz"
      sha256 "d9fc784f8b5f481962aa0dd5552fb1c6d395058a37b47163c473424bb8373738"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.0.0/let-go_1.0.0_linux_arm64.tar.gz"
      sha256 "c721090ad96b12f5c09f8526396903623df6460b8dd9963e5b45cc7185529553"
    end
  end

  def install
    bin.install "lg"
  end

  test do
    assert_equal "2", shell_output("#{bin}/lg -e '(+ 1 1)'").strip
  end
end
