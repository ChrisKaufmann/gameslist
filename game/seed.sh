#!/bin/bash
echo "delete from consoles" | mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test
echo "delete from games" | mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test
echo "delete from user_games" | mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test
echo "delete from user_consoles" | mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test

mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test < ../sql/seed_consoles.sql
mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test < ../sql/seed_2600_games.sql
mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test < ../sql/seed_nes_games.sql

echo "replace into games (id, name, console_name) values (1,'game1','NES'),(2,'game2','NES')" | mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test
echo "replace into users (id, email, admin) values (1,'demouser',false),(2,'adminuser',true)" | mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test
echo "replace into user_consoles (name,user_id,has, manual, box, rating,review) values ('NES',1,true,true,true,3,'is good')"| mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test
echo "replace into user_games (id,game_id,user_id,has, manual, box, rating,review) values (1,1,1,true,true,true,3,'is good')"| mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test
if [[ $? -ne 0 ]]
  then
	exit $?
fi
