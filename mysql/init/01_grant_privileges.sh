#!/bin/bash
set -e

# 環境変数を使って、Podmanのサブネットからの接続を許可する
# MYSQL_PWD は mysql コマンドがパスワードとして使用する特殊な環境変数です
export MYSQL_PWD=$MYSQL_ROOT_PASSWORD

mysql -u root <<EOSQL
  CREATE USER IF NOT EXISTS '$MYSQL_USER'@'10.89.0.%' IDENTIFIED BY '$MYSQL_PASSWORD';
  GRANT ALL PRIVILEGES ON $MYSQL_DATABASE.* TO '$MYSQL_USER'@'10.89.0.%';
  FLUSH PRIVILEGES;
EOSQL
