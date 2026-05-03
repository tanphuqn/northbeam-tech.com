# Northbeam Technologies

## Project Structure

- `wordpress/wp-theme/northbeam-technologies/` — WordPress theme source files
- `wordpress/static-html/` — Static HTML version of the site for local preview and demo
- `wordpress/wp-pages/` — WordPress page content copies and templates
- `wordpress/site.info` — WordPress environment snapshot

## How to run the static page

From the repository root:

```bash
cd wordpress/static-html
python3 -m http.server 8000
```

Then open:

```text
http://localhost:8000
```

If your machine uses `python` for Python 3, use:

```bash
cd wordpress/static-html
python -m http.server 8000
```

Or with Node if you prefer:

```bash
cd wordpress/static-html
npx http-server . -p 8000
```

## WordPress theme updates

The WordPress theme source is in:

```text
wordpress/wp-theme/northbeam-technologies/
```

To install or update the theme:

1. Zip the `northbeam-technologies` folder.
2. In WordPress Admin, go to **Appearance → Themes → Add New → Upload Theme**.
3. Activate the theme.

## Notes

- The static preview and WordPress theme share the same design and CSS logic.
- Use the static preview for quick frontend checks without running a full WordPress install.
- For theme-specific setup, see `wordpress/wp-theme/northbeam-technologies/README.md`.
