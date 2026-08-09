# Phu H. — Upwork Individual Profile

> Source of truth for individual freelancer profile copy.
>
> Profile: https://www.upwork.com/freelancers/~0154a937f03807cfa1
> Agency: https://www.upwork.com/agencies/2067609224625512954/

---

## Current State (from screenshot — Jun 2026)

| Field              | Current value                                                | Action needed                                                            |
| ------------------ | ------------------------------------------------------------ | ------------------------------------------------------------------------ |
| Title/Headline     | "Full Stack Developer: Laravel, PHP, Golang Expert 🚀"       | ✏️ replace — too generic                                                 |
| Overview           | Generic bio mentioning Laravel/Python/Flask/MongoDB          | ✏️ replace — doesn't reflect real projects                               |
| Hourly rate        | $30/hr                                                       | review — consider $45–60 given Top Rated Plus + 100% JSS                 |
| Job Success        | 100%                                                         | ✅ keep                                                                  |
| Badge              | Top Rated Plus                                               | ✅ keep                                                                  |
| Earned             | $500K+                                                       | ✅ keep                                                                  |
| Total jobs         | 16                                                           | ✅ keep                                                                  |
| Portfolio          | Only Verify365 visible                                       | ✏️ add kolel.org, SearchEye, fac4fac                                     |
| Skills             | PHP, MySQL, JavaScript, Laravel, Golang, Python, Docker, GCP | ✏️ replace Python/MySQL; add Rails, React Native, Stripe, Twilio, OpenAI |
| Employment history | Senior Software Developer / WCMS Solutions                   | ✅ keep                                                                  |
| Languages          | English                                                      | ✅ keep                                                                  |
| Associated with    | Northbeam Technologies                                       | ✅ keep                                                                  |

---

## 1. Title / Headline

**Replace** current title with:

```
Go, Laravel, Rails, Python Engineer - SaaS and API Integrations
```

_(63 chars — plain characters only, no special chars)_

---

## 2. Overview / Bio

> Individual profiles allow up to 5000 characters. Use the space to tell a story.
> Write in first person ("I"), unlike the agency "we".

**Replace entire overview with:**

```
I'm a senior backend engineer and founder of Northbeam Technologies, specializing
in Go, Laravel, Python and Ruby on Rails — building production SaaS platforms with complex
third-party integrations that go live and stay maintained.

Top Rated Plus · 100% Job Success Score · $100K+ earned on Upwork.

────────────────────────────────────────
WHAT I BUILD
────────────────────────────────────────

▸ Backend APIs & SaaS platforms
  Go (Gin, GORM), Laravel 9, Ruby on Rails 7, PHP, Python — REST APIs,
  multi-tenant billing, compliance workflows, background job pipelines.

▸ Third-party integrations
  Stripe Connect, Twilio WhatsApp, KYC/KYB APIs (KYC6, CreditSafe, LexisNexis),
  Open Banking (Basiq, Nordigen), e-signatures (BoldSign), OAuth2 practice
  management (Clio, ActionStep), Algolia, OpenAI, Google Cloud Video.

▸ Mobile & web frontends
  React Native (Expo) iOS & Android apps. Nuxt 2, Vue 2, React SPAs.
  SSR, RTL-aware, EAS build pipelines and Fastlane store releases.

▸ CRM customization
  SugarCRM module development, Twilio webhook infrastructure, PDF/Excel
  export engines.

────────────────────────────────────────
SELECTED PROJECTS
────────────────────────────────────────

◆ Verify365 — KYC/KYB compliance platform for UK law firms
  I built the entire Go REST API (Gin + GORM + PostgreSQL on AWS) powering
  4 client apps. 20+ integrations: KYC6 PEP/sanctions monitoring, CreditSafe
  company credit + UBO via SFTP, Stripe Connect multi-tenant payments, BoldSign
  e-signatures, Open Banking, LexisNexis IDU, Companies House, FormEvo UK legal
  forms, Clio and ActionStep practice management.

◆ Kolel — Full-stack Jewish video streaming platform
  Built from scratch: Rails 7 API, Nuxt 2 SSR web (Hebrew RTL), Vue 2 creator
  dashboard with FFmpeg.wasm in-browser transcoding, React Native iOS & Android
  apps. 4.9 stars, 100K+ downloads. Google Cloud Video Transcoder, Elasticsearch,
  AWS S3, Stripe, Firebase FCM, ActionCable WebSockets.

◆ SearchEye — AI-powered SEO & link-building SaaS
  Laravel 9 + Vue 2.7 platform for digital agencies. Algolia backlink
  marketplace, ChatCells proprietary OpenAI GPT-4 content engine, dual
  Stripe + PayPal billing, Ahrefs + DataForSEO + Diffbot integrations.
  172 Eloquent models, 457 API routes, 73 background jobs.

◆ fac4fac — Custom SugarCRM with WhatsApp automation
  15+ custom PHP modules on SugarCRM 6.5.23: student lifecycle CRM, two-way
  WhatsApp messaging via Twilio (templates, bulk send, approval sync), Hebrew
  RTL PDF export (TCPDF), dynamic Excel reporting. Dual AWS deployment.

────────────────────────────────────────
HOW I WORK
────────────────────────────────────────

I take ownership of systems end-to-end — not just tickets. Every integration
I write gets a dedicated package with typed interfaces, isolated auth, and
consistent error handling so the codebase is maintainable years later.

I'm based in Ho Chi Minh City, Vietnam (UTC+7) and communicate clearly in
English throughout every engagement. Available for long-term contracts.

Let's talk about your project.
```

_(~2,850 chars — well within 5000 limit, room to expand if needed)_

---

## 3. Skills

**Current:** PHP, MySQL, JavaScript, Laravel, Golang, Python, Docker, Google Cloud Platform

**Recommended 15:**

```
Golang · Ruby on Rails · Laravel · PHP · Python · Vue.js · React Native ·
PostgreSQL · Amazon EC2 · Google Cloud Platform · SugarCRM Development ·
API Integration · Stripe · Twilio · OpenAI API
```

**Remove:** MySQL (redundant with PostgreSQL), Docker (not a client search term), JavaScript (too generic), Redis (infrastructure detail — drop to make room for Python)

**Add:** Ruby on Rails, React Native, PostgreSQL, Python, SugarCRM Development, API Integration, Stripe, Twilio, OpenAI API

---

## 4. Portfolio

> Individual profile form limits: title 70 chars, description 600 chars, skills 5 max.
> Cover images are in `upwork-agency/`.

### Verify365

```
Title:       Verify365 - KYC/KYB Compliance Platform for UK Law Firms
Role:        Lead Backend Engineer
Cover:       upwork-agency/portfolio-verify365-cover.png

Description:
Built the entire Go REST API (Gin + GORM + PostgreSQL on AWS) for a KYC/KYB
compliance platform serving UK law firms. Delivered 20+ third-party integrations:
KYC6 PEP/sanctions screening, CreditSafe company credit + UBO via SFTP, Stripe
Connect multi-tenant payments, BoldSign e-signatures, Open Banking, LexisNexis IDU,
Companies House, FormEvo UK legal forms, Clio and ActionStep practice management.
Powers 4 client apps: firm dashboard, customer portal, iOS/Android mobile app,
and partner OAuth integrations.

Skills:      Golang · PostgreSQL · Amazon EC2 · API Integration · Stripe
```

### Kolel

```
Title:       Kolel - Full-Stack Video Streaming Platform (Rails + React Native)
Role:        Full-Stack Engineer
Cover:       upwork-agency/portfolio-kolel-cover.png

Description:
Built a full-stack video streaming platform from scratch for a global Jewish
audience. Rails 7 JSON API, Nuxt 2 SSR web with Hebrew RTL support, Vue 2 creator
dashboard with FFmpeg.wasm in-browser video transcoding, and React Native (Expo)
iOS & Android apps. The mobile apps reached 4.9 stars with 100K+ downloads.
Stack: Google Cloud Video Transcoder, Elasticsearch, AWS S3, Stripe, Firebase FCM,
ActionCable WebSockets, Sidekiq background jobs.

Skills:      Ruby on Rails · React Native · Vue.js · Google Cloud Platform · Amazon EC2
```

### SearchEye

```
Title:       SearchEye - AI-Powered SEO & Link-Building SaaS (Laravel + Algolia)
Role:        Lead Backend & Full-Stack Engineer
Cover:       upwork-agency/portfolio-searcheye-cover.png

Description:
Built and maintained a large-scale SaaS platform for SEO agencies: 172 Eloquent
models, 457 API routes, 73 background jobs. Features: Algolia-powered backlink
marketplace with instant faceted search; ChatCells — a proprietary database-driven
OpenAI GPT-4 content engine; dual payment processing (Stripe + PayPal); SEO data
from Ahrefs, DataForSEO, and Diffbot; multi-tenant role system for agencies,
clients, brokers, and publishers. Stack: Laravel 9, PHP 8, Vue 2.7, MySQL, Redis,
Laravel Horizon.

Skills:      Laravel · Vue.js · OpenAI API · Algolia · Stripe
```

### fac4fac

```
Title:       fac4fac - Custom SugarCRM with WhatsApp Automation (Twilio)
Role:        Full-Stack PHP & CRM Engineer
Cover:       upwork-agency/portfolio-fac4fac-cover.png

Description:
Built 15+ custom PHP modules on SugarCRM 6.5.23 for an Israeli educational
organization: full student lifecycle CRM, two-way WhatsApp messaging via Twilio
(templates, bulk send, delivery webhooks, approval sync), bulk Hebrew RTL PDF
export (TCPDF), dynamic color-coded Excel weekly reports, and a custom
query/report builder. Dual-instance AWS deployment across two organizations
(EC2 + RDS MySQL).

Skills:      PHP · SugarCRM Development · Twilio · Amazon EC2 · API Integration
```

---

## 5. Hourly Rate

**Current: $30/hr** (intentional decision — keeping for volume)

Consider raising to **$45–60/hr** given Top Rated Plus + 100% JSS + $100K+ earned.

---

## 6. Consultation

> Upwork 1-on-1 paid Zoom consultation listing.

**Category:** Web, Mobile & Software Dev

**Custom topics (5):** API Integration · SaaS Architecture · Backend Development · Stripe Integration · React Native

**Rate:** $50/30 min

**Description:**
```
Book a 1-on-1 call if you need expert advice on your backend architecture, API integrations, or mobile app structure.

I can help with:
• Go (Gin/GORM) REST API design and performance
• Laravel or Ruby on Rails SaaS architecture
• Third-party integrations: Stripe Connect, Twilio WhatsApp, KYC/KYB APIs, Open Banking, e-signatures
• React Native (Expo) iOS & Android app structure
• SugarCRM module customization and WhatsApp automation
• Multi-tenant billing and subscription systems
• Database schema design (PostgreSQL, MySQL)
• Code review and architectural feedback

I've built production systems with 20+ integrations for UK law firms, video streaming platforms (100K+ downloads), and SEO SaaS platforms. Top Rated Plus with 100% Job Success Score.

Come with your codebase, architecture diagram, or just a description of what you're building — I'll give you concrete, actionable advice.
```

**Documents for client:** Meeting summary ✅ · Delivery: 1 day

**Client requirement (250 chars max):**
```
Describe your project, the tech stack, and the main problem you're facing. Include any relevant links (repo, staging URL, or architecture doc) so I can prepare before the call.
```

**FAQ:**

Q: What tech stacks can you advise on?
```
Go, Laravel, Ruby on Rails, PHP, React Native, Vue.js, Nuxt. Integrations: Stripe Connect, Twilio, KYC/KYB APIs, Open Banking, Algolia, OpenAI, SugarCRM. Database design in PostgreSQL and MySQL.
```

Q: What should I prepare before the call?
```
Share your project description when booking. If you have a codebase, architecture diagram, or specific error, send it ahead of time. The more context you provide, the more concrete advice I can give.
```

Q: Can you help review my existing codebase or architecture?
```
Yes. Share your repo or architecture diagram before the call and I'll review it in advance. I can give feedback on structure, performance bottlenecks, integration design, and refactoring priorities.
```

**Gallery media:**
- `upwork-agency/agency-overview-banner.png`
- `upwork-agency/portfolio-verify365-cover.png`
- `upwork-agency/portfolio-kolel-cover.png`
- `upwork-agency/portfolio-searcheye-cover.png`
- `upwork-agency/portfolio-fac4fac-cover.png`

---

## Priority Order

- [x] Replace title/headline (Section 1)
- [x] Replace overview/bio (Section 2)
- [x] Update skills (Section 3)
- [ ] Add Kolel, SearchEye, fac4fac portfolio items (Section 4)
- [x] Consultation listing created (Section 6)
