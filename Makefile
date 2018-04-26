default:
	go run main.go init_logic_index.go init_logic_poll_ticker.go \
		init_logic_rules.go init_logic_settings.go init_resp.go

test:
	go test ./...
