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
build_sqlite: check-go main_sqlite frontend

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

main_sqlite:
	@echo "Building backend with sqlite.."
	@go build sqlitemain/main-sqlite.go
	@mv main-sqlite main

.PHONY: clean
clean:
	@echo "Cleaning build.."
	@rm -f main
	@rm test.db test.db-shm test.db-wal 2>/dev/null || true

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

test_barber_reg:
	@echo "Registering barber account.."
	@curl -s -X POST http://localhost:5173/api/v1/users/signup \
	  -H "Content-Type: application/json" \
	  -d '{"fname": "Tony", "lname": "Barber", "email": "tony@barber.gr", "password": "barberpassword123", "birthday": "1985-06-15", "admin_token": "poda-barber-admin-token"}' \
	  | tee /$(PREFIX)/tmp/barber_reg_response.json
	@echo
	@echo

test_bill_changepassword:
	@echo "Changing bill's password.."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/login_response.json))
	@$(eval UUID=$(shell jq -r '.user.id' /$(PREFIX)/tmp/login_response.json))
	@curl -s -X POST http://localhost:5173/api/v1/users/changepassword \
	  -H "Content-Type: application/json" \
	  -H "Authorization: Bearer $(TOKEN)" \
	  -d "{\"uuid\": \"$(UUID)\", \"password\": \"newpassword456\"}"
	@echo
	@echo

test_bill_changeinfo:
	@echo "Changing bill's info.."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/login_response.json))
	@$(eval UUID=$(shell jq -r '.user.id' /$(PREFIX)/tmp/login_response.json))
	@curl -s -X POST http://localhost:5173/api/v1/users/changeinfo \
	  -H "Content-Type: application/json" \
	  -H "Authorization: Bearer $(TOKEN)" \
	  -d "{\"uuid\": \"$(UUID)\", \"fname\": \"William\", \"lname\": \"Sideris\", \"birthday\": \"2002-12-01\"}"
	@echo
	@echo

test_bill_changecontact:
	@echo "Changing bill's contact.."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/login_response.json))
	@$(eval UUID=$(shell jq -r '.user.id' /$(PREFIX)/tmp/login_response.json))
	@curl -s -X POST http://localhost:5173/api/v1/users/changecontact \
	  -H "Content-Type: application/json" \
	  -H "Authorization: Bearer $(TOKEN)" \
	  -d "{\"uuid\": \"$(UUID)\", \"email\": \"bill88t@feline.gr\", \"phone\": \"+30 210 1234567\"}"
	@echo
	@echo

test_bill_change_no_auth:
	@echo "Testing changepassword without auth (should fail).."
	@curl -s -X POST http://localhost:5173/api/v1/users/changepassword \
	  -H "Content-Type: application/json" \
	  -d '{"password": "hackpass"}'
	@echo
	@echo

test_barber_login:
	@echo "Logging in as barber.."
	@curl -s -X POST http://localhost:5173/api/v1/users/login \
	  -H "Content-Type: application/json" \
	  -d '{"email": "tony@barber.gr", "password": "barberpassword123"}' \
	  | tee /$(PREFIX)/tmp/barber_login_response.json
	@echo
	@echo

test_barber_get_all_appointments:
	@echo "Barber fetching all appointments.."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/barber_login_response.json))
	@curl -s -X GET http://localhost:5173/api/v1/barber/appointments \
	  -H "Authorization: Bearer $(TOKEN)"
	@echo
	@echo

test_barber_cancel_appointment:
	@echo "Barber cancelling appointment.."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/barber_login_response.json))
	@$(eval APPT_ID=$(shell jq -r '.id' /$(PREFIX)/tmp/appointment_response.json))
	@curl -s -X DELETE http://localhost:5173/api/v1/barber/appointments/$(APPT_ID) \
	  -H "Authorization: Bearer $(TOKEN)"
	@echo OK
	@echo

test_barber_forbidden_for_client:
	@echo "Testing barber route as client (should fail).."
	@$(eval TOKEN=$(shell jq -r '.token' /$(PREFIX)/tmp/login_response.json))
	@curl -s -X GET http://localhost:5173/api/v1/barber/appointments \
	  -H "Authorization: Bearer $(TOKEN)"
	@echo
	@echo

runtests: test_bill_reg test_bill_login test_bill_profile test_bill_create_appointment test_bill_get_appointments test_bill_get_appointment_by_id test_bill_update_appointment test_bill_create_second_appointment test_bill_get_appointments test_bill_changepassword test_bill_changeinfo test_bill_changecontact test_barber_reg test_barber_login test_barber_get_all_appointments test_bill_delete_appointment

runtests_full: test_bill_reg test_bill_duplicate_reg test_bill_login test_bill_invalid_login test_bill_profile test_bill_profile_fail test_bill_create_appointment test_bill_get_appointments test_bill_get_appointment_by_id test_bill_update_appointment test_bill_create_second_appointment test_bill_get_appointments test_bill_appointment_not_found test_bill_invalid_appointment test_bill_appointment_no_auth test_bill_changepassword test_bill_changeinfo test_bill_changecontact test_bill_change_no_auth test_barber_reg test_barber_login test_barber_get_all_appointments test_barber_forbidden_for_client test_bill_create_appointment test_barber_cancel_appointment test_bill_get_appointments
