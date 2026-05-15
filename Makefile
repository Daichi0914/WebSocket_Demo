.PHONY: up down build rebuild logs ps restart clean e2e-up e2e-down test-backend test-integration test-frontend test-e2e stg-up stg-down stg-build prod-up prod-down prod-build

# --- Development Environment ---

up:
	touch ./.htpasswd
	podman volume create chat-dev_chat_mysql_dev_data || true
	podman volume create chat-dev_chat_redis_dev_data || true
	podman compose up -d --scale cloudflared=0
	rm -f ./.htpasswd

down:
	podman compose down

build:
	podman compose build

rebuild:
	touch ./.htpasswd
	podman volume create chat-dev_chat_mysql_dev_data || true
	podman volume create chat-dev_chat_redis_dev_data || true
	podman compose up -d --build --scale cloudflared=0
	rm -f ./.htpasswd

logs:
	podman compose logs -f

ps:
	podman compose ps

restart:
	podman compose restart

clean:
	podman compose down -v

# --- Testing ---

test-backend:
	cd backend && go test -v ./...

test-integration:
	cd backend && go test -v -tags=integration ./...

test-frontend:
	cd frontend && npm run test

test-e2e:
	$(MAKE) e2e-up
	@echo "Waiting for E2E backend to be ready..."
	@for i in $$(seq 1 30); do \
		if curl -s http://localhost:8081/api/messages > /dev/null 2>&1; then \
			echo "E2E backend is ready!"; \
			break; \
		fi; \
		echo "Waiting... ($$i/30)"; \
		sleep 2; \
		if [ $$i -eq 30 ]; then echo "Timeout waiting for E2E backend"; exit 1; fi; \
	done
	@cd frontend && E2E_ENV=true NEXT_PUBLIC_API_URL=http://localhost:8081 NEXT_PUBLIC_WS_URL=ws://localhost:8081/ws npm run test:e2e; \
	EXIT_CODE=$$?; $(MAKE) -C .. e2e-down; exit $$EXIT_CODE

# --- E2E Environment ---

e2e-up:
	podman compose -p websocket-chat-e2e --env-file .env.e2e --profile e2e up -d

e2e-down:
	podman compose -p websocket-chat-e2e --env-file .env.e2e --profile e2e down -v

