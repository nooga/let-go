# typed: false
# frozen_string_literal: true

class LetGo < Formula
  desc "A Clojure dialect implemented as a bytecode VM in Go"
  homepage "https://github.com/nooga/let-go"
  license "MIT"
  version "1.3.7"

  on_macos do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.7/let-go_1.3.7_darwin_amd64.tar.gz"
      sha256 "2b3f756e1914c8a326a4d5568e64ad3b137b34a510f41370537d69f97af06653"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.7/let-go_1.3.7_darwin_arm64.tar.gz"
      sha256 "15c251f265a911198e908e6a7821b289d4a8b8753932842a427a9c1dc6f8ce60"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.7/let-go_1.3.7_linux_amd64.tar.gz"
      sha256 "2c475145af01f1a87c3b4a3709af982a0fa83f6381657229f67c4c2ad0ff6eca"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.7/let-go_1.3.7_linux_arm64.tar.gz"
      sha256 "8b7cfcbef58bcb4d9713f7295aeb540897153f5826855399bbc5f94ce51faa45"
    end
  end

  def install
    bin.install "lg"
  end

  test do
    assert_equal "2", shell_output("#{bin}/lg -e '(+ 1 1)'").strip
  end
end
