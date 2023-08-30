FROM golang:1.21-alpine

WORKDIR /workspace

COPY custom/git/git-prompt.sh /custom/git/
COPY custom/git/git-completion.bash /custom/git/

RUN apk update \
    && apk add --no-cache bash git openssh curl wget musl-dev gcc jq \
    # .bashrcのカスタマイズ
    && { \
        # Git用Promptの設定
        echo 'source /custom/git/git-prompt.sh'; \
        echo 'source /custom/git/git-completion.bash'; \
        echo 'GIT_PS1_SHOWDIRTYSTATE=true'; \
        echo "export PS1='\[\033[32m\]\u@\h\[\033[00m\]:\[\033[34m\]\w\[\033[31m\]\$(__git_ps1)\[\033[00m\]\n\$ '"; \
        # aliasの定義
        echo "alias ll='ls -l'"; \
    } >> /root/.bashrc \
    # タイムゾーンの設定
    && ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    # golang hot reloadツールの導入
    && go install github.com/cosmtrek/air@latest \
    # VSCode golang プラグイン
    && go install github.com/cweill/gotests/gotests@v1.6.0 \
    && go install github.com/fatih/gomodifytags@v1.16.0 \
    && go install github.com/josharian/impl@v1.1.0 \
    && go install github.com/haya14busa/goplay/cmd/goplay@v1.0.0 \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest \
    && go install golang.org/x/tools/gopls@latest
ENTRYPOINT [ "bash" ]