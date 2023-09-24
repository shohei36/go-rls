fetch_user: 
	curl -X GET http://localhost:8080/users/1000 -H "tenant-id:tenant01"

update_user:
	curl -X POST http://localhost:8080/users/1000 \
		-H "Content-Type: application/json" \
		-H "tenant-id:tenant01" \
		-d "{\"name\":\"Elen\",\"gender\":\"female\",\"age\":25}"