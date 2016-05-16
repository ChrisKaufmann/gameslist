echo "select 'user_games';select * from user_games" | mysql -uroot gameslist_test
echo "select 'user_consoles';select * from user_consoles" | mysql -uroot gameslist_test
echo "select 'consoles';select * from consoles" | mysql -uroot gameslist_test
echo "select 'users';select * from users" | mysql -uroot gameslist_test
echo "select 'games';select * from games limit 2" | mysql -uroot gameslist_test
