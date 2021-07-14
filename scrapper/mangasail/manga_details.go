package mangasail

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bigscreen/manga-scrapper/domain"
	"github.com/bigscreen/manga-scrapper/errors"
)

type MangaDetailsPageScrapper interface {
	GetContent(path string) (domain.Manga, error)
}

type mangaDetailsScrapper struct{}

func NewMangaDetailsPageScrapper() MangaDetailsPageScrapper {
	return mangaDetailsScrapper{}
}

func (m mangaDetailsScrapper) GetContent(path string) (domain.Manga, error) {
	const op = "MangaDetailsPageScrapper.GetContent"
	waitSelector := `document.querySelector("#node-254222")`
	wantedSelector := `document.querySelector("body > section > div > div > div.main-table > div > section > div")`
	document, err := getCrawledHtmlDocument(buildPageURL(path), waitSelector, wantedSelector)
	if err != nil {
		return domain.Manga{}, errors.New(errors.WithOp(op), errors.WithError(err))
	}

	return m.buildDetailsContent(document), nil
}

func (m mangaDetailsScrapper) buildDetailsContent(document *goquery.Document) domain.Manga {
	if document == nil {
		return domain.Manga{}
	}

	attributesChannel := make(chan domain.Manga)
	chaptersChannel := make(chan domain.Chapters)
	defer func() {
		close(attributesChannel)
		close(chaptersChannel)
	}()

	go func() {
		attributesChannel <- m.getAttributes(*document)
	}()
	go func() {
		chaptersChannel <- m.getChapters(*document)
	}()

	manga := <-attributesChannel
	manga.Chapters = <-chaptersChannel
	return manga
}

func (m mangaDetailsScrapper) getAttributes(document goquery.Document) domain.Manga {
	manga := domain.Manga{}
	document.Find(".main-content-inner").Each(func(pos int, selection *goquery.Selection) {
		manga.Name = selection.Find("h1.page-header").Text()
		selection.Find("#node-254222 .content").Each(func(nPos int, nSelection *goquery.Selection) {
			manga.IconURL, _ = nSelection.Find(".field-name-field-image2 .field-items .field-item img").Attr("src")
			manga.ReleaseYear = nSelection.Find(".field-name-field-year-of-release .field-items .field-item").Text()
			manga.Status = nSelection.Find(".field-name-field-status .field-items .field-item").Text()
			manga.Author = nSelection.Find(".field-name-field-author .field-items .field-item").Text()
			manga.Description = m.formatMangaDescription(nSelection.Find(".field-name-body .field-items .field-item p").Text())
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

func (m mangaDetailsScrapper) formatMangaDescription(s string) string {
	ns := strings.ReplaceAll(s, `"`, "")
	ns = strings.ReplaceAll(s, "<br>", "")
	ns = strings.ReplaceAll(s, "\n", " ")
	return ns
}

func (m mangaDetailsScrapper) getChapters(document goquery.Document) domain.Chapters {
	var chapters domain.Chapters
	document.Find("#node-254222 table.chlist tbody tr").Each(func(pos int, selection *goquery.Selection) {
		chapter := domain.Chapter{}
		selection.Find("td").Each(func(tPos int, tSelection *goquery.Selection) {
			if tPos == 0 {
				aSelection := tSelection.Find("a")
				chapter.Title = aSelection.Text()
				chapterPath, _ := aSelection.Attr("href")
				chapter.ID = getIdFromPath(chapterPath)
			} else {
				chapter.LastModified = tSelection.Text()
			}
		})
		chapters = append(chapters, chapter)
	})
	return chapters
}
