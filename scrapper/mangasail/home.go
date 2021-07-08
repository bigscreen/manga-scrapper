package mangasail

import (
	"context"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/bigscreen/manga-scrapper/common"
	"github.com/bigscreen/manga-scrapper/domain"
)

type HomePageScrapper interface {
	GetContent() (domain.HomeMangas, error)
}

type homeScrapper struct {
	chromeCtx context.Context
}

func NewHomePageScrapper(chromeCtx context.Context) HomePageScrapper {
	return homeScrapper{chromeCtx: chromeCtx}
}

func (h homeScrapper) GetContent() (domain.HomeMangas, error) {
	waitSelector := `document.querySelector("#block-showmanga-hot-today")`
	wantedSelector := `document.querySelector("body > section > div > div > div.main-table > div")`
	document, err := getCrawledHtmlDocument(h.chromeCtx, common.MangasailBaseURL, waitSelector, wantedSelector)
	if err != nil {
		fmt.Println("GetHomeContent, failed to get html document, err:", err)
		return domain.HomeMangas{}, err
	}

	return h.buildHomeContent(document), nil
}

func (h homeScrapper) buildHomeContent(document *goquery.Document) domain.HomeMangas {
	if document == nil {
		return domain.HomeMangas{}
	}

	dailyHotMangasChannel := make(chan domain.Mangas)
	latestMangasChannel := make(chan domain.Mangas)
	popularMangasChannel := make(chan domain.Mangas)
	newMangasChannel := make(chan domain.Mangas)
	defer func() {
		close(dailyHotMangasChannel)
		close(latestMangasChannel)
		close(popularMangasChannel)
		close(newMangasChannel)
	}()

	execute := func(c chan domain.Mangas, f func(goquery.Document) domain.Mangas) {
		c <- f(*document)
	}
	go execute(dailyHotMangasChannel, h.getDailyHotMangas)
	go execute(latestMangasChannel, h.getLatestMangas)
	go execute(popularMangasChannel, h.getPopularMangas)
	go execute(newMangasChannel, h.getNewMangas)

	return domain.HomeMangas{
		DailyHotMangas: <-dailyHotMangasChannel,
		LatestMangas:   <-latestMangasChannel,
		PopularMangas:  <-popularMangasChannel,
		NewMangas:      <-newMangasChannel,
	}
}

func (h homeScrapper) getDailyHotMangas(document goquery.Document) domain.Mangas {
	var mangas domain.Mangas
	document.Find("#block-showmanga-hot-today .content #hottoday-list").
		Each(func(pos int, selection *goquery.Selection) {
			selection.Find("li").Each(func(liPos int, liSelection *goquery.Selection) {
				tSelection := liSelection.Find("a.mtitle")
				chapter := domain.Chapter{}
				chapter.Title = tSelection.Text()
				chapterPath, _ := tSelection.Attr("href")
				chapter.ID = getIdFromPath(chapterPath)
				manga := domain.Manga{}
				mangaPath, _ := liSelection.Find("a").Attr("href")
				manga.ID = getIdFromPath(mangaPath)
				manga.IconURL, _ = liSelection.Find("a img.img-responsive").Attr("src")
				manga.Chapters = domain.Chapters{chapter}
				mangas = append(mangas, manga)
			})
		})
	return mangas
}

func (h homeScrapper) getLatestMangas(document goquery.Document) domain.Mangas {
	var mangas domain.Mangas
	document.Find("#block-showmanga-lastest-list #latest-list").Each(func(pos int, selection *goquery.Selection) {
		selection.Find("li").Each(func(liPos int, liSelection *goquery.Selection) {
			manga := domain.Manga{}
			manga.IconURL, _ = liSelection.Find("a img.img-responsive").Attr("src")
			liSelection.Find("ul li").Each(func(liChildPos int, liChildSelection *goquery.Selection) {
				liChildSelection.Find(".tl").Each(func(tlLiChildPos int, tlLiChildSelection *goquery.Selection) {
					aSelection := tlLiChildSelection.Find("a")
					mangaPath, _ := aSelection.Attr("href")
					manga.Name = aSelection.Find("strong").Text()
					manga.ID = getIdFromPath(mangaPath)
				})
				liChildSelection.Find("#c-list li").Each(func(cListLiChildPos int, cListLiChildSelection *goquery.Selection) {
					if cListLiChildPos == 0 {
						aSelection := cListLiChildSelection.Find("a")
						chapter := domain.Chapter{}
						chapter.Title = aSelection.Text()
						chapterPath, _ := aSelection.Attr("href")
						chapter.ID = getIdFromPath(chapterPath)
						chapter.LastModified = cListLiChildSelection.Find("span.tm").Text()
						manga.Chapters = domain.Chapters{chapter}
					}
				})
			})
			mangas = append(mangas, manga)
		})

	})
	return mangas
}

func (h homeScrapper) getPopularMangas(document goquery.Document) domain.Mangas {
	var mangas domain.Mangas
	document.Find("#block-showmanga-hot-manga #new-list").Each(func(pos int, selection *goquery.Selection) {
		selection.Find("li").Each(func(liPos int, liSelection *goquery.Selection) {
			manga := domain.Manga{}
			manga.IconURL, _ = liSelection.Find("a img.img-responsive").Attr("src")
			liSelection.Find("ul li").Each(func(liChildPos int, liChildSelection *goquery.Selection) {
				liChildSelection.Find(".tl").Each(func(tlLiChildPos int, tlLiChildSelection *goquery.Selection) {
					aSelection := tlLiChildSelection.Find("a")
					mangaPath, _ := aSelection.Attr("href")
					manga.Name = aSelection.Text()
					manga.ID = getIdFromPath(mangaPath)
				})
				liChildSelection.Find(".cl").Each(func(clLiChildPos int, clLiChildSelection *goquery.Selection) {
					aSelection := clLiChildSelection.Find("a")
					chapter := domain.Chapter{}
					chapter.Title = aSelection.Text()
					chapterPath, _ := aSelection.Attr("href")
					chapter.ID = getIdFromPath(chapterPath)
					manga.Chapters = domain.Chapters{chapter}
				})
			})
			mangas = append(mangas, manga)
		})
	})
	return mangas
}

func (h homeScrapper) getNewMangas(document goquery.Document) domain.Mangas {
	var mangas domain.Mangas
	document.Find("#block-showmanga-new-manga #new-list").Each(func(pos int, selection *goquery.Selection) {
		selection.Find("li").Each(func(liPos int, liSelection *goquery.Selection) {
			manga := domain.Manga{}
			manga.IconURL, _ = liSelection.Find("a img.img-responsive").Attr("src")
			liSelection.Find("ul li").Each(func(liChildPos int, liChildSelection *goquery.Selection) {
				liChildSelection.Find(".tl").Each(func(tlLiChildPos int, tlLiChildSelection *goquery.Selection) {
					aSelection := tlLiChildSelection.Find("a")
					mangaPath, _ := aSelection.Attr("href")
					manga.Name = aSelection.Text()
					manga.ID = getIdFromPath(mangaPath)
				})
				liChildSelection.Find(".cl").Each(func(clLiChildPos int, clLiChildSelection *goquery.Selection) {
					aSelection := clLiChildSelection.Find("a")
					chapter := domain.Chapter{}
					chapter.Title = aSelection.Text()
					chapterPath, _ := aSelection.Attr("href")
					chapter.ID = getIdFromPath(chapterPath)
					manga.Chapters = domain.Chapters{chapter}
				})
			})
			mangas = append(mangas, manga)
		})
	})
	return mangas
}
