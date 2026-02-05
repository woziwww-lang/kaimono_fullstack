#!/bin/bash

echo "🚀 价格比較アプリ - 环境安装脚本 (macOS)"
echo "================================================"
echo ""

# 检查 Homebrew
if ! command -v brew &> /dev/null; then
    echo "❌ Homebrew 未安装"
    echo "📦 正在安装 Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
else
    echo "✅ Homebrew 已安装"
fi

echo ""
echo "📦 安装必需工具..."
echo ""

# 安装 Go
if ! command -v go &> /dev/null; then
    echo "🔧 安装 Go..."
    brew install go
else
    echo "✅ Go 已安装: $(go version)"
fi

# 安装 Docker
if ! command -v docker &> /dev/null; then
    echo "🐳 安装 Docker Desktop..."
    brew install --cask docker
    echo "⚠️  请启动 Docker Desktop 应用程序"
else
    echo "✅ Docker 已安装: $(docker --version)"
fi

# 安装 Flutter (可选)
read -p "是否安装 Flutter（移动端开发）? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if ! command -v flutter &> /dev/null; then
        echo "📱 安装 Flutter..."
        brew install --cask flutter
    else
        echo "✅ Flutter 已安装: $(flutter --version | head -1)"
    fi
fi

# 安装 Terraform (可选)
read -p "是否安装 Terraform（云部署）? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if ! command -v terraform &> /dev/null; then
        echo "☁️  安装 Terraform..."
        brew install terraform
    else
        echo "✅ Terraform 已安装: $(terraform --version | head -1)"
    fi
fi

echo ""
echo "================================================"
echo "✅ 安装完成！"
echo ""
echo "📋 下一步:"
echo "  1. 启动 Docker Desktop 应用"
echo "  2. 运行: make install"
echo "  3. 运行: make docker-up"
echo "  4. 运行: make server (在新终端)"
echo "  5. 运行: make web (在新终端)"
echo ""
echo "🔍 验证安装:"
echo "  node --version"
echo "  pnpm --version"
echo "  go version"
echo "  docker --version"
echo ""
