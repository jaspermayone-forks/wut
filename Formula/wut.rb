class Wut < Formula
  desc "Git worktree manager that keeps worktrees out of your repo"
  homepage "https://github.com/simonbs/wut"
  head "https://github.com/simonbs/wut.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/wut"
  end

  test do
    assert_match "wut", shell_output("#{bin}/wut --help")
  end
end
