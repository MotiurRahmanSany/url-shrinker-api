# 🚀 Project Setup & Execution Guide (Beginner Friendly)

Welcome to the **URL Shrinker** project! This guide is written to help you set up and run this fullstack application on a new computer from scratch, even if you have no prior coding experience.

---

## 📋 Table of Contents
1. [Prerequisites & Installation](#1-prerequisites--installation)
2. [First-Time Configuration](#2-first-time-configuration)
3. [Running the Application](#3-running-the-application)
4. [Testing & Verifying the Features](#4-testing--verifying-the-features)
5. [Troubleshooting Guide](#5-troubleshooting-guide)

---

## 1. Prerequisites & Installation

To run this application, you must install the following software packages. Download the installers for your Operating System and follow the standard installation wizards:

* **[Node.js (LTS Version)](https://nodejs.org/)**: Required to compile and run the frontend website.
* **[Go Language (v1.25+)](https://go.dev/dl/)**: Required to run the backend API.
* **[Docker Desktop](https://www.docker.com/products/docker-desktop/)**: Automatically spins up the database (PostgreSQL) and the cache layer (Redis) inside sandboxed containers.
* **[Git](https://git-scm.com/)**: Used to clone and manage the codebase files.

> [!IMPORTANT]
> Once Docker Desktop is installed, **open the Docker Desktop application** and make sure it is running in the background before continuing.

---

## 2. First-Time Configuration

Open your terminal or command prompt and run the following steps to initialize the environment:

### Step 1: Copy the Environment Variables Config
Create the `.env` configuration file in the root folder of the project:
* **Mac / Linux:**
  ```bash
  cp .env.example .env
  ```
* **Windows (Command Prompt):**
  ```cmd
  copy .env.example .env
  ```

### Step 2: Spin Up Database & Redis Containers
Run Docker to pull and start the pre-configured databases in the background:
```bash
docker compose up -d
```
*(You can verify that PostgreSQL and Redis are successfully running by checking the green status icons in Docker Desktop).*

### Step 3: Run Database Migrations (Build Table Schema)
Create the database tables (users, urls, click events) inside your PostgreSQL container:
```bash
make migrate-up
```
*(If you do not have `make` installed on Windows, you can install the `goose` migration tool and run `goose -dir ./sql/migrations postgres "host=localhost port=5433 user=postgres password=postgres-password dbname=url_shrinker_db sslmode=disable" up`).*

---

## 3. Running the Application

To run the full stack, you will need to open **two terminal windows** (one for the backend, one for the frontend).

### Terminal 1: Run the Go Backend Server
From the root of the project, start the backend server with live reload active:
```bash
make run
```
The backend server will start running and outputting logs at **`http://localhost:8090`**.

---

### Terminal 2: Run the Next.js Frontend Client
1. Open a new terminal window and change directory into the `client` folder:
   ```bash
   cd client
   ```
2. Install the necessary packages (only needed the first time):
   ```bash
   npm install
   ```
3. Start the website using the stable Webpack compiler:
   ```bash
   npm run dev
   ```
The frontend website will start compiling and will open at **`http://localhost:3000`**.

---

## 4. Testing & Verifying the Features

Open your web browser and navigate to **`http://localhost:3000`** to play around with the app:

### 1. Sign Up & Log In
* Click **Sign Up** in the header.
* Enter any email address and a password (minimum 6 characters).
* Log in with those credentials.

### 2. Shorten a URL
* Paste a long destination link (e.g., `https://github.com`) into the home page input box.
* Click **Advanced Options** to customize:
  * **Custom Code**: Type in a custom path suffix (e.g. `gh`), making your URL `http://localhost:8090/gh`.
  * **Expiration Date**: Select a date/time in the future.
  * **Max Clicks**: Set a limit (e.g., `2`). The link will break after being clicked twice.
* Click **Shrink It!** and copy the generated link.

### 3. Verify Redirection
* Paste the copied short link into a new browser tab.
* The backend will record your click and redirect you to the destination (e.g., GitHub).

### 4. Check Analytics
* Go to the **Dashboard** in the header.
* Click the **Analytics (graph icon)** button next to your link.
* View total clicks, today's clicks, and a visual daily timeline chart showing engagement.
* Click the **Edit (pencil icon)** to update original URLs or clear expiration dates/click limits.
* Click **Deactivate (trash icon)** to manually disable a link.

---

## 5. Troubleshooting Guide

* **Issue: "Docker Daemon is not running"**
  * *Solution:* Make sure Docker Desktop is launched and the status indicator at the bottom-left of the Docker UI is green.
* **Issue: "Turbopack crash / panic logs on Linux"**
  * *Solution:* Next.js Turbopack dev server has known bugs on Linux. Make sure you run the client using `npm run dev` (which executes `next dev --webpack`) instead of `npx next dev --turbo`.
* **Issue: "Address already in use / Port conflict"**
  * *Solution:* Ensure you do not have another app running on ports `8090` or `3000`. Stop any other Node or Go processes running in the background.
* **Issue: "Links immediately expire after editing"**
  * *Solution:* Leaving the expiration date empty during an edit resets the URL to have no expiry date constraint. Ensure your client code is updated to map empty fields to Go zero time `"0001-01-01T00:00:00Z"`.
