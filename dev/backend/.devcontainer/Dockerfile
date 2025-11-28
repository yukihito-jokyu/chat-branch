# ベースイメージ
FROM golang:1.25.4-bookworm

# パッケージ ---
RUN apt-get update \
    && apt-get install -y build-essential git curl zsh wget procps \
    && curl -1sLf "https://dl.cloudsmith.io/public/task/task/setup.deb.sh" | bash \
    && apt-get update \
    && apt-get install -y task \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# goimports
RUN go install golang.org/x/tools/cmd/goimports@latest

# goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh \
    | sh -s -- -b $(go env GOPATH)/bin v2.6.2

# ユーザー設定
ARG USERNAME=golang
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# ユーザー作成と権限設定
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && apt-get update \
    && apt-get install -y sudo \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# ワークスペース作成と権限設定
RUN mkdir -p /workspace && chown -R $USERNAME:$USERNAME /workspace

# ターミナルの環境を zsh に変更
RUN chsh -s $(which zsh) $USERNAME

USER $USERNAME
WORKDIR /home/$USERNAME

# oh-my-zsh インストール
RUN rm -rf ~/.oh-my-zsh && \
    sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended
# autosuggestions インストール
RUN git clone https://github.com/zsh-users/zsh-autosuggestions \
    ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions
# syntax-highlighting インストール
RUN git clone https://github.com/zsh-users/zsh-syntax-highlighting.git \
    ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting

RUN sed -i 's/plugins=(git)/plugins=(git zsh-autosuggestions zsh-syntax-highlighting)/' ~/.zshrc

WORKDIR /workspace
