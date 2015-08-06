echo "Building gameslist"
sh -c 'go build -o gameslist auth.go db.go main.go util.go'

