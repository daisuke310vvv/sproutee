class Sproutee < Formula
  desc "A powerful CLI tool for efficient Git worktree management"
  homepage "https://github.com/daisuke310vvv/sproutee"
  url "https://github.com/daisuke310vvv/sproutee/archive/v0.1.0.tar.gz"
  sha256 "" # Will be updated by GoReleaser
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/sproutee"
  end

  test do
    system "#{bin}/sproutee", "--help"
    assert_match "Sproutee", shell_output("#{bin}/sproutee --help")
  end
end