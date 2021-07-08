package mangasail

import (
	"context"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/bigscreen/manga-scrapper/domain"
)

type HomePageScrapper interface {
	GetContent() (domain.MangasailHomeMangas, error)
}

type homeScrapper struct {
	chromeCtx context.Context
}

func NewHomePageScrapper(chromeCtx context.Context) HomePageScrapper {
	return homeScrapper{chromeCtx: chromeCtx}
}

func (h homeScrapper) GetContent() (domain.MangasailHomeMangas, error) {
	waitSelector := `document.querySelector("#block-showmanga-hot-today")`
	wantedSelector := `document.querySelector("body > section > div > div > div.main-table > div")`
	document, err := getCrawledHtmlDocument(h.chromeCtx, HomeURL, waitSelector, wantedSelector)
	if err != nil {
		fmt.Println("GetHomeContent, failed to get html document, err:", err)
		return domain.MangasailHomeMangas{}, err
	}

	return h.buildHomeContent(document), nil
}

func (h homeScrapper) buildHomeContent(document *goquery.Document) domain.MangasailHomeMangas {
	if document == nil {
		return domain.MangasailHomeMangas{}
	}

	dailyHotMangasChannel := make(chan domain.MangasailMangas)
	latestMangasChannel := make(chan domain.MangasailMangas)
	hotMangasChannel := make(chan domain.MangasailMangas)
	newMangasChannel := make(chan domain.MangasailMangas)
	defer func() {
		close(dailyHotMangasChannel)
		close(latestMangasChannel)
		close(hotMangasChannel)
		close(newMangasChannel)
	}()

	execute := func(c chan domain.MangasailMangas, f func(goquery.Document) domain.MangasailMangas) {
		c <- f(*document)
	}
	go execute(dailyHotMangasChannel, h.getDailyHotMangas)
	go execute(latestMangasChannel, h.getLatestMangas)
	go execute(hotMangasChannel, h.getHotMangas)
	go execute(newMangasChannel, h.getNewMangas)

	return domain.MangasailHomeMangas{
		DailyHotMangas: <-dailyHotMangasChannel,
		LatestMangas:   <-latestMangasChannel,
		HotMangas:      <-hotMangasChannel,
		NewMangas:      <-newMangasChannel,
	}
}

func (h homeScrapper) getDailyHotMangas(document goquery.Document) domain.MangasailMangas {
	var mangas domain.MangasailMangas
	document.Find("#block-showmanga-hot-today .content #hottoday-list").
		Each(func(pos int, selection *goquery.Selection) {
			selection.Find("li").Each(func(liPos int, liSelection *goquery.Selection) {
				tSelection := liSelection.Find("a.mtitle")
				chapter := domain.MangasailChapter{}
				chapter.Title = tSelection.Text()
				chapter.Path, _ = tSelection.Attr("href")
				manga := domain.MangasailManga{}
				manga.Path, _ = liSelection.Find("a").Attr("href")
				manga.IconURL, _ = liSelection.Find("a img.img-responsive").Attr("src")
				manga.Chapters = domain.MangasailChapters{chapter}
				mangas = append(mangas, manga)
			})
		})
	return mangas
}

func (h homeScrapper) getLatestMangas(document goquery.Document) domain.MangasailMangas {
	var mangas domain.MangasailMangas
	document.Find("#block-showmanga-lastest-list #latest-list").
		Each(func(pos int, selection *goquery.Selection) {
			selection.Find("li").
				Each(func(liPos int, liSelection *goquery.Selection) {
					manga := domain.MangasailManga{}
					manga.IconURL, _ = liSelection.Find("a img.img-responsive").Attr("src")
					liSelection.Find("ul li").
						Each(func(liChildPos int, liChildSelection *goquery.Selection) {
							liChildSelection.Find(".tl").
								Each(func(tlLiChildPos int, tlLiChildSelection *goquery.Selection) {
									aSelection := tlLiChildSelection.Find("a")
									manga.Name = aSelection.Find("strong").Text()
									manga.Path, _ = aSelection.Attr("href")
								})
							liChildSelection.Find("#c-list li").
								Each(func(cListLiChildPos int, cListLiChildSelection *goquery.Selection) {
									if cListLiChildPos == 0 {
										aSelection := cListLiChildSelection.Find("a")
										chapter := domain.MangasailChapter{}
										chapter.Title = aSelection.Text()
										chapter.Path, _ = aSelection.Attr("href")
										chapter.LastModified = cListLiChildSelection.Find("span.tm").Text()
										manga.Chapters = domain.MangasailChapters{chapter}
									}
								})
						})
					mangas = append(mangas, manga)
				})

		})
	return mangas
}

func (h homeScrapper) getHotMangas(document goquery.Document) domain.MangasailMangas {
	var mangas domain.MangasailMangas
	document.Find("#block-showmanga-hot-manga #new-list").
		Each(func(pos int, selection *goquery.Selection) {
			selection.Find("li").Each(func(liPos int, liSelection *goquery.Selection) {
				manga := domain.MangasailManga{}
				manga.IconURL, _ = liSelection.Find("a img.img-responsive").Attr("src")
				liSelection.Find("ul li").
					Each(func(liChildPos int, liChildSelection *goquery.Selection) {
						liChildSelection.Find(".tl").
							Each(func(tlLiChildPos int, tlLiChildSelection *goquery.Selection) {
								aSelection := tlLiChildSelection.Find("a")
								manga.Name = aSelection.Text()
								manga.Path, _ = aSelection.Attr("href")
							})
						liChildSelection.Find(".cl").
							Each(func(clLiChildPos int, clLiChildSelection *goquery.Selection) {
								aSelection := clLiChildSelection.Find("a")
								chapter := domain.MangasailChapter{}
								chapter.Title = aSelection.Text()
								chapter.Path, _ = aSelection.Attr("href")
								manga.Chapters = domain.MangasailChapters{chapter}
							})
					})
				mangas = append(mangas, manga)
			})
		})
	return mangas
}

func (h homeScrapper) getNewMangas(document goquery.Document) domain.MangasailMangas {
	var mangas domain.MangasailMangas
	document.Find("#block-showmanga-new-manga #new-list").
		Each(func(pos int, selection *goquery.Selection) {
			selection.Find("li").Each(func(liPos int, liSelection *goquery.Selection) {
				manga := domain.MangasailManga{}
				manga.IconURL, _ = liSelection.Find("a img.img-responsive").Attr("src")
				liSelection.Find("ul li").
					Each(func(liChildPos int, liChildSelection *goquery.Selection) {
						liChildSelection.Find(".tl").
							Each(func(tlLiChildPos int, tlLiChildSelection *goquery.Selection) {
								aSelection := tlLiChildSelection.Find("a")
								manga.Name = aSelection.Text()
								manga.Path, _ = aSelection.Attr("href")
							})
						liChildSelection.Find(".cl").
							Each(func(clLiChildPos int, clLiChildSelection *goquery.Selection) {
								aSelection := clLiChildSelection.Find("a")
								chapter := domain.MangasailChapter{}
								chapter.Title = aSelection.Text()
								chapter.Path, _ = aSelection.Attr("href")
								manga.Chapters = domain.MangasailChapters{chapter}
							})
					})
				mangas = append(mangas, manga)
			})
		})
	return mangas
}
