# URL Shrinker API - Industrial Attachment Report (LaTeX Source)

This directory contains the complete LaTeX source code for your Industrial Attachment Report, structured and written according to your project guidelines.

## 📁 Directory Structure
*   `main.tex` - The primary driver file containing package imports, layout styles, and chapter inclusion configs.
*   `references.bib` - BibTeX database of your citations (IEEE/APA format) including Go, Postgres, Redis, Next.js, and ethics articles.
*   `chapters/` - Separate `.tex` files for each chapter:
    *   `titlepage.tex` - Title page template with student and supervisor fields.
    *   `declaration.tex` - Academic honesty statement.
    *   `acknowledgement.tex` - Acknowledging mentors at Texlab IT and the university.
    *   `executive_summary.tex` - 1-page summary of the report.
    *   `abbreviations.tex` - Table of acronyms used in the report.
    *   `introduction.tex` - Chapter 1: Background, POs, objectives, scope, methodology.
    *   `organization.tex` - Chapter 2: Texlab IT Rajshahi history, structure, services, compliance.
    *   `project_overview.tex` - Chapter 3: Redirection mechanics (302 found), base62, Cache-Aside pattern.
    *   `tools_techniques.tex` - Chapter 4: Go 1.25, Redis, PostgreSQL, sqlc, goose, Docker, ZSET rate limiting.
    *   `societal_considerations.tex` - Chapter 5: Workplace safety, GDPR privacy compliance, green computing (energy efficiency of Go).
    *   `ethics_responsibilities.tex` - Chapter 6: Preventing link hijacking/phishing, ACM/IEEE ethics code compliance, cascading logs deletes.
    *   `skills_learning.tex` - Chapter 7: Technical & self-directed learning, adaptation, lifelong learning reflection.
    *   `challenges_solutions.tex` - Chapter 8: Next.js Turbopack compiler fix, Docker startup checks, Redis atomic Tx pipelines.
    *   `contribution.tex` - Chapter 9: Database schema migrations, routing, cache service, async background jobs, supervisor feedback.
    *   `conclusion.tex` - Chapter 10: Summary, recommendations for university curriculum and host company.
    *   `appendices.tex` - Appendix containing `docker-compose.yaml` listings, Go source code snippets, and daily training logs.

---

## 🚀 How to Compile

### Option A: Overleaf (Recommended)
1.  **Zip the `report` folder:** Run a zip command on your terminal, or create a zip of the `report/` directory.
2.  **Upload to Overleaf:**
    *   Go to [Overleaf](https://www.overleaf.com/) and log in.
    *   Click **New Project** > **Upload Project**.
    *   Upload the zipped file.
3.  **Compile:**
    *   Ensure the compiler is set to **pdfLaTeX** (default in Overleaf).
    *   Click **Recompile** to generate your PDF report.

### Option B: Local Compilation
If you have a local LaTeX environment installed (such as TeX Live, MacTeX, or MiKTeX):
1.  Navigate into this `report/` directory.
2.  Run the compilation commands:
    ```bash
    pdflatex main.tex
    bibtex main
    pdflatex main.tex
    pdflatex main.tex
    ```
3.  This generates a `main.pdf` document in this folder.

---

## ✏️ Placeholders to Edit
Open the files inside the `chapters/` directory and replace the bracketed placeholders with your actual details:
*   `[Your Student ID]` in `titlepage.tex` and `declaration.tex`
*   `[Your University Name]` / `[Your University/Institute Name]` in `titlepage.tex`, `declaration.tex`, `acknowledgement.tex`
*   `[Supervisor's Name]` in `titlepage.tex`, `declaration.tex`, `acknowledgement.tex`
*   `[Supervisor's Designation]` in `titlepage.tex`, `declaration.tex`
