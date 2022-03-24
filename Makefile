#! bash

# run the dev server
run_api:
	docker-compose -f server/compose.yaml -p projetfedere up $(args)

# run_ui:
# 	# RUNNING APP UI...
# 	@cd appui && npm run dev
