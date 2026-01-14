class Termidoro < Formula
  desc "A terminal-based Pomodoro timer"
  homepage "https://github.com/user/termidoro"
  url "https://github.com/user/termidoro/archive/v1.0.0.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "make", "build"
    bin.install "termidoro"
    man1.install "termidoro.1"
  end

  test do
    system "#{bin}/termidoro", "--help"
  end
end
