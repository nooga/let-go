# typed: false
# frozen_string_literal: true

class LetGo < Formula
  desc "A Clojure dialect implemented as a bytecode VM in Go"
  homepage "https://github.com/nooga/let-go"
  license "MIT"
  version "1.3.4"

  on_macos do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.4/let-go_1.3.4_darwin_amd64.tar.gz"
      sha256 "822450a669880dbf2ce6455b79f1411fb2ffed32d71f1929177374dd4ba494d6"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.4/let-go_1.3.4_darwin_arm64.tar.gz"
      sha256 "b9168127a1aedce4d4157bf84160902aac628c85c12fc488aafb57f058408984"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.4/let-go_1.3.4_linux_amd64.tar.gz"
      sha256 "231095080d926201127828911adcec37b50676d135488433df3ae7db74778246"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.4/let-go_1.3.4_linux_arm64.tar.gz"
      sha256 "a8036eec0f34216322c6380ac8b136c86947210c7a15d77e1eabaf622b64b650"
    end
  end

  def install
    bin.install "lg"
  end

  test do
    assert_equal "2", shell_output("#{bin}/lg -e '(+ 1 1)'").strip
  end
end
