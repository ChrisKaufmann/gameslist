echo "Building gameslist"
sh -c 'go build -o gameslist db.go auth.go console.go game.go main.go util.go'

