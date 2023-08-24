FROM golang:1.19-alpine

WORKDIR /workspace

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
    && ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENTRYPOINT [ "bash" ]