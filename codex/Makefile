PYTHON ?= python3
PIP ?= $(PYTHON) -m pip
STORAGE ?= ./.demo/tasks.json
ARGS ?= list --status all
PIP_INSTALL_FLAGS ?=

.PHONY: help install uninstall run stop add list update done remove summary test build clean

help:
	@printf '%s\n' \
		'Available targets:' \
		'  make install                 Install the package into the current Python environment' \
		'  make uninstall               Uninstall the all-day package from the current Python environment' \
		'  make run ARGS="..."          Run the CLI against STORAGE (default: ./.demo/tasks.json)' \
		'  make add TITLE="..."         Add a task (optional: ADD_ARGS="--priority high --tag work")' \
		'  make list                    List tasks (optional: LIST_ARGS="--status all --date YYYY-MM-DD")' \
		'  make update TASK_ID=N        Update a task (required: UPDATE_ARGS="...")' \
		'  make done TASK_ID=N          Mark a task done' \
		'  make remove TASK_ID=N        Remove a task' \
		'  make summary                 Show a summary (optional: SUMMARY_ARGS="--date YYYY-MM-DD")' \
		'  make stop                    Clear the demo storage file under ./.demo/' \
		'  make test                    Run the unit test suite' \
		'  make build                   Build sdist and wheel artifacts' \
		'  make clean                   Remove local build/cache artifacts'

install:
	$(PIP) install $(PIP_INSTALL_FLAGS) .

uninstall:
	-$(PIP) uninstall -y all-day

run:
	@mkdir -p "$(dir $(STORAGE))"
	$(PYTHON) -m all_day --storage "$(STORAGE)" $(ARGS)

add:
	@if [ -z "$(TITLE)" ]; then echo "Usage: make add TITLE='Write report' [ADD_ARGS='--priority high --tag work']"; exit 1; fi
	@mkdir -p "$(dir $(STORAGE))"
	$(PYTHON) -m all_day --storage "$(STORAGE)" add "$(TITLE)" $(ADD_ARGS)

list:
	@mkdir -p "$(dir $(STORAGE))"
	$(PYTHON) -m all_day --storage "$(STORAGE)" list $(LIST_ARGS)

update:
	@if [ -z "$(TASK_ID)" ]; then echo "Usage: make update TASK_ID=1 UPDATE_ARGS='--title \"Final report\" --clear-tags'"; exit 1; fi
	@if [ -z '$(strip $(UPDATE_ARGS))' ]; then echo "Usage: make update TASK_ID=1 UPDATE_ARGS='--title \"Final report\" --clear-tags'"; exit 1; fi
	@mkdir -p "$(dir $(STORAGE))"
	$(PYTHON) -m all_day --storage "$(STORAGE)" update "$(TASK_ID)" $(UPDATE_ARGS)

done:
	@if [ -z "$(TASK_ID)" ]; then echo "Usage: make done TASK_ID=1"; exit 1; fi
	@mkdir -p "$(dir $(STORAGE))"
	$(PYTHON) -m all_day --storage "$(STORAGE)" done "$(TASK_ID)"

remove:
	@if [ -z "$(TASK_ID)" ]; then echo "Usage: make remove TASK_ID=1"; exit 1; fi
	@mkdir -p "$(dir $(STORAGE))"
	$(PYTHON) -m all_day --storage "$(STORAGE)" remove "$(TASK_ID)"

summary:
	@mkdir -p "$(dir $(STORAGE))"
	$(PYTHON) -m all_day --storage "$(STORAGE)" summary $(SUMMARY_ARGS)

stop:
	@case "$(STORAGE)" in \
		./.demo/*|.demo/*) ;; \
		*) echo "Refusing to clear non-demo storage path: $(STORAGE)"; exit 1 ;; \
	esac
	@rm -f "$(STORAGE)"
	@rmdir "$(dir $(STORAGE))" 2>/dev/null || true
	@echo "Cleared demo storage: $(STORAGE)"

test:
	$(PYTHON) -m unittest discover -s tests -v

build:
	$(PYTHON) -m build --sdist --wheel

clean:
	@find all_day tests -type d -name '__pycache__' -prune -exec rm -rf {} +
	@rm -rf build dist *.egg-info all_day.egg-info
