class Sentinel < Formula
  desc "SentinelStacks - AI Agent Management Platform"
  homepage "https://github.com/satishgonella2024/sentinelstacks"
  version "1.0.0"

  if Hardware::CPU.arm?
    url "https://github.com/satishgonella2024/sentinelstacks/releases/download/v1.0.0/sentinel-darwin-arm64.tar.gz"
    sha256 "9e938d9bad87501c4fea8123b55c7b65b70a1bce3503a5038e67ee9dc4b36145"
  else
    url "https://github.com/satishgonella2024/sentinelstacks/releases/download/v1.0.0/sentinel-darwin-amd64.tar.gz"
    sha256 "4a0ebc7dcfbc513fc0854291981eb9cd90d979d02d4e354249bfc28375df1e2c"
  end

  def install
    bin.install "sentinel"
  end

  test do
    system "#{bin}/sentinel", "--version"
  end
end 