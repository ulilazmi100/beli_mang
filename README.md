# Beli Mang

## 🚀 Welcome to Beli Mang!

Beli Mang is an amazing ride-hailing application designed specifically for order delivery, allowing users to conveniently purchase food and drinks from local merchants. Developed as part of Project Sprint Batch 2 Week 4, this project exemplifies rapid development of high-quality, scalable applications with rigorous testing and load testing.

## 🌄 Background

Beli Mang leverages advanced algorithms such as the Haversine formula and the Traveling Salesman Problem with Nearest Neighbor to optimize delivery routes. The backend is built using Golang with either the labstack/echo or gofiber/fiber frameworks, using pure SQL without an ORM, and the pgx database library.

## 📂 Project Structure

- **Main Folder:** Uses the gofiber/fiber framework.
- **beli_mang_echo Folder:** Original implementation using the labstack/echo framework.

---

## 🚀 Getting Started

### Prerequisites

- **Go** (1.22 or later)
- **PostgreSQL**
- **K6** for load testing ([Installation Guide](https://k6.io/docs/get-started/installation/))
- **WSL** (Windows Subsystem for Linux) if on Windows

### Installation and Setup

1. **Clone the repository:**

    ```bash
    git clone https://github.com/ulilazmi100/beli_mang.git
    cd beli_mang
    ```

2. **Set up environment variables:**

    Create a `.env` file or export these variables in your shell:

    ```bash
    export DB_NAME=your_db_name
    export DB_PORT=5432
    export DB_HOST=localhost
    export DB_USERNAME=your_db_user
    export DB_PASSWORD=your_db_password
    export DB_PARAMS="sslmode=disable"  # or "sslrootcert=rds-ca-rsa2048-g1.pem&sslmode=verify-full" for production
    export APP_PORT=8080 # or your designated app port
    export JWT_SECRET=your_jwt_secret
    export BCRYPT_SALT=8 or 10  # or higher for production

    export AWS_ACCESS_KEY_ID=
    export AWS_SECRET_ACCESS_KEY=
    export AWS_S3_BUCKET_NAME=projectsprint-bucket-public-read
    export AWS_REGION=ap-southeast-1
    ```

3. **Run database migrations:**

    ```bash
    migrate -database "postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?$DB_PARAMS" -path db/migrations up
    ```

4. **Start the server:**

    ```bash
    go run main.go
    ```

5. **Optional: Rollback migrations if needed:**

    ```bash
    migrate -database "postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?$DB_PARAMS" -path db/migrations down
    ```

### Running with Docker

1. **Build Docker image:**

    ```bash
    docker build -t beli_mang .
    ```

2. **Run Docker container:**

    ```bash
    docker run -p 8080:8080 --env-file .env beli_mang
    ```

---

## 🧪 Testing

### Prerequisites

- **K6** ([Installation Guide](https://k6.io/docs/get-started/installation/))
- Linux environment (WSL/MacOS)

### Running Tests

1. **Set environment variable:**

    ```bash
    export BASE_URL=http://localhost:8080
    ```

2. **Run the server.**

3. **Navigate to the `tests` folder and run the tests:**

    #### Regular testing
    1. Open two terminals.
    2. Run the caddy backend in one terminal:
        ```bash
        make run
        ```
    3. In the other terminal, run the tests:
        ```bash
        make test
        ```
    #### Debugging
    For detailed debugging:
    ```bash
    make test-debug
    ```

    > 💡 `make test` & `make test-debug` will automatically set the BASE_URL of the backend to `http://localhost:8080`. To change it, use:
    > ```bash
    > BASE_URL=http://backend.url:8080 make test
    > ```

    #### Load testing
    For load testing, use:
    ```bash
    make load-test
    ```

### Tests Result
- My personal test results or benchmarks can be seen in the `benchmark` folder

---

## 🛠 Database Migration

Use [golang-migrate](https://github.com/golang-migrate/migrate) for managing database migrations:

- **Create migration:**

    ```bash
    migrate create -ext sql -dir db/migrations add_user_table
    ```

- **Execute migration:**

    ```bash
    migrate -database "postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?$DB_PARAMS" -path db/migrations up
    ```

- **Rollback migration:**

    ```bash
    migrate -database "postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?$DB_PARAMS" -path db/migrations down
    ```

---

## 📝 Requirements

For detailed functional requirements, please refer to the [Notion page](https://openidea-projectsprint.notion.site/BeliMang-7979300c7ce54dbf8ecd0088806eff14).

---

## 👥 Contributing

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a pull request.

---

## 📝 License

-.

---

## 📚 Resources

- **Notion:** [Beli Mang's Requirements' Notion Page](https://openidea-projectsprint.notion.site/BeliMang-7979300c7ce54dbf8ecd0088806eff14)
- **Tests:** [Project Sprint Batch 2 Week 4 Test Cases](https://github.com/nandanugg/BeliMangTestCasesPB2W4)
- **Migrations:** [Golang Migration](https://github.com/golang-migrate/migrate)

---

## 📞 Contact

[Muhammad Ulil 'Azmi](https://github.com/ulilazmi100) - [@M_Ulil_Azmi](https://twitter.com/M_Ulil_Azmi)

Project Link: [https://github.com/ulilazmi100/beli_mang](https://github.com/ulilazmi100/beli_mang)

---

## ⚠️ Note

Please be aware that the Amazon service subscription may have already ended, which could result in pipeline failures in GitHub.

---
