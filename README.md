# user-system-microservice

Tech Stack: Go + Redis.

Go Libraries used: gin-gonic, go-redis

Database schema defined at- https://github.com/Aadithya-V/user-system-microservice/blob/master/internal/database/db.go

API defined informally in the package; strict definition and test cases to be added from Postman soon (insiders version- https://www.postman.com/flight-astronaut-52388155/workspace/user-system-microservice/collection/25140902-d4b0308d-e264-45b9-bfe2-e630d31242b6?action=share&creator=25140902)

Implements basic-auth (using Provos and Mazières's bcrypt adaptive hashing algorithm) and token-bearer authentication systems that are stateless, in-tune with REST principles.
