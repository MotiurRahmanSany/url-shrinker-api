# URL Shrinker Frontend — Next.js Client

This is the premium React/Next.js frontend for the **URL Shrinker API**. It features a modern, glassmorphic dark-themed user interface styled with **Tailwind CSS v4** and equipped with rich animations.

---

## 🚀 Features

* **🔗 Link Shrinking:** Shorten URLs with optional advanced features:
  * **Custom codes** (e.g., `/my-cool-link`)
  * **Temporal expiration limits** (automatic deactivation)
  * **Max click counts** (automatic deactivation when reached)
* **🔑 Authentication Flow:** Secure JWT token access featuring:
  * Signup and Login views
  * Client-side persistent session state via React Context
  * Automatic silent token rotation (access token refresh) on 401 response codes
* **📊 Management Dashboard:**
  * Listing user's links in a paginated table
  * Real-time click counting
  * Live status toggles (Active vs Inactive)
  * Quick access edit modal for updating constraints (URLs, limits, and expiration)
  * Soft deactivation support
* **📈 Deep Analytics Details:** Detailed statistics view (`/stats/{code}`) featuring:
  * Total click counter
  * Today's click count
  * CSS-animated daily timeline activity chart

---

## 🛠️ Technology Stack

* **Framework:** Next.js 16.2+ (App Router)
* **Library:** React 19
* **Styles:** Tailwind CSS v4 & custom dark glassmorphism
* **Icons:** Lucide React
* **Language:** TypeScript

---

## 🏃 Getting Started

### Prerequisites

* Node.js v18+
* npm or pnpm / yarn

### Installation

1. Navigate to the client directory:
   ```bash
   cd client
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. (Optional) Configure environment variables. By default, the app points to `http://localhost:8090` (matching the backend). To override this:
   Create a `.env.local` file inside the `client` directory:
   ```env
   NEXT_PUBLIC_API_URL=http://localhost:8090
   ```

### Running the Development Server

Due to a known bug in Next.js Turbopack compiler (`next dev --turbo`) under Linux environments that can cause compilation panics and loop requests, **always run the server using Webpack compiler**:

```bash
npm run dev
```

This runs the custom dev script: `next dev --webpack` to ensure a stable and reliable compilation.

The application will launch at [http://localhost:3000](http://localhost:3000).

---

## 📂 Project Structure

```text
client/
├── app/
│   ├── components/      # UI components (Navbar, ShrinkForm)
│   ├── context/         # AuthContext handling JWT state
│   ├── dashboard/       # Dashboard list page
│   ├── login/           # Login interface
│   ├── register/        # Signup interface
│   ├── stats/[code]/    # Dynamic route for click analytics and timelines
│   ├── globals.css      # Core styles and dark theme overrides
│   └── layout.tsx       # Root wrapper injecting Navbar and AuthProvider
├── lib/
│   └── api.ts           # Axios-like fetch wrapper with auto token refresh logic
└── package.json
```

---

## 🔧 Fullstack Integration Details

The client communicates with the backend via the centralized fetch wrapper located at [lib/api.ts](file:///home/lazy/Projects/golang__/url-shrinker-api/client/lib/api.ts).

### Silent Token Rotation

The fetch wrapper intercepts any unauthorized `401` errors, securely reads the `refresh_token` from `localStorage`, hits `/auth/refresh` to get a fresh token pair, and seamlessly retries the original failed user request. If the refresh token itself has expired or is invalid, it clears credentials and redirects the user to the login screen.

### Time Expiration Zero-Time Handling

To remove an expiration date constraint on a link during editing, the client sends `"0001-01-01T00:00:00Z"` (Go zero time representation). The backend parses this, detects it as a zero value, and updates the database record column to `NULL` (no limit).
