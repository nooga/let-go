# typed: false
# frozen_string_literal: true

class LetGo < Formula
  desc "A Clojure dialect implemented as a bytecode VM in Go"
  homepage "https://github.com/nooga/let-go"
  license "MIT"
  version "1.3.2"

  on_macos do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.2/let-go_1.3.2_darwin_amd64.tar.gz"
      sha256 "9f7fcb261bd4400994194d76d98590aecfaccefd4738e4be04178bd6142eb080"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.2/let-go_1.3.2_darwin_arm64.tar.gz"
      sha256 "6818c9227d8aa0175e12db275106930cb0918cfb8f9f8dde5a6958cc13cc851b"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.2/let-go_1.3.2_linux_amd64.tar.gz"
      sha256 "49c79160c9eb6d89fa381b034398d1edbcbc4c27877c630c9ad6c0b1ba25f94d"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.2/let-go_1.3.2_linux_arm64.tar.gz"
      sha256 "e6894b54e42f896f19c7db12c38275cfd922e8d13484b0fd4716f98e81e1ab4f"
    end
  end

  def install
    bin.install "lg"
  end

  test do
    assert_equal "2", shell_output("#{bin}/lg -e '(+ 1 1)'").strip
  end
end
