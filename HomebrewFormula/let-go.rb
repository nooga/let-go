# typed: false
# frozen_string_literal: true

class LetGo < Formula
  desc "A Clojure dialect implemented as a bytecode VM in Go"
  homepage "https://github.com/nooga/let-go"
  license "MIT"
  version "1.2.0"

  on_macos do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.2.0/let-go_1.2.0_darwin_amd64.tar.gz"
      sha256 "1cef42ab364f32e680f8766d04018cfb3c27e5188da3440da3125bcbbbfcd490"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.2.0/let-go_1.2.0_darwin_arm64.tar.gz"
      sha256 "0dd61bc58bd0ca5b0606d95c3746c5f542b33db3c1fab3bcd26fecd42715c4ed"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/nooga/let-go/releases/download/v1.2.0/let-go_1.2.0_linux_amd64.tar.gz"
      sha256 "d2070f6a085ba0f978c4811c84e77db14a51a929855b5be56751f964c1fdd3ea"
    end
    on_arm do
      url "https://github.com/nooga/let-go/releases/download/v1.2.0/let-go_1.2.0_linux_arm64.tar.gz"
      sha256 "fc22ec2525cd969682821eb9f7a0718cfad6c2d5942c6a7720cb65260f33554c"
    end
  end

  def install
    bin.install "lg"
  end

  test do
    assert_equal "2", shell_output("#{bin}/lg -e '(+ 1 1)'").strip
  end
end
