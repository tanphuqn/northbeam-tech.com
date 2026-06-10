# Insights Post — Common Requirements

This file documents the standard requirements for every post published to **northbeam-tech.com/insights/**.
Follow these rules for every new post so the marketing team receives a consistent, ready-to-publish package.

---

## 1. Folder Structure

Each post gets its own folder inside `docs/insights/`. Send the **entire folder** to the marketing team.

```
docs/insights/
├── common.md                          ← this file (do not send to marketing)
└── <post-slug>/                       ← one folder per post (send this to marketing)
    ├── post-<post-slug>.html          ← WordPress post brief + content
    ├── <post-slug>-thumbnail.png      ← featured image (1200×628 px)
    ├── <post-slug>-thumbnail.svg      ← thumbnail source (keep for re-export)
    ├── <post-slug>-<diagram-name>.png ← inline diagram(s) used inside post
    └── <post-slug>-<diagram-name>.svg ← diagram source (keep for re-export)
```

**Example:**
```
whatsapp-sugarcrm-integration/
├── post-whatsapp-sugarcrm.html
├── whatsapp-sugarcrm-integration-thumbnail.png
├── whatsapp-sugarcrm-integration-thumbnail.svg
├── whatsapp-sugarcrm-integration-flow-diagram.png
└── whatsapp-sugarcrm-integration-flow-diagram.svg
```

> Send only the `.png` + `.html` files to marketing. Keep `.svg` sources in the repo.

---

## 2. File Naming Rules (SEO)

| Rule | Example |
|---|---|
| Lowercase only | `whatsapp-sugarcrm-integration` |
| Words separated by hyphens | `flow-diagram` not `flow_diagram` |
| Include the focus keyword | `whatsapp-sugarcrm-integration-thumbnail.png` |
| Descriptive suffix | `-thumbnail`, `-flow-diagram`, `-architecture` |
| No spaces, no special characters | never `My Post Image (1).png` |

---

## 3. Required Files per Post

| File | Purpose | Notes |
|---|---|---|
| `post-<slug>.html` | Full WordPress brief for marketing | See Section 6 for structure |
| `<slug>-thumbnail.png` | WordPress Featured Image | 1200×628 px (Open Graph standard) |
| `<slug>-<diagram>.png` | Inline image(s) inside post body | Width depends on diagram |
| `.svg` sources | Re-export at any size later | Keep in repo, do not send to marketing |

---

## 4. Image Specifications

### Thumbnail (Featured Image)
- **Size:** 1200 × 628 px
- **Format:** PNG (exported from SVG source)
- **Style:** Dark background (#0f1117), branded, includes post title + tech stack tags
- **Must include:** post title text, focus keyword visible, site URL (`northbeam-tech.com/insights`)
- **Alt text format:** `<focus keyword> <short description> cover image`
  - Example: `WhatsApp SugarCRM integration automated student onboarding cover image`

### Inline Diagrams
- **Format:** PNG (exported from SVG source)
- **Style:** Technical/minimal — black & white, monospace font, labeled nodes (START/PROCESS/DECISION/END)
- **Width:** match the SVG `width` attribute (typically 700 px)
- **Alt text format:** `<focus keyword> <diagram type> diagram`
  - Example: `WhatsApp SugarCRM integration flow diagram`

### Exporting SVG → PNG (Node.js + sharp)
```bash
cd /tmp
npm install sharp   # one-time

node -e "
const sharp = require('sharp');
const fs = require('fs');
const files = [
  { svg: 'path/to/<slug>-thumbnail.svg',   png: 'path/to/<slug>-thumbnail.png',   width: 1200 },
  { svg: 'path/to/<slug>-flow-diagram.svg', png: 'path/to/<slug>-flow-diagram.png', width: 700  },
];
(async () => {
  for (const f of files) {
    await sharp(fs.readFileSync(f.svg)).resize(f.width).png({ quality: 95 }).toFile(f.png);
    console.log('Created:', f.png);
  }
})();
"
```

---

## 5. SEO Fields (required for every post)

| Field | Rule |
|---|---|
| **Focus keyword** | 2–4 words, matches the core topic. Include in title, meta description, and first paragraph. |
| **Slug** | Same as the folder name. Lowercase, hyphenated, starts with the focus keyword. |
| **Meta description** | 120–160 characters. Must contain the focus keyword. Plain sentence, no markdown. |
| **Excerpt** | 1–2 sentences. Same meaning as meta description but can be slightly shorter. |
| **Tags** | 4–6 tags: tech stack items + topic keywords. |
| **Category** | Always `Insights` (create if missing in WordPress). |

---

## 6. Post HTML File Structure

The `post-<slug>.html` file must contain these sections in order:

1. **Yellow instruction box** — step-by-step publishing guide for marketing
2. **Post meta table** — title, slug, focus keyword, meta description, category, tags, excerpt
3. **Thumbnail table** — file names, sizes, alt texts for all images
4. **AI rewrite prompt** — copy-paste prompt with full project context (see template below)
5. **`<!-- POST CONTENT START -->`** — WordPress-ready HTML block
6. **`<!-- POST CONTENT END -->`**

### Post Content Structure (inside the HTML block)

```
H1 — Post title
Tagline — one sentence summary

H2 — The Problem
H2 — What We Built
H2 — How It Works   ← diagram goes here
H2 — Key Technical Decisions   ← bullet list
H2 — Results / What Changed
H2 — Stack   ← monospace tag list

--- separator ---
Blockquote CTA → /contact
```

---

## 7. AI Rewrite Prompt Template

Use this template when creating the AI rewrite prompt section in `post-<slug>.html`.
Replace the `[BRACKETS]` with post-specific details.

```
You are a technical content writer for a software consultancy called Northbeam Tech (northbeam-tech.com/insights/).

Here is a blog post draft about [SHORT PROJECT DESCRIPTION]. Please rewrite or adjust it based on my instructions below.

---
CONTEXT:
- Project name: [PROJECT NAME]
- Platform / CRM: [e.g. SugarCRM, HubSpot, custom]
- Integration: [e.g. Twilio WhatsApp API, Stripe, etc.]
- What it does: [1–2 sentences describing the core feature]
- Key features: [comma-separated list]
- Tech stack: [PHP / Node / Python, database, APIs]
- Focus keyword for SEO: "[FOCUS KEYWORD]"
- Target audience: [who reads this — e.g. "education businesses using SugarCRM"]

---
CURRENT POST TITLE:
[POST TITLE]

---
CURRENT POST CONTENT:
[paste the post HTML content here]

---
MY INSTRUCTIONS:
[describe what you want to change — e.g. "make it shorter", "add a pricing section",
"make the tone more casual", "translate to Vietnamese", "add FAQ section", etc.]
```

---

## 8. Checklist Before Sending to Marketing

- [ ] Folder named with the post slug (lowercase, hyphenated)
- [ ] `post-<slug>.html` includes all 6 sections (instruction box, meta, thumbnail table, AI prompt, content)
- [ ] Focus keyword appears in: title, meta description, first paragraph, at least one H2
- [ ] Meta description is 120–160 characters
- [ ] Thumbnail PNG is 1200×628 px
- [ ] All image alt texts filled in
- [ ] Diagram `src` uses relative path (e.g. `src="<slug>-flow-diagram.png"`) so it previews locally — marketing replaces with WordPress URL after upload
- [ ] SVG sources saved in repo (not required in the marketing zip)
- [ ] Only `.png` + `.html` files sent to marketing
