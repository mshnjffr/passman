class Passman < Formula
  desc "Beautiful terminal password manager built with Go and Bubble Tea"
  homepage "https://github.com/mshnjffr/passman"
  url "https://github.com/mshnjffr/passman/archive/refs/tags/v1.0.2.tar.gz"
  sha256 "4ba1fc2bb5e8a2c03115418778c566d7222d0ff0f9b2a2bb61057f308b6e5698"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-s -w", "-o", bin/"passman"
  end

  test do
    assert_match "passman", shell_output("#{bin}/passman --version")
  end
end
