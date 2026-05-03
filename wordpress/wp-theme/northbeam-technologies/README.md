# Northbeam Technologies WordPress Theme

## Install

1. Zip folder `northbeam-technologies`.
2. In WordPress Admin: Appearance -> Themes -> Add New -> Upload Theme.
3. Activate theme.

## Setup

1. Create pages: Home, About, Services, Projects, Blog, News, Contact.
2. Settings -> Reading:
   - Your homepage displays: A static page
   - Homepage: Home
   - Posts page: Blog
3. Appearance -> Menus:
   - Create menu with Home/About/Services/Projects/Blog/News/Contact
   - Assign to `Primary Menu`.

## Notes

- **Social links (footer):** Defaults match the company profile — Facebook `https://www.facebook.com/northbeam-technologies`, LinkedIn `https://www.linkedin.com/company/northbeam-technologies`, X `https://x.com/northbeam-technologies`. Override under **Appearance → Customize → Social links**.
- Theme assets are in `assets/`.
- Main logo file currently used: `assets/logo-short.svg`.
- Favicon file: `assets/favicon.svg`.
- Always update the theme version in `style.css` after any WordPress theme change.

## Local static preview

To run the static page locally:

```bash
cd wordpress/static-html
python3 -m http.server 8000
```

Open `http://localhost:8000` in your browser.
