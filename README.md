# spinney.dev

The home page of spinney, a small lab: an animation of a circle rolling
along a line, drawing cycloids, then the name, a line about the lab, and
two links, projects and contact, each its own page under the same
animation, which carries on from page to page.

The words are markdown, and a small Go program sets them in one layout:

    content/            one .md file per page: home, projects, contact
    templates/          layout.html, around every page
    main.go             the build: content/*.md into public/
    public/lab.css      the look
    public/pgp.asc      the PGP key for hello@spinney.dev, linked from contact
    public/js/          the animation: cycloid-horizontal.js and loop.js,
                        copied from obk.dev; change them there and copy
                        them back
    public/             the site; index.html, projects/ and contact/ are built

A page starts with its title between two `---` lines:

    ---
    title: projects
    ---

    - [canopy](https://canopy.spinney.dev/): a small, read-only web viewer for git

`content/home.md` is the home page; any other `content/<name>.md` is
served at `/<name>/`. To add a project, add a line to
`content/projects.md`.

To try it:

    go run . -serve     # builds, then serves on 127.0.0.1:8000 and
                        # rebuilds on every page load

Cloudflare Pages is connected to this repository: on every push to main
it runs `go run .` and serves `public/`. The built pages are not
committed.

The site for canopy, the git viewer, is a separate repository, served at
canopy.spinney.dev.

## License

[GNU GPL v3](LICENSE), or any later version.

    spinney.dev, the spinney home page and the program that builds it.
    Copyright (C) 2026 The spinney authors

    This program is free software: you can redistribute it and/or modify it
    under the terms of the GNU General Public License as published by the
    Free Software Foundation, either version 3 of the License, or (at your
    option) any later version.

    This program is distributed in the hope that it will be useful, but
    WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
    Public License for more details.
