default:
	go run main.go init_config.go init_constant.go init_exchange.go \
		init_logic_bot.go init_logic_exchange.go init_logic_index.go \
		init_logic_portfolio.go init_logic_setup.go init_poll_ticker.go \
		init_repository.go init_resp.go init_routes.go \
		init_socket_account.go init_account

test:
	go test ./...
