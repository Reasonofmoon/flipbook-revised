package handlers

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"time"

	"github.com/jonradoff/flipbook/internal/database"
	"github.com/jonradoff/flipbook/internal/i18n"
	"github.com/jonradoff/flipbook/internal/models"
	"github.com/jonradoff/flipbook/internal/stats"
)

// StatsHandler serves the operations dashboard and per-academy view reports.
type StatsHandler struct {
	db      *database.DB
	tmpl    *template.Template
	baseURL string
	dataDir string
	uiLang  string
}

func NewStatsHandler(db *database.DB, tmpl *template.Template, baseURL, dataDir, uiLang string) *StatsHandler {
	return &StatsHandler{db: db, tmpl: tmpl, baseURL: baseURL, dataDir: dataDir, uiLang: uiLang}
}

// unassignedClient is the query value for flipbooks with no academy set.
const unassignedClient = "-"

// bar is one day in the daily chart.
type bar struct {
	Date  time.Time
	Count int
	Pct   int // height as a percentage of the busiest day
}

type refRow struct {
	Name  string // "" means direct visit
	Count int
	Pct   int // width relative to the top referrer
}

type clientRow struct {
	Name      string // "" means unassigned
	Key       string // query value for the report link
	Flipbooks int
	Views     int
}

type flipbookRow struct {
	FB    *models.Flipbook
	Views int
}

func (h *StatsHandler) lang(r *http.Request) string {
	return i18n.ResolveAdmin(h.uiLang, r.Header.Get("Accept-Language"))
}

func (h *StatsHandler) selfHost() string {
	if u, err := url.Parse(h.baseURL); err == nil {
		return u.Hostname()
	}
	return ""
}

// load fetches flipbooks and the period's views, already summarized.
func (h *StatsHandler) load(ctx context.Context, r *http.Request, ids []string) (stats.Period, stats.Summary, error) {
	t := i18n.Admin(h.lang(r))
	p := stats.ParsePeriod(r.URL.Query().Get("period"), time.Now())
	views, err := h.db.ListViews(ctx, ids, p.From, p.To)
	if err != nil {
		return p, stats.Summary{}, err
	}
	return p, stats.Summarize(views, p, h.selfHost(), t["ref_self"], 8), nil
}

func (h *StatsHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	t := i18n.Admin(h.lang(r))
	flipbooks, err := h.db.ListFlipbooks()
	if err != nil {
		http.Error(w, t["err_list_failed"], http.StatusInternalServerError)
		return
	}
	p, sum, err := h.load(ctx, r, nil)
	if err != nil {
		log.Printf("stats: list views: %v", err)
		http.Error(w, t["err_stats_failed"], http.StatusInternalServerError)
		return
	}

	byClient := map[string]*clientRow{}
	var top []flipbookRow
	for _, fb := range flipbooks {
		key := fb.Client
		if key == "" {
			key = unassignedClient
		}
		row, ok := byClient[key]
		if !ok {
			row = &clientRow{Name: fb.Client, Key: key}
			byClient[key] = row
		}
		row.Flipbooks++
		row.Views += sum.ByFlipbook[fb.ID]
		if v := sum.ByFlipbook[fb.ID]; v > 0 {
			top = append(top, flipbookRow{FB: fb, Views: v})
		}
	}
	clients := make([]clientRow, 0, len(byClient))
	for _, c := range byClient {
		clients = append(clients, *c)
	}
	sort.Slice(clients, func(i, j int) bool {
		if clients[i].Views != clients[j].Views {
			return clients[i].Views > clients[j].Views
		}
		return clients[i].Name < clients[j].Name
	})
	sortFlipbookRows(top)
	if len(top) > 10 {
		top = top[:10]
	}

	dbBytes, err := h.db.StorageUsage(ctx)
	if err != nil {
		log.Printf("stats: db storage: %v", err)
	}

	h.tmpl.ExecuteTemplate(w, "admin_stats", map[string]interface{}{
		"Lang":      h.lang(r),
		"T":         t,
		"Period":    p,
		"Periods":   periodKeys,
		"PeriodTxt": periodText(p),
		"Summary":   sum,
		"Bars":      bars(sum),
		"Refs":      refRows(sum),
		"Clients":   clients,
		"Top":       top,
		"Books":     len(flipbooks),
		"DBSize":    humanBytes(dbBytes),
		"FileSize":  humanBytes(dirSize(h.dataDir)),
	})
}

func (h *StatsHandler) ClientReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	t := i18n.Admin(h.lang(r))
	key := r.URL.Query().Get("name")
	client := key
	if key == unassignedClient {
		client = ""
	}
	flipbooks, err := h.db.ListFlipbooks()
	if err != nil {
		http.Error(w, t["err_list_failed"], http.StatusInternalServerError)
		return
	}
	var mine []*models.Flipbook
	ids := []string{} // non-nil: an academy with no flipbooks matches no views
	for _, fb := range flipbooks {
		if fb.Client == client {
			mine = append(mine, fb)
			ids = append(ids, fb.ID)
		}
	}
	if len(mine) == 0 {
		http.Error(w, t["err_client_not_found"], http.StatusNotFound)
		return
	}
	p, sum, err := h.load(ctx, r, ids)
	if err != nil {
		log.Printf("stats: list views: %v", err)
		http.Error(w, t["err_stats_failed"], http.StatusInternalServerError)
		return
	}
	rows := make([]flipbookRow, 0, len(mine))
	for _, fb := range mine {
		rows = append(rows, flipbookRow{FB: fb, Views: sum.ByFlipbook[fb.ID]})
	}
	sortFlipbookRows(rows)

	h.tmpl.ExecuteTemplate(w, "admin_report", map[string]interface{}{
		"Lang":      h.lang(r),
		"T":         t,
		"Client":    client,
		"Key":       key,
		"Period":    p,
		"Periods":   periodKeys,
		"PeriodTxt": periodText(p),
		"Issued":    time.Now().In(stats.KST).Format("2006.01.02"),
		"Summary":   sum,
		"Bars":      bars(sum),
		"Refs":      refRows(sum),
		"Rows":      rows,
		"BaseURL":   h.baseURL,
	})
}

func sortFlipbookRows(rows []flipbookRow) {
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].Views != rows[j].Views {
			return rows[i].Views > rows[j].Views
		}
		return rows[i].FB.CreatedAt.After(rows[j].FB.CreatedAt)
	})
}

func bars(s stats.Summary) []bar {
	max := s.MaxDaily()
	out := make([]bar, 0, len(s.Daily))
	for _, d := range s.Daily {
		out = append(out, bar{Date: d.Date, Count: d.Count, Pct: d.Count * 100 / max})
	}
	return out
}

func refRows(s stats.Summary) []refRow {
	if len(s.Referrers) == 0 {
		return nil
	}
	top := s.Referrers[0].Count
	out := make([]refRow, 0, len(s.Referrers))
	for _, ref := range s.Referrers {
		out = append(out, refRow{Name: ref.Name, Count: ref.Count, Pct: ref.Count * 100 / top})
	}
	return out
}

// periodText renders the inclusive date range, e.g. "2026.09.12 – 2026.10.11".
func periodText(p stats.Period) string {
	return p.From.Format("2006.01.02") + " – " + p.To.AddDate(0, 0, -1).Format("2006.01.02")
}

// dirSize sums file sizes under dir (page images, thumbnails, originals).
func dirSize(dir string) int64 {
	var total int64
	filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if info, err := d.Info(); err == nil {
			total += info.Size()
		}
		return nil
	})
	return total
}

func humanBytes(n int64) string {
	switch {
	case n >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(n)/(1<<30))
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%d KB", n>>10)
	default:
		return fmt.Sprintf("%d B", n)
	}
}

// periodKeys lists the selectable report periods in display order.
var periodKeys = []string{"30d", "this-month", "last-month"}
