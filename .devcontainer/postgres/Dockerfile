FROM postgres:14

COPY ./initdb /docker-entrypoint-initdb.d/
RUN localedef -i ja_JP -c -f UTF-8 -A /usr/share/locale/locale.alias ja_JP.UTF-8