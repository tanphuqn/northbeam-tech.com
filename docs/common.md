# northbeam-tech.com — Post Creation Guide

Common requirements for every post published to **northbeam-tech.com**.
Follow these rules so the marketing team always receives a consistent, ready-to-publish package.

There are two post types:
- **Insights** (`/insights/`) — technical how-to posts about a specific integration or technique
- **Projects** (`/projects/`) — project introduction posts showcasing a full product we built

---

## 1. Folder Structure

```
docs/
├── common.md                              ← this file (do not send to marketing)
│
├── insights/
│   └── <post-slug>/                       ← one folder per post
│       ├── post-<post-slug>.html          ← WordPress brief + content
│       ├── <slug>-thumbnail.png           ← featured image (1200×628 px)
│       ├── <slug>-flow-diagram.png        ← inline diagram inside post
│       └── <third-party>-logo.png         ← integration logo (if applicable)
│
└── projects/
    └── <project-name>/                    ← one folder per project
        ├── post-<project-name>-project.html
        ├── <project-name>-thumbnail.png   ← featured image (1200×628 px)
        ├── <project-name>-architecture-diagram.png
        └── <project-name>-logo.png        ← project's own logo (if available)
```

> Send only `.png` + `.html` files to marketing. Keep `.svg` sources in the repo if they exist.

---

## 2. File Naming Rules

| Rule | Example |
|---|---|
| Lowercase, hyphens only | `verify365-stripe-connect-multitenant` |
| Include focus keyword in slug | `pep-sanctions-monitoring-kyc6-go` |
| Suffix describes role | `-thumbnail`, `-flow-diagram`, `-architecture-diagram` |
| Third-party logos: `<name>-logo.png` | `stripe-logo.png`, `creditsafe-logo.png`, `kyc6-logo.png` |
| Project logo: `<project>-logo.png` | `verify365-logo.png`, `fac4fac-logo.png` |
| No spaces, no special characters | never `My Post Image (1).png` |

---

## 3. Required Files per Post

### Insights post

| File | Purpose | Notes |
|---|---|---|
| `post-<slug>.html` | Full WordPress brief for marketing | See Section 6 |
| `<slug>-thumbnail.png` | Featured Image | 1200×628 px |
| `<slug>-flow-diagram.png` | Inline diagram in post body | ~700 px wide |
| `<third-party>-logo.png` | Integration logo shown above H1 | 48 px tall in post |

### Projects post

| File | Purpose | Notes |
|---|---|---|
| `post-<name>-project.html` | Full WordPress brief for marketing | See Section 6 |
| `<name>-thumbnail.png` | Featured Image | 1200×628 px |
| `<name>-architecture-diagram.png` | System architecture inline diagram | 860×980 px |
| `<name>-logo.png` | Project's own logo shown above H1 | 56 px tall in post |

---

## 4. Image Specifications

### Thumbnail (Featured Image)
- **Size:** 1200 × 628 px
- **Style:** Dark background (`#0f1117` or project brand color), post title + tech stack tags
- **Alt text:** `<focus keyword> <short description> cover image`

### Flow Diagram (Insights)
- **Style:** Technical/minimal — dark or white background, labeled nodes
- **Width:** ~700 px

### Architecture Diagram (Projects)
- **Style:** Layered — client apps → API → integrations → data
- **Size:** 860 × 980 px (portrait)

### Logos
- **Third-party logos** (Stripe, CreditSafe, KYC6 etc.): created as SVG then exported to PNG using Node.js + sharp at `/tmp`
- **Project logos**: copy from project's own assets folder (prefer `color-logo.png` or `assets/logo.png`)
- **Placement in post HTML:** `<div style="margin-bottom:24px;"><img src="<name>-logo.png" alt="<Name> logo" style="height:48px;width:auto;"/></div>` directly above the `<h1>` tag

### Exporting SVG → PNG (Node.js + sharp)

All image generation uses inline SVG strings in a Node.js script at `/tmp`:

```bash
cd /tmp && npm install sharp 2>/dev/null

node -e "
const sharp = require('sharp');

const thumbnailSvg = \`<svg xmlns='http://www.w3.org/2000/svg' width='1200' height='628'>
  <!-- design here -->
</svg>\`;

const diagramSvg = \`<svg xmlns='http://www.w3.org/2000/svg' width='700' height='900'>
  <!-- design here -->
</svg>\`;

const BASE = '/Users/tanphuqn/Projects/me/northbeam-tech.com/docs/insights/YOUR-SLUG';
Promise.all([
  sharp(Buffer.from(thumbnailSvg)).png().toFile(BASE + '/YOUR-SLUG-thumbnail.png'),
  sharp(Buffer.from(diagramSvg)).png().toFile(BASE + '/YOUR-SLUG-flow-diagram.png'),
]).then(() => console.log('Done')).catch(console.error);
"
```

---

## 5. SEO Fields (required for every post)

| Field | Rule |
|---|---|
| **Focus keyword** | 2–4 words. Include in title, meta description, and first paragraph. |
| **Slug** | Same as folder name. Lowercase, hyphenated, focus keyword first. |
| **Meta description** | 120–160 characters. Contains focus keyword. Plain sentence, no markdown. |
| **Excerpt** | 1–2 sentences. Same meaning as meta description, can be shorter. |
| **Tags** | 4–8 tags: tech stack + topic keywords. |
| **Category** | `Insights` or `Projects` (create in WordPress if missing). |

---

## 6. Post HTML File Structure

Every `post-*.html` must contain these sections in order:

1. **Yellow instruction box** — step-by-step publishing guide for marketing (includes which files to upload)
2. **Post meta table** — title, slug, focus keyword, meta description, category, tags, excerpt
3. **Images table** — file names, sizes, alt texts for all images (including logos)
4. **AI rewrite prompt** — full context prompt for the marketing team to adjust copy
5. **`<!-- POST CONTENT START -->`** … `<!-- POST CONTENT END -->`

### Insights post content structure

```
[Logo image — above H1, height 48px]
H1 — Post title
Tagline — one sentence summary

H2 — The Problem
H2 — What We Built
H2 — How It Works   ← flow diagram goes here
H2 — Key Technical Decisions   ← real code snippets + explanation
H2 — Results / What Changed
H2 — Stack   ← monospace tag list

--- separator ---
Blockquote CTA → /contact
```

### Projects post content structure

```
[Project logo — above H1, height 56px]
H1 — Project name + one-line descriptor
Tagline — one sentence summary
Live badge (e.g. "Live at verify365.app")

H2 — About the Project
H2 — System Architecture   ← architecture diagram goes here
H2 — Client Applications
H2 — [Key Feature 1]
H2 — [Key Feature 2]
...
H2 — Technical Highlights   ← stack box + bullet list

--- separator ---
Related Insights links (link to posts about specific integrations)
Blockquote CTA → /contact
```

---

## 7. Marketing Instructions Block (copy template)

Paste this yellow box at the top of every post HTML, adjusting file names:

```html
<div class="instructions">
  <h2>📋 Instructions for Marketing</h2>
  <ol>
    <li>WordPress → <strong>Posts → Add New</strong>, set category to <code>Insights</code></li>
    <li>Paste the title below</li>
    <li>Switch to <strong>Code editor</strong>, paste content between <code>&lt;!-- POST CONTENT START --&gt;</code> markers</li>
    <li>Upload <code>SLUG-logo.png</code> → replace <code>src="SLUG-logo.png"</code> with WordPress URL</li>
    <li>Upload <code>SLUG-flow-diagram.png</code> → replace <code>src="SLUG-flow-diagram.png"</code> with WordPress URL</li>
    <li>Upload <code>SLUG-thumbnail.png</code> → set as <strong>Featured Image</strong></li>
    <li>Fill Yoast/Rank Math with Focus Keyword, Slug, Meta Description below</li>
    <li>Category + Tags below → Preview → Publish</li>
  </ol>
</div>
```

---

## 8. AI Rewrite Prompt Template

```
You are a technical content writer for a software consultancy called Northbeam Tech (northbeam-tech.com/insights/).

Here is a blog post draft about [SHORT PROJECT DESCRIPTION]. Please rewrite or adjust it based on my instructions below.

---
CONTEXT:
- Project name: [PROJECT NAME]
- Platform / framework: [e.g. SugarCRM, Go, PHP]
- Integration: [e.g. Twilio WhatsApp API, Stripe Connect, KYC6]
- What it does: [1–2 sentences describing the core feature]
- Key features: [comma-separated list]
- Tech stack: [language, database, APIs]
- Focus keyword for SEO: "[FOCUS KEYWORD]"
- Target audience: [who reads this]

---
CURRENT POST TITLE:
[POST TITLE]

CURRENT POST CONTENT:
[paste the post HTML content here]

MY INSTRUCTIONS:
[describe changes — e.g. "make it shorter", "add FAQ", "translate to Vietnamese"]
```

---

## 9. Existing Posts

### Projects

| Folder | Project | Live URL |
|---|---|---|
| `projects/verify365/` | Verify365 — KYC/KYB platform for UK law firms (Go) | verify365.app |
| `projects/fac4fac/` | fac4fac — SugarCRM customisation for Israeli org (PHP) | — |

### Insights

| Folder | Topic | Project |
|---|---|---|
| `insights/verify365-pep-sanctions-kyc6/` | PEP & sanctions monitoring via KYC6 API | Verify365 |
| `insights/verify365-stripe-connect-multitenant/` | Multi-tenant Stripe Connect payments | Verify365 |
| `insights/verify365-creditsafe-ubo-sftp/` | CreditSafe REST API + UBO SFTP feed | Verify365 |
| `insights/whatsapp-sugarcrm-integration/` | WhatsApp integration with SugarCRM via Twilio | fac4fac |
| `insights/whatsapp-student-followup-sugarcrm/` | Automated student follow-up via WhatsApp | fac4fac |
| `insights/whatsapp-template-approval-sync/` | WhatsApp template approval & sync | fac4fac |
| `insights/sugarcrm-excel-report-export/` | Dynamic Excel report export from SugarCRM | fac4fac |
| `insights/sugarcrm-pdf-student-export/` | RTL Hebrew PDF student export (TCPDF) | fac4fac |

---

## 10. Checklist Before Sending to Marketing

- [ ] Folder named with the post slug (lowercase, hyphenated)
- [ ] `post-*.html` has all 5 sections: instruction box, meta, images table, AI prompt, content
- [ ] Logo image placed above `<h1>` in post content (48 px Insights / 56 px Projects)
- [ ] Logo filename added to the upload instructions in the yellow box
- [ ] Focus keyword in: title, meta description, first paragraph, at least one H2
- [ ] Meta description 120–160 characters
- [ ] Thumbnail PNG is 1200×628 px
- [ ] All image `src` attributes use relative paths (marketing replaces with WordPress URL after upload)
- [ ] All image alt texts filled in
- [ ] Only `.png` + `.html` files sent to marketing (no `.svg` sources)
