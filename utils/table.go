package utils

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func GenerateTable(headers []string, rows [][]string, columnWidths []int) string {
	// 定义表格样式
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	// 创建表格
	t := table.New(
		table.WithColumns(createColumns(headers, columnWidths)),
		table.WithRows(createRows(rows)),
		table.WithFocused(false),
		table.WithHeight(len(rows)+1),
	)

	// 设置表格样式
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.Foreground(lipgloss.Color("255"))
	t.SetStyles(s)

	// 直接渲染并打印表格
	return baseStyle.Render(t.View())
}

func createColumns(headers []string, widths []int) []table.Column {
	columns := make([]table.Column, len(headers))
	for i, header := range headers {
		columns[i] = table.Column{Title: header, Width: widths[i]}
	}
	return columns
}

func createRows(rows [][]string) []table.Row {
	tableRows := make([]table.Row, len(rows))
	for i, row := range rows {
		tableRows[i] = row
	}
	return tableRows
}
