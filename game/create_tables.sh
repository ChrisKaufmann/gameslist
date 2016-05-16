#!/bin/bash
echo "drop table if exists consoles" | mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test
echo "drop table if exists games" | mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test
echo "drop table if exists users" | mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test
echo "drop table if exists sessions" | mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test
echo "drop table if exists user_games" | mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test
echo "drop table if exists user_consoles" | mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test
mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test < ../sql/consoles.sql
mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test < ../sql/games.sql
mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test < ../sql/sessions.sql
mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test < ../sql/user_games.sql
mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test < ../sql/users.sql
mysql -u gameslist_test -pKoZMEEh49nMts2T4XgEWVWC1 gameslist_test < ../sql/user_consoles.sql
if [[ $? -ne 0 ]]
  then
	exit $?
fi
