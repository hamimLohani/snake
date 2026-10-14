class Snake < Formula
  desc "Modern developer-grade terminal Snake game with replays and boss key"
  homepage "https://github.com/hamimlohani/snake"
  url "https://github.com/hamimlohani/snake/archive/refs/tags/v1.0.3.tar.gz"
  sha256 "bc78bb21db83ab8718efb25239bcdf1e1f1afdeabf4f8dd929c9f16ebef008da"
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
