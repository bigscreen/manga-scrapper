package mangasail

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bigscreen/manga-scrapper/domain"
)

type DetailsPageScrapper interface {
	GetContent(path string) (domain.MangasailManga, error)
}

type detailsScrapper struct {
	chromeCtx context.Context
}

func NewDetailsPageScrapper(chromeCtx context.Context) DetailsPageScrapper {
	return detailsScrapper{chromeCtx: chromeCtx}
}

func (d detailsScrapper) GetContent(path string) (domain.MangasailManga, error) {
	waitSelector := `document.querySelector("#node-254222")`
	wantedSelector := `document.querySelector("body > section > div > div > div.main-table > div > section > div")`
	document, err := getCrawledHtmlDocument(d.chromeCtx, d.buildDetailsURL(path), waitSelector, wantedSelector)
	if err != nil {
		fmt.Println("GetDetailsContent, failed to get html document, err:", err)
		return domain.MangasailManga{}, err
	}

	return d.buildDetailsContent(document), nil
}

func (d detailsScrapper) buildDetailsURL(path string) string {
	return HomeURL + path
}

func (d detailsScrapper) buildDetailsContent(document *goquery.Document) domain.MangasailManga {
	if document == nil {
		return domain.MangasailManga{}
	}

	attributesChannel := make(chan domain.MangasailManga)
	chaptersChannel := make(chan domain.MangasailChapters)
	defer func() {
		close(attributesChannel)
		close(chaptersChannel)
	}()

	go func() {
		attributesChannel <- d.getAttributes(*document)
	}()
	go func() {
		chaptersChannel <- d.getChapters(*document)
	}()

	manga := <-attributesChannel
	manga.Chapters = <-chaptersChannel
	return manga
}

func (d detailsScrapper) getAttributes(document goquery.Document) domain.MangasailManga {
	manga := domain.MangasailManga{}
	document.Find(".main-content-inner").Each(func(pos int, selection *goquery.Selection) {
		manga.Name = selection.Find("h1.page-header").Text()
		selection.Find("#node-254222 .content").Each(func(nPos int, nSelection *goquery.Selection) {
			manga.IconURL, _ = nSelection.Find(".field-name-field-image2 .field-items .field-item img").Attr("src")
			manga.ReleaseYear = nSelection.Find(".field-name-field-year-of-release .field-items .field-item").Text()
			manga.Status = nSelection.Find(".field-name-field-status .field-items .field-item").Text()
			manga.Author = nSelection.Find(".field-name-field-author .field-items .field-item").Text()
			manga.Description = d.formatMangaDescription(nSelection.Find(".field-name-body .field-items .field-item p").Text())
			var genres []string
			nSelection.Find(".field-name-field-genres .field-items .field-item").
				Each(func(gPos int, gSelection *goquery.Selection) {
					genres = append(genres, gSelection.Find("a").Text())
				})
			manga.Genres = genres
		})
	})
	return manga
}

func (d detailsScrapper) formatMangaDescription(s string) string {
	ns := strings.ReplaceAll(s, `"`, "")
	ns = strings.ReplaceAll(s, "<br>", "")
	ns = strings.ReplaceAll(s, "\n", " ")
	return ns
}

func (d detailsScrapper) getChapters(document goquery.Document) domain.MangasailChapters {
	var chapters domain.MangasailChapters
	document.Find("#node-254222 table.chlist tbody tr").Each(func(pos int, selection *goquery.Selection) {
		chapter := domain.MangasailChapter{}
		selection.Find("td").Each(func(tPos int, tSelection *goquery.Selection) {
			if tPos == 0 {
				aSelection := tSelection.Find("a")
				chapter.Title = aSelection.Text()
				chapter.Path, _ = aSelection.Attr("href")
			} else {
				chapter.LastModified = tSelection.Text()
			}
		})
		chapters = append(chapters, chapter)
	})
	return chapters
}
