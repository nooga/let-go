# typed: false
# frozen_string_literal: true

class LetGo < Formula
  desc "A Clojure dialect implemented as a bytecode VM in Go"
  homepage "https://github.com/nooga/let-go"
  license "MIT"
  version "1.3.0"

  on_macos do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.0/let-go_1.3.0_darwin_amd64.tar.gz"
      sha256 "ac6bde8e79a0e8fa9156eddc4a8035a128e21fdd47fb254941b824339ca71a54"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.0/let-go_1.3.0_darwin_arm64.tar.gz"
      sha256 "78f838a973df84e1f831c3d4908ef0f1c3064637b4b9bfb8cbd59a456d458310"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.0/let-go_1.3.0_linux_amd64.tar.gz"
      sha256 "e85e0031d064f03030bc5097f15bdf8550e04de65aa15fcc6c25288baa6cd170"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.0/let-go_1.3.0_linux_arm64.tar.gz"
      sha256 "612935d84084c70f3ac51a23a4ec198ced34d9a031bcfa966149bfdb25ca4875"
    end
  end

  def install
    bin.install "lg"
  end

  test do
    assert_equal "2", shell_output("#{bin}/lg -e '(+ 1 1)'").strip
  end
end
