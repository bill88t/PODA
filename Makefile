PREFIX := $(or $(PREFIX), )

PORT=5173

.PHONY: all
all: clean build

.PHONY: check-go
check-go:
	@command -v go >/dev/null 2>&1 || { echo >&2 "Error: go is not installed."; exit 1; }

# tsc already is running from vite
# tsc is not global always but it is always local to node nodules
# .PHONY: check-tsc
# check-tsc:
# 	@command -v frontend/node_modules/typescript/bin/tsc >/dev/null 2>&1 || { echo >&2 "Error: typescript is not installed."; exit 1; }

.PHONY: breadport
breadport:
	@$(eval PORT = 80)

.PHONY: frontend
frontend: # check-tsc - already from vite
	@echo "Building frontend.."
	@cd frontend; npm run build

build: check-go main frontend

npm:
	@echo "Refreshing NPM packages"
	@cd frontend; npm install

run:
	@echo "Running on port $PORT"
	@PORT=$(PORT) ./main

breadrun: breadport | run

main:
	@echo "Building backend.."
	@go build -o main .

clean:
	@echo "Cleaning build.."
	@rm -f main

full_clean: rmdb clean

rmdb:
	@echo "Wiping database.."
	@rm -v poda.db-wal 2>/dev/null || true
	@rm -v poda.db-shm 2>/dev/null || true
	@rm -v poda.db 2>/dev/null || true

test_bill_reg:
	@echo "Registering bill88t.."
	curl -X POST http://localhost:5173/api/v1/users/signup \
	  -H "Content-Type: application/json" \
	  -d '{"fname": "Bill", "lname": "Sideris", "email": "bill88t@feline.gr", "password": "securepassword123", "birthday": "2002-12-01" }'
	@echo
	@echo

test_bill_login:
	@echo "Testing login.."
	@curl -s -X POST http://localhost:5173/api/v1/users/login \
	  -H "Content-Type: application/json" \
	  -d '{ "email":  "bill88t@feline.gr", "password": "securepassword123" }' \
	  | tee /$(PREFIX)/tmp/login_response.json
	@echo
	@echo

test_bill_profile:
	@echo "Getting profile.."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/login_response.json))
	@curl -s -X GET http://localhost:5173/api/v1/profile/ \
	  -H "Authorization: Bearer $(TOKEN)"
	@echo
	@echo

test_bill_profile_fail:
	@echo "Testing profile without auth (should fail).."
	@curl -X GET http://localhost:5173/api/v1/profile/
	@echo
	@echo

test_bill_duplicate_reg:
	@echo "Testing duplicate registration (should fail).."
	@curl -X POST http://localhost:5173/api/v1/users/signup \
	  -H "Content-Type: application/json" \
	  -d '{"fname": "Bill", "lname": "Sideris", "email": "bill88t@feline.gr", "password": "securepassword123", "birthday": "2002-12-01" }'
	@echo
	@echo

test_bill_invalid_login:
	@echo "Testing invalid login (should fail).."
	@curl -X POST http://localhost:5173/api/v1/users/login \
	  -H "Content-Type: application/json" \
	  -d '{ "email":  "bill88t@feline.gr", "password": "wrongpassword" }'
	@echo
	@echo

test_bill_create_appointment:
	@echo "Creating appointment.."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/login_response.json))
	@curl -s -X POST http://localhost:5173/api/v1/profile/appointments/ \
	  -H "Content-Type: application/json" \
	  -H "Authorization: Bearer $(TOKEN)" \
	  -d '{"datetime": "2026-02-15T10:00:00Z", "kind": "haircut"}' \
	  | tee /$(PREFIX)/tmp/appointment_response.json
	@echo
	@echo

test_bill_get_appointments:
	@echo "Getting all appointments.."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/login_response.json))
	@curl -s -X GET http://localhost:5173/api/v1/profile/appointments/ \
	  -H "Authorization: Bearer $(TOKEN)"
	@echo
	@echo

test_bill_get_appointment_by_id:
	@echo "Getting appointment by ID.."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/login_response.json))
	@$(eval APPT_ID=$(shell jq -r '.id' /$(PREFIX)/tmp/appointment_response.json))
	@curl -s -X GET http://localhost:5173/api/v1/profile/appointments/$(APPT_ID) \
	  -H "Authorization: Bearer $(TOKEN)"
	@echo
	@echo

test_bill_update_appointment:
	@echo "Updating appointment.."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/login_response.json))
	@$(eval APPT_ID=$(shell jq -r '.id' /$(PREFIX)/tmp/appointment_response.json))
	@curl -s -X PUT http://localhost:5173/api/v1/profile/appointments/$(APPT_ID) \
	  -H "Content-Type: application/json" \
	  -H "Authorization: Bearer $(TOKEN)" \
	  -d '{"datetime": "2026-02-15T14:00:00Z", "kind": "shaving"}'
	@echo
	@echo

test_bill_delete_appointment:
	@echo "Deleting appointment.."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/login_response.json))
	@$(eval APPT_ID=$(shell jq -r '.id' /$(PREFIX)/tmp/appointment_response.json))
	@curl -s -X DELETE http://localhost:5173/api/v1/profile/appointments/$(APPT_ID) \
	  -H "Authorization: Bearer $(TOKEN)"
	@echo
	@echo

test_bill_create_second_appointment:
	@echo "Creating second appointment for testing.."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/login_response.json))
	@curl -s -X POST http://localhost:5173/api/v1/profile/appointments/ \
	  -H "Content-Type: application/json" \
	  -H "Authorization: Bearer $(TOKEN)" \
	  -d '{"datetime": "2026-03-20T09:00:00Z", "kind": "therapy lmao"}'
	@echo
	@echo

test_bill_appointment_not_found:
	@echo "Testing get non-existent appointment (should fail).."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/login_response.json))
	@curl -X GET http://localhost:5173/api/v1/profile/appointments/99999 \
	  -H "Authorization: Bearer $(TOKEN)"
	@echo
	@echo

test_bill_appointment_no_auth:
	@echo "Testing create appointment without auth (should fail).."
	@curl -X POST http://localhost:5173/api/v1/profile/appointments/ \
	  -H "Content-Type: application/json" \
	  -d '{"datetime": "2026-02-15T10:00:00Z", "kind": "checkup"}'
	@echo
	@echo

test_bill_invalid_appointment:
	@echo "Testing create appointment with invalid data (should fail).."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/login_response.json))
	@curl -X POST http://localhost:5173/api/v1/profile/appointments/ \
	  -H "Content-Type: application/json" \
	  -H "Authorization: Bearer $(TOKEN)" \
	  -d '{"datetime": "invalid-date", "kind": "robbery"}'
	@echo
	@echo

runtests: test_bill_reg test_bill_login test_bill_profile test_bill_create_appointment test_bill_get_appointments test_bill_get_appointment_by_id test_bill_update_appointment test_bill_create_second_appointment test_bill_get_appointments test_bill_delete_appointment

runtests_full: test_bill_reg test_bill_duplicate_reg test_bill_login test_bill_invalid_login test_bill_profile test_bill_profile_fail test_bill_create_appointment test_bill_get_appointments test_bill_get_appointment_by_id test_bill_update_appointment test_bill_create_second_appointment test_bill_get_appointments test_bill_appointment_not_found test_bill_invalid_appointment test_bill_appointment_no_auth test_bill_delete_appointment
