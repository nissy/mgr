NAME = mgr

redis-up:
	docker run -d \
		-p 6379:6379 \
		--name=$(NAME)-redis redis:4

redis-down:
	docker rm -f $(NAME)-redis
