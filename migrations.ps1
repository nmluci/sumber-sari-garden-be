$cmd = $args[0]
$params = $args[1]
$db_addr = "mysql://root:@tcp(localhost:3306)/main_db?parseTime=true&multiStatements=true"

if (($cmd -eq "up") -or ($cmd -eq "down" )) {
    if ($params -ne "") {
        migrate -source file://./db/migrations -database $db_addr $cmd $params
    } else {
        migrate -source file://./db/migrations -database $db_addr $cmd
    }
} elseif ($cmd -eq "drop") {
    migrate -source file://./db/migrations -database $db_addr $cmd
} elseif ($cmd -eq "new") {
    migrate create -ext sql -dir ./db/migrations $params
}