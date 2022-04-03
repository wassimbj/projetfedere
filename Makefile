#! bash

# run the dev server
up:
	docker-compose -f server/compose.yaml -p projetfedere up $(args)
