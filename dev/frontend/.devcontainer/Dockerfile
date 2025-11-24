# ベースイメージ
FROM node:24-slim

# パッケージ ---
RUN apt-get update \
    && apt-get install -y build-essential git curl zsh wget procps \
    && curl -1sLf "https://dl.cloudsmith.io/public/task/task/setup.deb.sh" | bash \
    && apt-get update \
    && apt-get install -y task \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

ARG USERNAME=node

RUN mkdir -p /workspace && chown -R $USERNAME:$USERNAME /workspace

# ターミナルの環境を zsh に変更
RUN chsh -s $(which zsh) node

USER $USERNAME
WORKDIR /home/node

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

ENV PNPM_HOME="/home/$USERNAME/.local/share/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

# pnpm インストール
RUN export SHELL=/bin/bash \
    && wget -qO- https://get.pnpm.io/install.sh | sh -

WORKDIR /workspace
