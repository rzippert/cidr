package main

import (
	"fmt"
	"net"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type cidrInfo struct {
	Netmask        string
	Wildcard       string
	TotalAddresses string
	MaxAddresses   string
	Network        string
	CIDRNotation   string
	AddrRange      string
	Error          string
}

type model struct {
	ipInput   textinput.Model
	maskInput textinput.Model
	cursor    int
	results   cidrInfo
}

func calculateCIDR(ipStr, maskStr string) cidrInfo {
	if ipStr == "" || maskStr == "" {
		return cidrInfo{}
	}

	ip, ipnet, err := net.ParseCIDR(ipStr + "/" + maskStr)
	if err != nil {
		if net.ParseIP(ipStr) == nil {
			return cidrInfo{Error: "Invalid IP Address"}
		}
		return cidrInfo{Error: "Invalid Mask"}
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return cidrInfo{Error: "Only IPv4 addresses are supported."}
	}

	mask := ipnet.Mask
	ones, _ := mask.Size()

	// The number of addresses in the subnet
	totalAddresses := uint64(1) << (32 - ones) // Terminology from website

	var maxUsableAddresses uint64
	if ones == 32 {
		maxUsableAddresses = 1
	} else if ones == 31 {
		maxUsableAddresses = 2
	} else {
		maxUsableAddresses = totalAddresses - 2
	}

	// The netmask is the last 4 bytes of the 16-byte mask
	netmask := mask[len(mask)-4:]

	// The wildcard is the bitwise NOT of the netmask
	wildcard := make(net.IP, 4)
	for i := 0; i < 4; i++ {
		wildcard[i] = ^netmask[i]
	}

	// The broadcast address is the network address OR'd with the wildcard mask
	broadcast := make(net.IP, 4)
	networkIP := ipnet.IP.To4()
	for i := 0; i < 4; i++ {
		broadcast[i] = networkIP[i] | wildcard[i]
	}

	return cidrInfo{
		Netmask:        fmt.Sprintf("%d.%d.%d.%d", netmask[0], netmask[1], netmask[2], netmask[3]),
		Wildcard:       wildcard.String(),
		TotalAddresses: fmt.Sprintf("%d", totalAddresses),
		MaxAddresses:   fmt.Sprintf("%d", maxUsableAddresses),
		Network:        networkIP.String(),
		CIDRNotation:   ipnet.String(),
		AddrRange:      fmt.Sprintf("%s - %s", networkIP.String(), broadcast.String()),
	}
}

func initialModel() model {
	ip := textinput.New()
	ip.Placeholder = "IP Address"
	ip.Focus()
	ip.CharLimit = 15
	ip.Width = 20

	mask := textinput.New()
	mask.Placeholder = "Mask Bits"
	mask.CharLimit = 2
	mask.Width = 20

	return model{
		ipInput:   ip,
		maskInput: mask,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "shift+tab":
			if m.cursor == 0 {
				m.cursor = 1
				m.ipInput.Blur()
				m.maskInput.Focus()
			} else {
				m.cursor = 0
				m.maskInput.Blur()
				m.ipInput.Focus()
			}
		}
	}

	if m.cursor == 0 {
		m.ipInput, cmd = m.ipInput.Update(msg)
	} else {
		m.maskInput, cmd = m.maskInput.Update(msg)
	}

	m.results = calculateCIDR(m.ipInput.Value(), m.maskInput.Value())

	return m, cmd
}

func (m model) View() string {
	var results string
	if m.results.Error != "" {
		results = m.results.Error
	} else {
		results = fmt.Sprintf(
			"CIDR Netmask: %s\nWildcard Mask: %s\nTotal Addresses: %s\nMaximum Addresses: %s\nCIDR Network (Route): %s\nNet: CIDR Notation: %s\nCIDR Address Range: %s",
			m.results.Netmask,
			m.results.Wildcard,
			m.results.TotalAddresses,
			m.results.MaxAddresses,
			m.results.Network,
			m.results.CIDRNotation,
			m.results.AddrRange,
		)
	}

	return fmt.Sprintf(
		"CIDR Calculator\n\n%s\n%s\n\n%s\n\n%s",
		m.ipInput.View(),
		m.maskInput.View(),
		results,
		"(q to quit)",
	)
}

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("cidr version %s, commit %s, built at %s\n", version, commit, date)
		os.Exit(0)
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
