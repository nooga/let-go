# typed: false
# frozen_string_literal: true

class LetGo < Formula
  desc "A Clojure dialect implemented as a bytecode VM in Go"
  homepage "https://github.com/nooga/let-go"
  license "MIT"
  version "1.3.6"

  on_macos do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.6/let-go_1.3.6_darwin_amd64.tar.gz"
      sha256 "a03b53a9a4744adb56898716ab5d90020850264e1a2779c083a43291fe2c18bf"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.6/let-go_1.3.6_darwin_arm64.tar.gz"
      sha256 "d852114011286b890fcc4229b4fe69ac2eed37256815496be96b12302ad38221"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.3.6/let-go_1.3.6_linux_amd64.tar.gz"
      sha256 "1d1325186bf27c6f4fcaa56c287efb13648db668277ec8b2f891926171739225"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.3.6/let-go_1.3.6_linux_arm64.tar.gz"
      sha256 "8bdaf8e368781be953fb593d5aa78118526e304f1ea2936b8a020963436e461f"
    end
  end

  def install
    bin.install "lg"
  end

  test do
    assert_equal "2", shell_output("#{bin}/lg -e '(+ 1 1)'").strip
  end
end
