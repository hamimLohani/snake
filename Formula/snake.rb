class Snake < Formula
  desc "Modern developer-grade terminal Snake game with replays and boss key"
  homepage "https://github.com/hamimlohani/snake"
  url "https://github.com/hamimlohani/snake/archive/refs/tags/v1.0.7.tar.gz"
  sha256 "979bfa9dd770e7933a87bca3c1bf9b0e1d7a6bfeb0a07162c4a20fe63738bb9b"
  license "MIT"
  head "https://github.com/hamimlohani/snake.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version}"), "./cmd/snake"
    man6.install "man/snake.6"
    man1.install "man/snake.1"
  end

  test do
    assert_match "snake version", shell_output("#{bin}/snake --version")
    assert_match "Headless AI Bot Simulation Complete", shell_output("#{bin}/snake bot --headless --ticks 20")
  end
end
