# typed: false
# frozen_string_literal: true

class LetGo < Formula
  desc "A Clojure dialect implemented as a bytecode VM in Go"
  homepage "https://github.com/nooga/let-go"
  license "MIT"
  version "1.1.0"

  on_macos do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.1.0/let-go_1.1.0_darwin_amd64.tar.gz"
      sha256 "83b71844526c8d4d551d6bce253e877592464300607dbb1b27ed36cad09ff0cc"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.1.0/let-go_1.1.0_darwin_arm64.tar.gz"
      sha256 "bd2c8338e8905218ec989fe13f36a9d054fb738f2206c7bdc61dcf903dc9336a"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.1.0/let-go_1.1.0_linux_amd64.tar.gz"
      sha256 "12cd5fe2321c466ea4d099d3db1d4c16dfeb2e2eba87f98fd6ce93aad8e14a46"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.1.0/let-go_1.1.0_linux_arm64.tar.gz"
      sha256 "51e56255f9a5bcd970e2deeda89a44cc34d6efddb04e5c2ba31020066121d3cb"
    end
  end

  def install
    bin.install "lg"
  end

  test do
    assert_equal "2", shell_output("#{bin}/lg -e '(+ 1 1)'").strip
  end
end
