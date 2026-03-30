# typed: false
# frozen_string_literal: true

class LetGo < Formula
  desc "A Clojure dialect implemented as a bytecode VM in Go"
  homepage "https://github.com/nooga/let-go"
  license "MIT"
  version "1.3.1"

  on_macos do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.1/let-go_1.3.1_darwin_amd64.tar.gz"
      sha256 "02e79a7f3167f978e1e4bff923699120a5616aa4a6ba80eb54e76a9c9bf153a8"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.1/let-go_1.3.1_darwin_arm64.tar.gz"
      sha256 "ccc54e23535d6bde9ccab690f2f18ddfcd31ed9247732ee7460d83af17515c4c"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.1/let-go_1.3.1_linux_amd64.tar.gz"
      sha256 "5fd2aabe41ab1886f62ffce5729dcebd7d5799913a75b71f54b0879b7004a259"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.1/let-go_1.3.1_linux_arm64.tar.gz"
      sha256 "a872b953de4d9fd07c773279408f374d202686735b7485c0d210c0ccf836a7d1"
    end
  end

  def install
    bin.install "lg"
  end

  test do
    assert_equal "2", shell_output("#{bin}/lg -e '(+ 1 1)'").strip
  end
end
