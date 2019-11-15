#Trie

## Ausführen mit Docker

-   Images bauen

    ```
    make docker
    ```

-   ein (Docker)-Netzwerk `actors` erzeugen

    ```
    docker network create actors
    ```

-   Starten des Tree-Services und binden an den Port 8090 des Containers mit dem DNS-Namen
    `treeservice` (entspricht dem Argument von `--name`) im Netzwerk `actors`:

    ```
    docker run --rm --net actors --name treeservice treeservice \
      --bind="treeservice.actors:8090"
    ```

    Damit das funktioniert, müssen Sie folgendes erst im Tree-Service implementieren:

    -   die `main` verarbeitet Kommandozeilenflags und
    -   der Remote-Actor nutzt den Wert des Flags
    -   wenn Sie einen anderen Port als `8090` benutzen wollen,
        müssen Sie das auch im Dockerfile ändern (`EXPOSE...`)

-   Starten des Tree-CLI, Binden an `treecli.actors:8091` und nutzen des Services unter
    dem Namen und Port `treeservice.actors:8090`:

    ```
    docker run --rm --net actors --name treecli treecli --bind="treecli.actors:8091" \
      --remote="treeservice.actors:8090" trees
    ```

    Hier sind wieder die beiden Flags `--bind` und `--remote` beliebig gewählt und
    in der Datei `treeservice/main.go` implementiert. `trees` ist ein weiteres
    Kommandozeilenargument, dass z.B. eine Liste aller Tree-Ids anzeigen soll.

    Zum Ausprobieren können Sie den Service dann laufen lassen. Das CLI soll ja jedes
    Mal nur einen Befehl abarbeiten und wird dann neu gestartet.

-   Zum Beenden, killen Sie einfach den Tree-Service-Container mit `Ctrl-C` und löschen
    Sie das Netzwerk mit

    ```
    docker network rm actors
    ```

## Ausführen mit Docker ohne vorher die Docker-Images zu bauen

Nach einem Commit baut der Jenkins, wenn alles durch gelaufen ist, die beiden
Docker-Images. Sie können diese dann mit `docker pull` herunter laden. Schauen Sie für die
genaue Bezeichnung in die Consolenausgabe des Jenkins-Jobs.

Wenn Sie die Imagenamen oben (`treeservice` und `treecli`) durch die Namen aus der
Registry ersetzen, können Sie Ihre Lösung mit den selben Kommandos wie oben beschrieben,
ausprobieren.

## CLI commands

-   Starten des Treeservice mit `go run ./treeservice`
-   Treeclient:
    -   Erstellen eines neuen Baums: `go run ./treecli -create trie 10`
    -   Einfügen eines Wertes: `go run ./treecli -insert -key 8 -value="hello" -id 545 -token 37`
    -   Löschen eines Wertes: `go run ./treecli -delete -key -id 545 token 37`
    -   Traversieren: `go run ./treecli -traverse -id 545 token 37`
    -   Suchen: `go run ./treecli -find -key 8 -id 545 token 37`
    -   Baum löschen (muss mit yes bestätigt werden): `go run ./treecli -deleteTree -id 545 token 37`
-  Stimmt der Token mit die ID nicht überein, wird eine Fehlermeldung ausgegeben.

## Flags CLI:
```
  -id int
        ID des Baums (default 0)
  -delete
        löscht einen key und value aus dem Baum (default -1)
  -delete-trie bool
        löscht den ganzen Baum (default false)
  -insert bool
        fügt einen key und value in den Baum ein (default false)
  -key int
        Key der für Insert/Search/Delete benötigt wird (default -1)
  -create-trie int
        legt einen neuen Baum mit der angegebenen Größe an und gibt die id und den token aus (default -1)
  -find int
        sucht den value zu einem key (default -1)
  -token string
        Token des Baums (default "")
  -traverse bool
        gibt den Inhalt des Baums aus (default false)
  -value string
        value der im Baum abgelegt wird/ist (default "")
