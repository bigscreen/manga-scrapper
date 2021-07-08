package domain

import (
	"fmt"
	"strings"
)

type (
	MangasailMangas   []MangasailManga
	MangasailChapters []MangasailChapter
)

type MangasailManga struct {
	Name     string
	Path     string
	IconURL  string
	Chapters MangasailChapters
}

type MangasailChapter struct {
	Title        string
	Path         string
	LastModified string
}

type MangasailHomeMangas struct {
	DailyHotMangas MangasailMangas
	LatestMangas   MangasailMangas
	HotMangas      MangasailMangas
	NewMangas      MangasailMangas
}

func (m MangasailManga) ToString() string {
	var fields []string
	fields = append(fields, fmt.Sprintf(`Name: "%s"`, m.Name))
	fields = append(fields, fmt.Sprintf(`Path: "%s"`, m.Path))
	fields = append(fields, fmt.Sprintf(`IconURL: "%s"`, m.IconURL))
	fields = append(fields, fmt.Sprintf(`Chapters: [%s]`, m.Chapters.ToString()))
	return fmt.Sprintf("{%s}", strings.Join(fields, ", "))
}

func (m MangasailChapter) ToString() string {
	var fields []string
	fields = append(fields, fmt.Sprintf(`Title: "%s"`, m.Title))
	fields = append(fields, fmt.Sprintf(`Path: "%s"`, m.Path))
	fields = append(fields, fmt.Sprintf(`LastModified: "%s"`, m.LastModified))
	return fmt.Sprintf("{%s}", strings.Join(fields, ", "))
}

func (m MangasailMangas) ToString() string {
	var mt []string
	for _, manga := range m {
		mt = append(mt, manga.ToString())
	}
	return strings.Join(mt, ", ")
}

func (m MangasailChapters) ToString() string {
	var ct []string
	for _, chapter := range m {
		ct = append(ct, chapter.ToString())
	}
	return strings.Join(ct, ", ")
}

func (m MangasailHomeMangas) ToString() string {
	var fields []string
	fields = append(fields, fmt.Sprintf("DailyHotMangas: [\n%s\n]", m.DailyHotMangas.ToString()))
	fields = append(fields, fmt.Sprintf("LatestMangas: [\n%s\n]", m.LatestMangas.ToString()))
	fields = append(fields, fmt.Sprintf("HotMangas: [\n%s\n]", m.HotMangas.ToString()))
	fields = append(fields, fmt.Sprintf("NewMangas: [\n%s\n]", m.NewMangas.ToString()))
	return fmt.Sprintf("{\n%s\n}", strings.Join(fields, ",\n"))
}
