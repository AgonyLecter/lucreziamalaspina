# Compilare l'eseguibile per Windows

Serve Go (1.22 o più recente). Da questa cartella:

```bash
python3 make_game.py game            # copia il gioco (con i caratteri inclusi) in ./game
go run github.com/tc-hib/go-winres@latest make --arch amd64   # icona e informazioni dell'exe
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-H windowsgui -s -w" -o LucreziaMalaspina.exe .
```

Per la versione browser di itch.io: `python3 make_game.py web sw`, poi comprimere il contenuto di `web/`.
