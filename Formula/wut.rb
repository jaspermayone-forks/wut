class Wut < Formula
  desc "Git worktree manager that keeps worktrees out of your repo"
  homepage "https://github.com/simonbs/wut"
  url "https://github.com/simonbs/wut/archive/refs/tags/v0.3.2.tar.gz"
  sha256 "b9ed3814895d7043eb0e95a736ce8c36c527ee50b8ccf9904c8dcd551d1ce281"
  head "https://github.com/simonbs/wut.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/wut"
  end

  test do
    assert_match "wut", shell_output("#{bin}/wut --help")
  end
end
