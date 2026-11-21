# Promotion Kit — profileforge

Ready-to-post content to drive stars. Post to each platform; track with UTM if desired.

## 1. Hacker News — Show HN

Title: Show HN: profileforge – generate a GitHub profile README + dev card in one command

Body:
Hi HN, I built profileforge, a small Go CLI that turns your public GitHub data into a polished profile README and an SVG dev card in one command.

- Zero dependencies (Go stdlib only), single static binary.
- Fetches repos, stars, followers, top languages, featured projects.
- Dark/light themes with matching badges.
- Works unauthenticated (60 req/hr) or with GITHUB_TOKEN.

Install: go install github.com/v01dst/profileforge@latest
Usage: profileforge <username> -theme dark -out ./profile

Repo: https://github.com/v01dst/profileforge
Release: https://github.com/v01dst/profileforge/releases

Why: I maintain 100+ small zero-dependency CLI tools and wanted a one-command way to generate a profile README instead of hand-editing badges and stats. Feedback welcome — especially on output format and which stats matter most.

---

## 2. Reddit — r/golang

Title: [Project] profileforge — generate a GitHub profile README + dev card in one command (Go, zero deps)

Body:
I just released profileforge: a Go CLI that generates a polished GitHub profile README (stats, top languages, featured projects) and an SVG dev card from your public GitHub data.

- Zero dependencies, single static binary (~7MB)
- Dark/light themes
- Optional GITHUB_TOKEN for higher rate limits

Install: `go install github.com/v01dst/profileforge@latest`
Try: `profileforge v01dst`

Repo: https://github.com/v01dst/profileforge
Would love feedback on the output and what stats you'd want added.

---

## 3. Reddit — r/github

Title: One-command GitHub profile README + dev card generator

Body:
Tired of hand-editing profile READMEs? profileforge generates a polished README (stats, languages, top projects) and an SVG card from your public GitHub data in one command.

`profileforge <username> -theme dark`

Zero deps, single binary. Repo: https://github.com/v01dst/profileforge

---

## 4. Reddit — r/commandline

Title: profileforge — generate a GitHub profile README + dev card from your terminal

Body:
Small CLI: profileforge <username> → README.md + card.svg. Zero dependencies, Go. https://github.com/v01dst/profileforge

---

## 5. Dev.to article

Title: Generate your GitHub profile README in one command with profileforge

Tags: github, go, cli, opensource

Body:
# Generate your GitHub profile README in one command

Hand-editing a GitHub profile README is tedious: badges, stats, language bars, featured repos... profileforge does it in one command.

## Install

```bash
go install github.com/v01dst/profileforge@latest
```

## Use

```bash
profileforge v01dst -theme dark -out ./profile
```

Produces `README.md` and `card.svg`. Drop the README into a repo named after your username and GitHub renders it on your profile.

## Features

- Stats: repos, stars, followers, forks, account age
- Top languages from your actual repos
- Featured projects (auto-selected by stars)
- Dark/light themes
- SVG dev card you can embed anywhere
- Zero dependencies, single static binary
- Optional GITHUB_TOKEN

## Links

- Repo: https://github.com/v01dst/profileforge
- Releases: https://github.com/v01dst/profileforge/releases

If you find it useful, a ⭐ helps a lot!

---

## 6. Twitter / X thread

1/ I just shipped profileforge — generate a beautiful GitHub profile README + dev card in one command 🧵

2/ `profileforge v01dst` → README.md + card.svg. Stats, top languages, featured projects, badges. Zero deps, single Go binary.

3/ Dark & light themes. Optional GITHUB_TOKEN. Works unauthenticated.

4/ Install: `go install github.com/v01dst/profileforge@latest`

5/ Repo: https://github.com/v01dst/profileforge — if it's useful, a ⭐ means a lot!

---

## 7. LinkedIn

I just released profileforge — a small Go CLI that generates a polished GitHub profile README and an SVG dev card from your public GitHub data in one command.

Zero dependencies, single static binary. Dark/light themes. Optional token for higher rate limits.

Repo: https://github.com/v01dst/profileforge

If you've ever hand-edited a profile README, this is for you. Feedback welcome!

---

## 8. Product Hunt

Name: profileforge
Tagline: Generate a GitHub profile README + dev card in one command
Description: profileforge fetches your public GitHub data and writes a polished profile README (stats, top languages, featured projects, badges) and an SVG dev card. Zero dependencies, single static Go binary. Dark/light themes.
Link: https://github.com/v01dst/profileforge
Topics: Developer Tools, GitHub, Open Source

---

## 9. Directories to submit (manual)

- Terminal Trove: https://terminaltrove.com/submit  (submit profileforge)
- awesome-cli-apps: https://github.com/agarrharr/awesome-cli-apps (PR adding profileforge under "Developer Tools" / "Git")
- awesome-go: https://github.com/avelino/awesome-go (PR under "Command Line" → Standard CLI, if accepted)
- Console (console.dev): submit via their form
- Hacker News (Show HN): post above
- Lobsters: submit link

## 10. GitHub repo settings checklist

- [x] Description set
- [x] Topics set
- [x] README with badges + demo image
- [x] LICENSE (MIT)
- [x] CI workflow
- [x] Release v0.1.0 with binaries
- [ ] Social preview image (upload card-preview.svg as PNG 1280x640 in repo settings)
- [ ] Enable Discussions
- [ ] Pin repo on profile
- [ ] Add to profile README

## Note on Starstruck

Starstruck requires 16+ real stars. The above promotion drives real stars. Do NOT buy/fake stars — GitHub detects and bans. Track progress: gh repo view v01dst/profileforge --json stargazerCount