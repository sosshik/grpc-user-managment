# gRPC User Service

## Overview

This is my solution to problem **4.1 User management: grpc**.
This project implements a gRPC user service with the following features:
- Create user
- Get user by email
- Get user by ID
- Get list of users
- Update user
- Delete user


## How to run

Clone the repo: 

    git clone https://git.foxminded.ua/foxstudent106264/task-4.1.git

Create `.env` file in cmd directory with parameters: 
- `PORT` - port where you wish to start the bot
- `DATABASE_URL` - your MongoDB connection string
- `CONN_CHECK` - use true or false to enable connection check
- `RECONN_TIME` - time before next connection check
- `LOG_LEVEL` - used to set log level
- `RECONN_TRIES` - used to set amount of reconnections in a row

Run the app from cmd directory:

    go run main.go
