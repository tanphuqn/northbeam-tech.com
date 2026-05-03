# WordPress page exports

Paste each file’s content into the matching WordPress **Page** (Code editor or block paste). The **Northbeam Technologies** theme expects:

- **Home**: `front-page.php` wraps content in `.northbeam-front` (see **Reading** below).
- **Other pages**: `page.php` wraps content in `.northbeam-inner` (full-width blocks; title comes from your blocks, not the template).
- **Block colors**: Exports use preset classes such as `has-background-background-color` and `has-tertiary-background-color`. The theme ships **`theme.json`** with a matching color palette so WordPress outputs the correct background, text, and border colors. Without it, the home page can look flat, unreadable, or “broken.” Deploy theme **v1.0.3+** and clear any page cache after updating.

## Reading

- **Settings → Reading**: set **Your homepage displays** to **A static page**, and choose the page whose content is pasted from `home.html`.

## Slug checklist

WordPress permalinks should match in-content links (especially from `home.html`).

| Paste from        | Suggested slug   | Notes |
|-------------------|------------------|--------|
| `home.html`       | *(front page)*   | Assign as static front page. |
| `about-us.html`   | `about-us`       | Optional: fix placeholder `href="#"` buttons in the editor. |
| `contact.html`    | `contact`        | See **Contact Form 7** below. |
| `service.html`    | `service`        | Home links use `/service` (not `/services`). Use theme **v1.0.7+** so **`services-page`** body styles apply to pillar images (`.service-card`) and **How We Work** (`.how-step`). |
| `project.html`    | `project`        | Home links use `/project`. |
| `privacy-policy.html` | `privacy-policy` | Footer secondary link; **legal review** before launch. |
| `terms.html`      | `terms`          | Footer secondary link; **legal review** before launch. |
| `careers.html`    | `careers`        | Footer secondary link. |
| `sitemap.html`    | `sitemap`        | Human-readable site map; footer secondary link. |
| `support.html`    | `support`        | Footer secondary link. |
| `insights.html`   | `insights`       | **Query Loop** listing blog posts; see **Insights page (Query Loop)** below. |

### Insights page (Query Loop)

Use this when you want an **Insights** menu item pointing to a **Page** at `/insights/` that lists **Posts** (optionally filtered to one category).

1. **Posts → Categories**: create **Insights** (slug `insights`). Assign posts to this category.
2. **Pages → Add New**, slug **`insights`**. Switch to **Code editor** (⋮ → **Code editor**).
3. Paste the full contents of **`insights.html`** (including the HTML comment at the top — WordPress ignores it).
4. Switch to **Visual editor** or **List view** → select the **Query Loop** block (parent of “Post template”).
5. In the **block sidebar**, open **Filters** / **Taxonomies** (wording depends on WP version) and set **Category** to **Insights** so the list is not “all posts”. If you skip this step, the loop shows every post.
6. **Publish** and add **Insights** to **Appearance → Menus** as this page.

**Footer note:** The theme’s default footer link for “Insights” uses **`/insights/`** (theme v1.2.1+). If you use a different URL, assign a **Footer Menu** or filter **`northbeam_footer_fallback_links`**.

### Create the legal & support pages (quick steps)

1. **Pages → Add New** for each page below (or **Add** five pages in one go from **Pages**).
2. Set the **permalink slug** in the sidebar to match the **Suggested slug** column (WordPress may auto-title from the H1 after paste — you can still edit the slug).
3. Switch the editor to **Code editor** (⋮ menu → **Code editor**) or paste into a **Custom HTML** block if needed, then paste the full contents of the matching `wp-pages/*.html` file.
4. **Publish**. The theme footer already links to `/privacy-policy/`, `/terms/`, `/careers/`, `/sitemap/`, and `/support/`.

If you use different slugs, update links inside `home.html` (or edit those links in the block editor after paste), and use the **`northbeam_footer_secondary_links`** filter if URLs must change.

## Contact Form 7 (`contact.html`)

The export embeds a CF7 block with **example** metadata from the source site:

- Block JSON may reference `id` **40** and `hash` **`a56ab6e`**.
- The shortcode line looks like: `[contact-form-7 id="a56ab6e" title="Contact form 1"]`.

On your live site, form IDs differ. After pasting:

1. Install **Contact Form 7** and create (or reuse) your form.
2. In the page editor, **select the Contact Form 7 block** and choose the correct form from the dropdown, or replace the shortcode with your form’s **id** from **Contact → Contact Forms** in wp-admin.

Until the block points at a real form, the contact page may show nothing or an error.

## Footer menu

The theme registers **Footer Menu**. If you **do not** assign a menu, the footer still shows **default links** (About Us, Our Services, Projects, Insights, Contact Us) using the URLs below.

**Create a custom footer menu:** **Appearance → Menus → Create a new menu** → add **Custom Links** or **Pages** → check **Footer Menu** under **Menu Settings** → **Save Menu**.

### Pages to include (recommended)

| Priority | Page | Suggested slug | Role |
|----------|------|----------------|------|
| Core | Home | *(front page)* | Linked from logo; optional extra footer link |
| Core | About Us | `about-us` | Company story, team |
| Core | Our Services | `service` | Offerings (paste `service.html`) |
| Core | Projects | `project` | Portfolio (paste `project.html`) |
| Core | Insights | `insights` | Use **`insights.html`** (Query Loop) or a category archive; default footer uses **`/insights/`** |
| Core | Contact Us | `contact` | Form + details (paste `contact.html`) |
| Optional | Privacy Policy | `privacy-policy` | Legal / cookies / GDPR |
| Optional | Terms of Service | `terms` | Legal |
| Optional | Careers | `careers` | Jobs (footer secondary link) |
| Optional | Site map (HTML) | `sitemap` | Human-readable list of pages (WordPress also exposes an XML sitemap; see **SEO** plugins or `/wp-sitemap.xml` in WP 5.5+) |
| Optional | Support | `support` | Help for clients (footer secondary link) |

Adjust slugs in **Settings → Permalinks** and edit page slugs so they match in-content links from `home.html` (`/service`, `/project`, `/contact`, etc.). To change default footer URLs without a menu, use the **`northbeam_footer_fallback_links`** filter in a child theme or small plugin.

### Footer: secondary row (Privacy, Terms, Careers, Sitemap, Support)

The theme shows a second row in the footer bar (next to the copyright) with links to **`/privacy-policy/`**, **`/terms/`**, **`/careers/`**, **`/sitemap/`**, and **`/support/`**. Create those **Pages** in WordPress with matching slugs (or override the list).

Developers can replace or hide the row with the **`northbeam_footer_secondary_links`** filter: pass an empty array to hide it, or an array of `url` / `label` pairs.

Static HTML mirrors: `northbeam-tech/static-html/privacy.html`, `terms.html`, `careers.html`, `sitemap.html`, `support.html`.

### Footer: Facebook, LinkedIn & X

**Appearance → Customize → Social links** sets the **Facebook**, **LinkedIn**, and **X** URLs shown under “Follow us” in the footer (defaults: Facebook `northbeam-technologies`, LinkedIn `company/northbeam-technologies`, X `x.com/northbeam-technologies`). Clear a field to hide that link. Developers can override with filters **`northbeam_facebook_url`**, **`northbeam_linkedin_url`**, and **`northbeam_x_url`**.
