#!/bin/bash

# インストールする拡張機能のリスト
extensions=(
    "esbenp.prettier-vscode"
    "oderwat.indent-rainbow"
    "streetsidesoftware.code-spell-checker"
    "mhutchie.git-graph"
    "editorconfig.editorconfig"
    "shardulm94.trailing-spaces"
    "ms-ceintl.vscode-language-pack-ja"
    "golang.go"
    "windmilleng.vscode-go-autotest"
    "romanpeshkov.go-test-outline"
)

echo "拡張機能のインストールを開始します..."

# リストをループして agy コマンドを実行
for ext in "${extensions[@]}"; do
    echo "Installing: $ext"

    # agyコマンドの実行
    agy --install-extension $ext

    if [ $? -eq 0 ]; then
        echo "✓ $ext のインストールに成功しました"
    else
        echo "✗ $ext のインストールに失敗しました"
    fi
    echo "----------------------------------------"
done

echo "すべての処理が完了しました。"
