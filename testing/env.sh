function clean() {
	docker-compose down -t 1
}
function start() {
	docker-compose up -d
}

function enter_goland() {
	docker-compose exec  golang  bash
}
