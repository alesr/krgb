package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sstallion/go-hid"

	"github.com/alesr/krgb/config"
	"github.com/alesr/krgb/via"
	"github.com/alesr/krgb/xerrors"
)

const (
	focusEffect = iota
	focusHue
	focusSat
	focusBrightness
	focusSpeed
	focusCount
)

type effectEntry struct {
	name string
	mode byte
}

var effects = []effectEntry{
	{name: "Off", mode: 0},
	{name: "Solid", mode: 1},
	{name: "Alphas Mods", mode: 2},
	{name: "Gradient Up Down", mode: 3},
	{name: "Gradient Left Right", mode: 4},
	{name: "Breathing", mode: 5},
	{name: "Band (Sat)", mode: 6},
	{name: "Band (Val)", mode: 7},
	{name: "Cycle All", mode: 12},
	{name: "Cycle Left Right", mode: 13},
	{name: "Cycle Up Down", mode: 14},
	{name: "Dual Beacon", mode: 19},
	{name: "Rainbow Beacon", mode: 20},
	{name: "Rainbow Pinwheels", mode: 21},
	{name: "Raindrops", mode: 22},
	{name: "Jellybean", mode: 23},
	{name: "Hue Breathing", mode: 25},
	{name: "Hue Wave", mode: 27},
	{name: "Pixel Rain", mode: 30},
	{name: "Starlight", mode: 31},
	{name: "Solid Reactive", mode: 37},
	{name: "Reactive Cross", mode: 40},
	{name: "Splash", mode: 44},
}

var (
	focusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	titleStyle = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
)

type model struct {
	device     *hid.Device
	deviceInfo *hid.DeviceInfo

	focus      int
	effectIdx  int
	hue        int
	sat        int
	brightness int
	speed      int

	statusMsg string
}

type clearStatusMsg struct{}

func New(dev *hid.Device, info *hid.DeviceInfo, s *config.Settings) tea.Model {
	idx := max(s.EffectIdx, 0)
	if idx >= len(effects) {
		idx = len(effects) - 1
	}

	return &model{
		device:     dev,
		deviceInfo: info,
		effectIdx:  idx,
		hue:        clamp(s.Hue),
		sat:        clamp(s.Sat),
		brightness: clamp(s.Brightness),
		speed:      clamp(s.Speed),
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			m.focus--
			if m.focus < 0 {
				m.focus = focusCount - 1
			}

		case "down", "j":
			m.focus++
			if m.focus >= focusCount {
				m.focus = 0
			}

		case "left", "h":
			m.changeFocused(-1)

		case "right", "l":
			m.changeFocused(1)

		case "w":
			s := &config.Settings{
				EffectIdx:  m.effectIdx,
				Hue:        m.hue,
				Sat:        m.sat,
				Brightness: m.brightness,
				Speed:      m.speed,
			}

			if err := s.Save(); err != nil {
				m.statusMsg = fmt.Sprintf("Error [%d]: %v", xerrors.CodeOf(err), err)
				return m, nil
			}

			if m.device != nil {
				var errs []string
				if err := via.SaveValue(m.device, via.ValEffect, []byte{effects[m.effectIdx].mode}); err != nil {
					errs = append(errs, err.Error())
				}

				if err := via.SaveValue(m.device, via.ValColor, []byte{byte(m.hue), byte(m.sat)}); err != nil {
					errs = append(errs, err.Error())
				}

				if err := via.SaveValue(m.device, via.ValBrightness, []byte{byte(m.brightness)}); err != nil {
					errs = append(errs, err.Error())
				}

				if err := via.SaveValue(m.device, via.ValSpeed, []byte{byte(m.speed)}); err != nil {
					errs = append(errs, err.Error())
				}

				if len(errs) > 0 {
					m.statusMsg = fmt.Sprintf("Config saved, but keyboard save failed: %s", strings.Join(errs, "; "))
					return m, func() tea.Msg {
						time.Sleep(3000 * time.Millisecond)
						return clearStatusMsg{}
					}
				}
			}

			m.statusMsg = "Config saved"

			return m, func() tea.Msg {
				time.Sleep(1500 * time.Millisecond)
				return clearStatusMsg{}
			}
		}

	case clearStatusMsg:
		m.statusMsg = ""
		return m, nil
	}
	return m, nil
}

func (m *model) changeFocused(dir int) {
	switch m.focus {
	case focusEffect:
		m.effectIdx += dir

		if m.effectIdx < 0 {
			m.effectIdx = len(effects) - 1
		} else if m.effectIdx >= len(effects) {
			m.effectIdx = 0
		}

		if m.device != nil {
			via.SetValue(m.device, via.ValEffect, []byte{effects[m.effectIdx].mode})
		}

	case focusHue:
		m.hue = clamp(m.hue + dir*4)

		if m.device != nil {
			via.SetValue(m.device, via.ValColor, []byte{byte(m.hue), byte(m.sat)})
		}

	case focusSat:
		m.sat = clamp(m.sat + dir*4)

		if m.device != nil {
			via.SetValue(m.device, via.ValColor, []byte{byte(m.hue), byte(m.sat)})
		}

	case focusBrightness:
		m.brightness = clamp(m.brightness + dir*4)

		if m.device != nil {
			via.SetValue(m.device, via.ValBrightness, []byte{byte(m.brightness)})
		}

	case focusSpeed:
		m.speed = clamp(m.speed + dir*4)

		if m.device != nil {
			via.SetValue(m.device, via.ValSpeed, []byte{byte(m.speed)})
		}
	}
}

func clamp(v int) int {
	if v < 0 {
		return 0
	}

	if v > 255 {
		return 255
	}
	return v
}

func (m model) View() string {
	var s string
	s = titleStyle.Render("krgb — Keychron LED Control") + "\n\n"

	s += fmt.Sprintf("  Device: %s\n", m.deviceInfo.ProductStr)
	s += "  Status: \u25cf Connected\n\n"

	s += m.renderRow(0, "Effect", effects[m.effectIdx].name)
	s += m.renderSliderRow(1, "Hue", m.hue)
	s += m.renderSliderRow(2, "Sat", m.sat)
	s += m.renderSliderRow(3, "Bright", m.brightness)
	s += m.renderSliderRow(4, "Speed", m.speed)

	s += "\n"

	if m.statusMsg != "" {
		s += fmt.Sprintf("  %s\n\n", okStyle.Render(m.statusMsg))
	}

	s += "  w: save config  \u00b7  q: quit\n"
	s += "  \u2191\u2193: nav  \u00b7  \u2190\u2192: change"
	return s
}

func (m model) renderLabel(index int) string {
	if m.focus == index {
		return focusStyle.Render(">")
	}
	return " "
}

func (m model) renderRow(index int, label, value string) string {
	marker := m.renderLabel(index)
	return fmt.Sprintf("  %s %-6s %s\n", marker, label, value)
}

func (m model) renderSliderRow(index int, label string, value int) string {
	marker := m.renderLabel(index)

	const barLen = 12
	filled := min(max(value*barLen/255, 0), barLen)
	bar := strings.Repeat("\u2588", filled) + strings.Repeat("\u2591", barLen-filled)

	return fmt.Sprintf("  %s %-6s %s %3d\n", marker, label, bar, value)
}
