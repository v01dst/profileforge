package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

var version = "0.1.0"

// ---------- GitHub API types ----------

type ghUser struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	Bio         string `json:"bio"`
	Location    string `json:"location"`
	Blog        string `json:"blog"`
	Company     string `json:"company"`
	AvatarURL   string `json:"avatar_url"`
	HTMLURL     string `json:"html_url"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	CreatedAt   string `json:"created_at"`
}

type ghRepo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	Stargazers  int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Archived    bool   `json:"archived"`
	Fork        bool   `json:"fork"`
	Private     bool   `json:"private"`
	UpdatedAt   string `json:"updated_at"`
	Topics      []string `json:"topics"`
}

// ---------- HTTP helper ----------

func ghGet(path string) ([]byte, error) {
	url := "https://api.github.com" + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "profileforge/"+version)
	// optional token
	if t := os.Getenv("GITHUB_TOKEN"); t != "" {
		req.Header.Set("Authorization", "Bearer "+t)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("github api %s: %s (status %d)", path, string(body), resp.StatusCode)
	}
	return body, nil
}

// ---------- Data collection ----------

type Profile struct {
	User        ghUser
	Repos       []ghRepo
	TopRepos    []ghRepo
	Languages   map[string]int // language -> bytes/weight
	TotalStars  int
	TotalForks  int
}

func collect(username string) (*Profile, error) {
	p := &Profile{Languages: map[string]int{}}

	// user
	b, err := ghGet("/users/" + username)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &p.User); err != nil {
		return nil, err
	}

	// repos (paginated, up to 300)
	page := 1
	for page <= 3 {
		b, err := ghGet(fmt.Sprintf("/users/%s/repos?per_page=100&sort=updated&page=%d&type=owner", username, page))
		if err != nil {
			return nil, err
		}
		var repos []ghRepo
		if err := json.Unmarshal(b, &repos); err != nil {
			return nil, err
		}
		if len(repos) == 0 {
			break
		}
		for _, r := range repos {
			if r.Fork || r.Private || r.Archived {
				continue
			}
			p.Repos = append(p.Repos, r)
			p.TotalStars += r.Stargazers
			p.TotalForks += r.Forks
			if r.Language != "" {
				p.Languages[r.Language] += 1
			}
		}
		if len(repos) < 100 {
			break
		}
		page++
	}

	// top repos by stars (then updated)
	sorted := append([]ghRepo(nil), p.Repos...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Stargazers != sorted[j].Stargazers {
			return sorted[i].Stargazers > sorted[j].Stargazers
		}
		return sorted[i].UpdatedAt > sorted[j].UpdatedAt
	})
	if len(sorted) > 6 {
		sorted = sorted[:6]
	}
	p.TopRepos = sorted

	return p, nil
}

// ---------- Helpers ----------

func displayName(u ghUser) string {
	if u.Name != "" {
		return u.Name
	}
	return u.Login
}

func langColor(lang string) string {
	colors := map[string]string{
		"Go":         "#00ADD8",
		"Python":     "#3572A5",
		"JavaScript": "#f1e05a",
		"TypeScript": "#3178c6",
		"Rust":       "#dea584",
		"Java":       "#b07219",
		"C":          "#555555",
		"C++":        "#f34b7d",
		"C#":         "#178600",
		"Ruby":       "#701516",
		"PHP":        "#4F5D95",
		"Swift":      "#F05138",
		"Kotlin":     "#A97BFF",
		"Shell":      "#89e051",
		"HTML":       "#e34c26",
		"CSS":        "#563d7c",
		"Vue":        "#41b883",
		"Svelte":     "#ff3e00",
		"Dart":       "#00B4AB",
		"Lua":        "#000080",
		"Zig":        "#ec915c",
		"Elixir":     "#6e4a7e",
		"Haskell":    "#5e5086",
	}
	if c, ok := colors[lang]; ok {
		return c
	}
	return "#8b949e"
}

func topLanguages(langs map[string]int, n int) []struct {
	Name  string
	Count int
} {
	type kv struct {
		Name  string
		Count int
	}
	var arr []kv
	for k, v := range langs {
		arr = append(arr, kv{k, v})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].Count > arr[j].Count })
	if len(arr) > n {
		arr = arr[:n]
	}
	var out []struct {
		Name  string
		Count int
	}
	for _, a := range arr {
		out = append(out, struct {
			Name  string
			Count int
		}{a.Name, a.Count})
	}
	return out
}

// ---------- README generation ----------

func genREADME(p *Profile, theme string) string {
	u := p.User
	var b strings.Builder

	b.WriteString("# Hi, I'm " + displayName(u) + " 👋\n\n")

	if u.Bio != "" {
		bio := strings.ReplaceAll(u.Bio, "\n", " · ")
		bio = strings.ReplaceAll(bio, "\r", "")
		b.WriteString("> " + bio + "\n\n")
	}

	// badges row
	b.WriteString("![Profile views](https://komarev.com/ghpvc/?username=" + u.Login + "&style=flat-square)\n\n")

	// about
	b.WriteString("## About\n\n")
	if u.Company != "" {
		b.WriteString("- 🏢 " + u.Company + "\n")
	}
	if u.Location != "" {
		b.WriteString("- 📍 " + u.Location + "\n")
	}
	if u.Blog != "" {
		blog := u.Blog
		if !strings.HasPrefix(blog, "http") {
			blog = "https://" + blog
		}
		b.WriteString("- 🔗 " + blog + "\n")
	}
	repoWord := "repos"
		if u.PublicRepos == 1 { repoWord = "repo" }
		b.WriteString("- 📦 " + fmt.Sprintf("%d", u.PublicRepos) + " public " + repoWord + "\n")
	b.WriteString("- ⭐ " + fmt.Sprintf("%d", p.TotalStars) + " stars received\n")
	b.WriteString("- 👥 " + fmt.Sprintf("%d", u.Followers) + " followers\n")
	if u.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, u.CreatedAt); err == nil {
			b.WriteString("- 🗓️ On GitHub since " + t.Format("January 2006") + "\n")
		}
	}
	b.WriteString("\n")

	// stats
	b.WriteString("## GitHub Stats\n\n")
	b.WriteString("![GitHub stats](https://github-readme-stats.vercel.app/api?username=" + u.Login + "&show_icons=true&theme=" + theme + "&hide_border=true)\n\n")
	b.WriteString("![Top languages](https://github-readme-stats.vercel.app/api/top-langs/?username=" + u.Login + "&layout=compact&theme=" + theme + "&hide_border=true)\n\n")
	b.WriteString("![Streak](https://github-readme-streak-stats.herokuapp.com/?user=" + u.Login + "&theme=" + theme + "&hide_border=true)\n\n")

	// top languages (from own data)
	if len(p.Languages) > 0 {
		b.WriteString("## Top Languages\n\n")
		for _, l := range topLanguages(p.Languages, 8) {
			b.WriteString(fmt.Sprintf("- **%s** — %d repos\n", l.Name, l.Count))
		}
		b.WriteString("\n")
	}

	// top repos
	if len(p.TopRepos) > 0 {
		b.WriteString("## Featured Projects\n\n")
		for _, r := range p.TopRepos {
			desc := r.Description
			if desc == "" {
				desc = "No description"
			}
			b.WriteString(fmt.Sprintf("- [**%s**](%s) — %s ⭐%d\n", r.Name, r.HTMLURL, desc, r.Stargazers))
		}
		b.WriteString("\n")
	}

	b.WriteString("---\n\n")
	b.WriteString("*Generated with [profileforge](https://github.com/v01dst/profileforge)*\n")

	return b.String()
}

// ---------- SVG card ----------

func genCard(p *Profile, theme string) string {
	u := p.User
	name := displayName(u)
	bio := u.Bio
	if len(bio) > 60 {
		bio = bio[:57] + "..."
	}

	bg, fg, muted, accent := "#0d1117", "#e6edf3", "#8b949e", "#58a6ff"
	if theme == "light" {
		bg, fg, muted, accent = "#ffffff", "#1f2328", "#656d76", "#0969da"
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="420" height="200" viewBox="0 0 420 200" role="img" aria-label="%s">
`, name))
	b.WriteString(fmt.Sprintf(`<rect width="420" height="200" rx="12" fill="%s" stroke="#30363d"/>
`, bg))
	// avatar circle
	b.WriteString(fmt.Sprintf(`<circle cx="50" cy="50" r="32" fill="%s"/>
`, accent))
	b.WriteString(fmt.Sprintf(`<text x="50" y="58" text-anchor="middle" font-family="system-ui,Segoe UI,Roboto" font-size="24" font-weight="700" fill="#ffffff">%s</text>
`, strings.ToUpper(string([]rune(name)[0:1]))))
	// name
	b.WriteString(fmt.Sprintf(`<text x="95" y="40" font-family="system-ui,Segoe UI,Roboto" font-size="18" font-weight="700" fill="%s">%s</text>
`, fg, escapeXML(name)))
	if bio != "" {
		b.WriteString(fmt.Sprintf(`<text x="95" y="62" font-family="system-ui,Segoe UI,Roboto" font-size="12" fill="%s">%s</text>
`, muted, escapeXML(bio)))
	}
	// stats
	stats := []struct {
		Label string
		Value int
	}{
		{"Repos", u.PublicRepos},
		{"Stars", p.TotalStars},
		{"Followers", u.Followers},
		{"Forks", p.TotalForks},
	}
	x := 20
	for _, s := range stats {
		b.WriteString(fmt.Sprintf(`<text x="%d" y="115" font-family="system-ui,Segoe UI,Roboto" font-size="20" font-weight="700" fill="%s">%d</text>
`, x, fg, s.Value))
		b.WriteString(fmt.Sprintf(`<text x="%d" y="132" font-family="system-ui,Segoe UI,Roboto" font-size="11" fill="%s">%s</text>
`, x, muted, s.Label))
		x += 100
	}
	// language bar
	if len(p.Languages) > 0 {
		tl := topLanguages(p.Languages, 5)
		total := 0
		for _, l := range tl {
			total += l.Count
		}
		bx := 20
		bw := 380
		cur := bx
		for _, l := range tl {
			w := int(float64(l.Count) / float64(total) * float64(bw))
			b.WriteString(fmt.Sprintf(`<rect x="%d" y="155" width="%d" height="8" fill="%s"/>
`, cur, w, langColor(l.Name)))
			cur += w
		}
		// legend
		lx := 20
		for _, l := range tl {
			b.WriteString(fmt.Sprintf(`<circle cx="%d" cy="180" r="4" fill="%s"/>
`, lx, langColor(l.Name)))
			b.WriteString(fmt.Sprintf(`<text x="%d" y="184" font-family="system-ui,Segoe UI,Roboto" font-size="10" fill="%s">%s</text>
`, lx+8, muted, escapeXML(l.Name)))
			lx += 78
		}
	}
	b.WriteString("</svg>")
	return b.String()
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// ---------- main ----------

func usage() {
	fmt.Println(`profileforge — generate a beautiful GitHub profile README + dev card

Usage:
  profileforge <username> [flags]

Flags:
  -out string     output directory (default ".")
  -theme string   theme: dark | light (default "dark")
  -card           also generate an SVG dev card (default true)
  -readme         generate README.md (default true)
  -token string   GitHub token (or set GITHUB_TOKEN)
  -version        print version

Examples:
  profileforge torvalds
  profileforge v01dst -theme light -out ./out
  GITHUB_TOKEN=ghp_xxx profileforge octocat`)
}

func main() {
	outDir := flag.String("out", ".", "output directory")
	theme := flag.String("theme", "dark", "theme: dark | light")
	genCardFlag := flag.Bool("card", true, "generate SVG dev card")
	genReadme := flag.Bool("readme", true, "generate README.md")
	showVersion := flag.Bool("version", false, "print version")
	flag.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Println("profileforge", version)
		return
	}
	if flag.NArg() < 1 {
		usage()
		os.Exit(2)
	}
	username := flag.Arg(0)

	fmt.Printf("🔍 Fetching profile for @%s...\n", username)
	p, err := collect(username)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(*outDir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if *genReadme {
		md := genREADME(p, *theme)
		path := *outDir + "/README.md"
		if err := os.WriteFile(path, []byte(md), 0644); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Wrote", path)
	}
	if *genCardFlag {
		svg := genCard(p, *theme)
		path := *outDir + "/card.svg"
		if err := os.WriteFile(path, []byte(svg), 0644); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Wrote", path)
	}

	fmt.Printf("🎉 Done. %d repos, %d stars, %d languages.\n", len(p.Repos), p.TotalStars, len(p.Languages))
}
