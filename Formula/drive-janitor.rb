class DriveJanitor < Formula
  desc "Clean up tools for managing and cleaning up your drive"
  homepage "https://github.com/secatscale/drive-janitor"
  url "https://github.com/secatscale/drive-janitor/archive/refs/tags/v0.0.22.tar.gz"
  sha256 "0019dfc4b32d63c1392aa264aed2253c1e0c2fb09216f8e2cc269bbfb8bb49b5"
  license "GPL-3.0"

  depends_on "go" => :build
  depends_on "pkg-config" => :build
  depends_on "yara" => :build

  def install
    system "go", "build"
  end

  test do
    system "#{bin}/drive-janitor", "-h"
  end
end
