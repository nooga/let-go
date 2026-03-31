# typed: false
# frozen_string_literal: true

class LetGo < Formula
  desc "A Clojure dialect implemented as a bytecode VM in Go"
  homepage "https://github.com/nooga/let-go"
  license "MIT"
  version "1.3.3"

  on_macos do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.3/let-go_1.3.3_darwin_amd64.tar.gz"
      sha256 "dd56318ffbc47d0954308d92adc6bb5f5b01235ba9b852961d89728c59c249fb"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.3/let-go_1.3.3_darwin_arm64.tar.gz"
      sha256 "a0b52fd8ab71e21b6f398ca75592ce2a25119efc9064840c75b8a7fdcfcf7ba4"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.3/let-go_1.3.3_linux_amd64.tar.gz"
      sha256 "20032f1b3476cca52723ab94800d1ec6e3db1f013933e2bc133c6883cda9d485"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.3/let-go_1.3.3_linux_arm64.tar.gz"
      sha256 "c3b30c14dfb86c523d4f11449d251556c67485f86800560dd032c324ddbfdbf8"
    end
  end

  def install
    bin.install "lg"
  end

  test do
    assert_equal "2", shell_output("#{bin}/lg -e '(+ 1 1)'").strip
  end
end
