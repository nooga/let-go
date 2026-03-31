# typed: false
# frozen_string_literal: true

class LetGo < Formula
  desc "A Clojure dialect implemented as a bytecode VM in Go"
  homepage "https://github.com/nooga/let-go"
  license "MIT"
  version "1.3.5"

  on_macos do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.5/let-go_1.3.5_darwin_amd64.tar.gz"
      sha256 "e68e174d980779ca4ad6aafafc185d9ea9de0bd27b1f6ff4b6ecd6444b20c496"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.5/let-go_1.3.5_darwin_arm64.tar.gz"
      sha256 "aab5a67a594877f2bd3046af3e40f2170447181ecc317597b79af737fd3ca189"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.5/let-go_1.3.5_linux_amd64.tar.gz"
      sha256 "da5529826013bcf415a15287381aa9b75e91336302aa297ba91b3f8f0e66a906"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.5/let-go_1.3.5_linux_arm64.tar.gz"
      sha256 "cf1dba11de607d4b8e5fc167ba5786281e9c3a75e17808f6c507ef147a320b7f"
    end
  end

  def install
    bin.install "lg"
  end

  test do
    assert_equal "2", shell_output("#{bin}/lg -e '(+ 1 1)'").strip
  end
end
