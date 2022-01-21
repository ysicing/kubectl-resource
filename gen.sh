#!/bin/bash

version=$(cat version.txt)
macosamd64sha=$(cat dist/kr_darwin_amd64.sha256sum | awk '{print $1}')
macosarm64sha=$(cat dist/kr_darwin_arm64.sha256sum | awk '{print $1}')
linuxamd64sha=$(cat dist/kr_linux_amd64.sha256sum | awk '{print $1}')
linuxarm64sha=$(cat dist/kr_linux_arm64.sha256sum | awk '{print $1}')


cat > kr.rb <<EOF
class KR < Formula
    desc "Devops tools kube-resource"
    homepage "https://github.com/ysicing/kube-resource"
    version "${version}"

    if OS.mac?
      if Hardware::CPU.arm?
        url "https://github.com/ysicing/kube-resource/releases/download/#{version}/kr_darwin_arm64"
        sha256 "${macosarm64sha}"
      else
        url "https://github.com/ysicing/kube-resource/releases/download/#{version}/kr_darwin_amd64"
        sha256 "${macosamd64sha}"
      end  
    elsif OS.linux?
      if Hardware::CPU.intel?
        url "https://github.com/ysicing/kube-resource/releases/download/#{version}/kr_linux_amd64"
        sha256 "${linuxamd64sha}"
      end
      if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
        url "https://github.com/ysicing/kube-resource/releases/download/#{version}/kr_linux_arm64"
        sha256 "${linuxarm64sha}"
      end
    end

    def install
      if OS.mac?
        if Hardware::CPU.intel?
          bin.install "kr_darwin_amd64" => "kr"
        else
          bin.install "kr_darwin_arm64" => "kr"
        end 
      elsif OS.linux?
        if Hardware::CPU.intel?
          bin.install "kr_linux_amd64" => "kr"
        else
          bin.install "kr_linux_arm64" => "kr"
        end 
      end
    end

    test do
      assert_match "kr vervion v#{version}", shell_output("#{bin}/kr version")
    end
end
EOF

docker build -t ysicing/taprb:kr -f hack/brew/Dockerfile .
docker push ysicing/taprb:kr