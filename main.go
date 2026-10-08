//
// Inkcopi - Inkscape Color Picker
// - program pemilih warna berbasis TUI yang mengambil palet warna dari Inkscape
//
// Oleh: Mongkee Lutfi <mongkee.lutfi@gmail.com>
//       program ditulis bersama dengan Gemini AI.
//
// July 2026

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Struktur Data Warna
type Color struct {
	Name    string
	R, G, B int
}

// Struktur Data Dokumen Palet
type Palette struct {
	Name   string
	Colors []Color
}

// Konversi Warna: RGB ke CMYK
func rgbToCmyk(r, g, b int) (c, m, y, k int) {
	rf, gf, bf := float64(r)/255.0, float64(g)/255.0, float64(b)/255.0
	max := math.Max(rf, math.Max(gf, bf))
	k = int(math.Round((1.0 - max) * 100))
	if k == 100 {
		return 0, 0, 0, 100
	}
	c = int(math.Round(((1.0 - rf - (float64(k)/100.0)) / (1.0 - (float64(k)/100.0))) * 100))
	m = int(math.Round(((1.0 - gf - (float64(k)/100.0)) / (1.0 - (float64(k)/100.0))) * 100))
	y = int(math.Round(((1.0 - bf - (float64(k)/100.0)) / (1.0 - (float64(k)/100.0))) * 100))
	return
}

// Konversi Warna: CMYK ke RGB
func cmykToRgb(c, m, y, k int) (r, g, b int) {
	cf, mf, yf, kf := float64(c)/100.0, float64(m)/100.0, float64(y)/100.0, float64(k)/100.0
	r = int(math.Round(255 * (1 - cf) * (1 - kf)))
	g = int(math.Round(255 * (1 - mf) * (1 - kf)))
	b = int(math.Round(255 * (1 - yf) * (1 - kf)))
	return
}

// Parser File .gpl Inkscape
func loadPalettes() []Palette {
	var list []Palette
	inkscapePath := "/usr/share/inkscape/palettes/*.gpl"
	files, _ := filepath.Glob(inkscapePath)

	if len(files) == 0 {
		files, _ = filepath.Glob("palettes/*.gpl")
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			continue
		}
		
		p := Palette{
			Name:   strings.TrimSuffix(filepath.Base(file), ".gpl"), 
			Colors: []Color{},
		}
		
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "GIMP") || strings.HasPrefix(line, "Name:") || strings.HasPrefix(line, "Columns:") || strings.HasPrefix(line, "#") {
				continue
			}
			
			var r, g, b int
			_, err := fmt.Sscanf(line, "%d %d %d", &r, &g, &b)
			if err == nil {
				parts := strings.Fields(line)
				var name string
				if len(parts) > 3 {
					name = strings.Join(parts[3:], " ")
				} else {
					name = fmt.Sprintf("RGB(%d,%d,%d)", r, g, b)
				}
				p.Colors = append(p.Colors, Color{Name: name, R: r, G: g, B: b})
			}
		}
		f.Close()

		if len(p.Colors) > 0 {
			list = append(list, p)
		}
	}

	if len(list) == 0 {
		list = append(list, Palette{
			Name: "Fallback (No Inkscape Installed)",
			Colors: []Color{
				{"Black", 0, 0, 0}, {"White", 255, 255, 255}, {"Red", 255, 0, 0}, {"Green", 0, 255, 0}, {"Blue", 0, 0, 255},
			},
		})
	}
	return list
}

type model struct {
	mode          int 
	palettes      []Palette
	currentPalIdx int
	colorIndex    int
	rgbCursor     int
	cmykCursor    int
	manualR, manualG, manualB int
	manualC, manualM, manualY, manualK int
	
	// Menyimpan ukuran terminal dinamis
	termWidth     int
	termHeight    int
	calculatedCols int // Jumlah kolom adaptif
}

func initialModel() model {
	pals := loadPalettes()
	initColor := pals[0].Colors[0]
	c, m, y, k := rgbToCmyk(initColor.R, initColor.G, initColor.B)
	return model{
		mode:           0,
		palettes:       pals,
		currentPalIdx:  0,
		colorIndex:     0,
		manualR:        initColor.R, manualG: initColor.G, manualB: initColor.B,
		manualC:        c, manualM: m, manualY: y, manualK: k,
		termWidth:      80, // Default sebelum di-update sistem
		termHeight:     24,
		calculatedCols: 10,
	}
}

func (m model) Init() tea.Cmd { 
	return nil 
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currPalette := m.palettes[m.currentPalIdx]

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Menangkap ukuran terminal saat pertama dibuka atau saat di-resize (di-drag) oleh user
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		
		// " █ " memakan 3 karakter space di terminal, kita hitung berapa kolom yang muat di layar
		m.calculatedCols = msg.Width / 3
		if m.calculatedCols < 1 {
			m.calculatedCols = 1
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.mode = (m.mode + 1) % 3
			return m, nil
		case "p":
			if m.mode == 0 {
				m.currentPalIdx = (m.currentPalIdx + 1) % len(m.palettes)
				m.colorIndex = 0
				m.syncManual(m.palettes[m.currentPalIdx].Colors[0])
			}
		}

		cols := m.calculatedCols

		switch m.mode {
		case 0:
			switch msg.String() {
			case "left", "h":
				if m.colorIndex > 0 { m.colorIndex-- }
			case "right", "l":
				if m.colorIndex < len(currPalette.Colors)-1 { m.colorIndex++ }
			case "up", "k":
				if m.colorIndex >= cols { m.colorIndex -= cols }
			case "down", "j":
				if m.colorIndex+cols < len(currPalette.Colors) { m.colorIndex += cols }
			}
			m.syncManual(currPalette.Colors[m.colorIndex])

		case 1:
			switch msg.String() {
			case "up", "k":
				if m.rgbCursor > 0 { m.rgbCursor-- }
			case "down", "j":
				if m.rgbCursor < 2 { m.rgbCursor++ }
			case "right", "l": m.adjustRGB(1)
			case "left", "h": m.adjustRGB(-1)
			case "shift+right", "L": m.adjustRGB(10)
			case "shift+left", "H": m.adjustRGB(-10)
			}

		case 2:
			switch msg.String() {
			case "up", "k":
				if m.cmykCursor > 0 { m.cmykCursor-- }
			case "down", "j":
				if m.cmykCursor < 3 { m.cmykCursor++ }
			case "right", "l": m.adjustCMYK(1)
			case "left", "h": m.adjustCMYK(-1)
			case "shift+right", "L": m.adjustCMYK(10)
			case "shift+left", "H": m.adjustCMYK(-10)
			}
		}
	}
	return m, nil
}

func (m *model) syncManual(col Color) {
	m.manualR, m.manualG, m.manualB = col.R, col.G, col.B
	m.manualC, m.manualM, m.manualY, m.manualK = rgbToCmyk(col.R, col.G, col.B)
}

func (m *model) adjustRGB(amount int) {
	switch m.rgbCursor {
	case 0: m.manualR = clamp(m.manualR+amount, 0, 255)
	case 1: m.manualG = clamp(m.manualG+amount, 0, 255)
	case 2: m.manualB = clamp(m.manualB+amount, 0, 255)
	}
	m.manualC, m.manualM, m.manualY, m.manualK = rgbToCmyk(m.manualR, m.manualG, m.manualB)
}

func (m *model) adjustCMYK(amount int) {
	switch m.cmykCursor {
	case 0: m.manualC = clamp(m.manualC+amount, 0, 100)
	case 1: m.manualM = clamp(m.manualM+amount, 0, 100)
	case 2: m.manualY = clamp(m.manualY+amount, 0, 100)
	case 3: m.manualK = clamp(m.manualK+amount, 0, 100)
	}
	m.manualR, m.manualG, m.manualB = cmykToRgb(m.manualC, m.manualM, m.manualY, m.manualK)
}

func clamp(val, min, max int) int {
	if val < min { return min }
	if val > max { return max }
	return val
}

func (m model) View() string {
	// 1. Pengaturan Tema Warna Utama (LipGloss)
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFF")).Background(lipgloss.Color("#1A1A2E")).Padding(0, 1)
	activeTabStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00")).Underline(true)
	inactiveTabStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	
	currPalette := m.palettes[m.currentPalIdx]
	
	var currentR, currentG, currentB int
	var currentName string
	
	if m.mode == 0 {
		selectedColor := currPalette.Colors[m.colorIndex]
		currentR = selectedColor.R
		currentG = selectedColor.G
		currentB = selectedColor.B
		currentName = selectedColor.Name
	} else {
		currentR, currentG, currentB = m.manualR, m.manualG, m.manualB
		currentName = "Custom User Color"
	}
	
	currentC, currentM, currentY, currentK := rgbToCmyk(currentR, currentG, currentB)
	hexStr := fmt.Sprintf("#%02X%02X%02X", currentR, currentG, currentB)

	tabs := []string{"[1] Inkscape Palette", "[2] Manual RGB", "[3] Manual CMYK"}
	for i, t := range tabs {
		if m.mode == i { tabs[i] = activeTabStyle.Render(t) } else { tabs[i] = inactiveTabStyle.Render(t) }
	}
	
	header := titleStyle.Render(" INKCOPI: INKSCAPE COLOR PICKER ") + fmt.Sprintf("  (Palet: %s)\n\n", currPalette.Name) + tabs[0] + "  " + tabs[1] + "  " + tabs[2] + "\n\n"

	// Hitung tinggi yang terpakai oleh komponen tetap (Header, Info Panel, Footer)
	// Header + Tabs memakan sekitar 4 baris, Preview memakan 4 baris, Footer memakan 2 baris. Total ~10 baris.
	fixedLines := 10
	availableContentHeight := m.termHeight - fixedLines
	if availableContentHeight < 3 {
		availableContentHeight = 3 // Batas aman minimum jika terminal terlalu kecil
	}

	var content string
	switch m.mode {
	case 0:
		content = "Ganti Palet: [P] | Navigasi: [Arrow Keys / HJKL]\n\n"
		
		// Hitung seberapa banyak warna yang bisa dirender agar pas memenuhi tinggi layar
		maxColorsToRender := availableContentHeight * m.calculatedCols
		if maxColorsToRender > len(currPalette.Colors) {
			maxColorsToRender = len(currPalette.Colors)
		}
		
		for i := 0; i < maxColorsToRender; i++ {
			col := currPalette.Colors[i]
			blockColor := lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", col.R, col.G, col.B))
			
			if i == m.colorIndex {
				content += lipgloss.NewStyle().Background(lipgloss.Color("#333333")).Foreground(blockColor).Render(" █ ")
			} else {
				content += lipgloss.NewStyle().Foreground(blockColor).Render(" █ ")
			}
			
			if (i+1)%m.calculatedCols == 0 { content += "\n" }
		}
		if maxColorsToRender%m.calculatedCols != 0 { content += "\n" }
		
	case 1:
		content = "Manual RGB (Up/Down pilih baris, Left/Right geser nilai, Shift+Left/Right [+/- 10]):\n\n"
		labels := []string{"Red  ", "Green", "Blue "}
		vals := []int{m.manualR, m.manualG, m.manualB}
		// Skalakan panjang slider bar agar responsif mengikuti lebar terminal
		barMaxLength := m.termWidth - 25
		if barMaxLength < 10 { barMaxLength = 10 }
		if barMaxLength > 50 { barMaxLength = 50 }

		for i := 0; i < 3; i++ {
			cursor := " "
			if m.rgbCursor == i { cursor = ">" }
			content += fmt.Sprintf("%s %s: [%3d] ", cursor, labels[i], vals[i])
			
			bars := int((float64(vals[i]) / 255.0) * float64(barMaxLength))
			for b := 0; b < barMaxLength; b++ {
				if b < bars { content += "■" } else { content += " " }
			}
			content += "\n"
		}
	case 2:
		content = "Manual CMYK (Up/Down pilih baris, Left/Right geser nilai, Shift+Left/Right [+/- 10]):\n\n"
		labels := []string{"Cyan   ", "Magenta", "Yellow ", "Black  "}
		vals := []int{m.manualC, m.manualM, m.manualY, m.manualK}
		barMaxLength := m.termWidth - 25
		if barMaxLength < 10 { barMaxLength = 10 }
		if barMaxLength > 50 { barMaxLength = 50 }

		for i := 0; i < 4; i++ {
			cursor := " "
			if m.cmykCursor == i { cursor = ">" }
			content += fmt.Sprintf("%s %s: [%3d%%] ", cursor, labels[i], vals[i])
			
			bars := int((float64(vals[i]) / 100.0) * float64(barMaxLength))
			for b := 0; b < barMaxLength; b++ {
				if b < bars { content += "■" } else { content += " " }
			}
			content += "\n"
		}
	}

	// 2. PANEL PREVIEW (Dinamis dan selalu nempel di bawah isi content)
	previewStyle := lipgloss.NewStyle().Background(lipgloss.Color(hexStr)).Padding(1, 6).MarginRight(4)
	if (currentR*299 + currentG*587 + currentB*114) / 1000 > 128 {
		previewStyle = previewStyle.Foreground(lipgloss.Color("#000")) 
	} else {
		previewStyle = previewStyle.Foreground(lipgloss.Color("#FFF")) 
	}
	previewBox := previewStyle.Render("COLOR\nPREVIEW")
	
	infoText := fmt.Sprintf(
		"NAMA : %s\nHEX  : %s  |  RGB  : rgb(%d, %d, %d)  |  CMYK : cmyk(%d%%, %d%%, %d%%, %d%%)", 
		currentName, hexStr, currentR, currentG, currentB, currentC, currentM, currentY, currentK,
	)
	infoBox := lipgloss.NewStyle().PaddingTop(1).Render(infoText)
	previewRow := "\n" + lipgloss.JoinHorizontal(lipgloss.Top, previewBox, infoBox) + "\n"

	// 3. FOOTER HELP STRIP
	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333")).Render("[Tab] Mode | [P] Ganti Palet | [Q] Keluar")

	// 4. ABSOLUTE PADDING SYSTEM (Sistem Pengunci Layar Penuh)
	// Kita gabungkan teks di atas, lalu hitung total baris yang terpakai saat ini.
	currentOutput := header + content + previewRow + footer
	actualLines := strings.Count(currentOutput, "\n")
	
	// Jika output asli lebih pendek dari tinggi terminal, kita suntikkan baris kosong (\n) 
	// agar footer terdorong paksa ke baris paling bawah terminal secara presisi.
	padding := ""
	if actualLines < m.termHeight {
		padding = strings.Repeat("\n", m.termHeight-actualLines-1)
	}

	return header + content + padding + previewRow + footer
}

func main() {
	// Menjalankan aplikasi dengan mengaktifkan mode AltScreen (Mengambil alih seluruh terminal)
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Fatal Error: %v", err)
		os.Exit(1)
	}
}
