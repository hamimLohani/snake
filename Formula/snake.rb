class Snake < Formula
  desc "Modern developer-grade terminal Snake game with replays and boss key"
  homepage "https://github.com/hamimlohani/snake"
  url "https://github.com/hamimlohani/snake/archive/refs/tags/v1.0.7.tar.gz"
  sha256 "4c09834bf8f212d0e1f5084e51d3a4f0d234c35c1d44d9b2692a871c1a53e580"
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
