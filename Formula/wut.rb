class Wut < Formula
  desc "Git worktree manager that keeps worktrees out of your repo"
  homepage "https://github.com/simonbs/wut"
  url "https://github.com/simonbs/wut/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "88a67f61a59bb0d7733ba2328b0b0a575a821330389171838a371f26dd6a7809"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/wut"
  end

  def caveats
    <<~EOS
      Add shell integration to ~/.zshrc or ~/.bashrc:
        eval "$(wut init)"
    EOS
  end

  test do
    assert_match "wut", shell_output("#{bin}/wut --help")
  end
end
