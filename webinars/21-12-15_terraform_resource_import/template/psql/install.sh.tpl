#!/bin/bash
echo "${SSH_PUB_KEY}" | sudo tee -a /etc/ssh/sdm_ca.pub
echo "TrustedUserCAKeys /etc/ssh/sdm_ca.pub" | sudo tee -a /etc/ssh/sshd_config
sudo systemctl restart ssh
sudo apt-get update && sudo apt-get upgrade -y
sudo apt-get install -y unzip curl wget
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt-get install -y  postgresql
echo "listen_addresses = '*'" | sudo tee -a /etc/postgresql/*/main/postgresql.conf
echo "host all all 10.0.2.0/16 md5" | sudo tee -a /etc/postgresql/*/main/pg_hba.conf
sudo systemctl start postgresql && sudo systemctl enable postgresql
curl -O https://www.postgresqltutorial.com/wp-content/uploads/2019/05/dvdrental.zip
unzip ./dvd*.zip && rm -rf dvdrental.zip
sudo -u postgres createdb dvdrental
sudo -u postgres pg_restore --no-owner --role=postgres -U postgres --dbname dvdrental --verbose dvdrental.tar
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'notastrongpassword123';"
sudo systemctl restart postgresql