#!/bin/sh
echo "🔧 Setting root user auth method..."

mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" <<EOF
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '${MYSQL_ROOT_PASSWORD}';
FLUSH PRIVILEGES;
EOF
