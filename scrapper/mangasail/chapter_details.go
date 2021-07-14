package mangasail

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/bigscreen/manga-scrapper/domain"
	"github.com/bigscreen/manga-scrapper/errors"
)

type ChapterDetailsPageScrapper interface {
	GetContent(path string) (domain.Chapter, error)
}

type chapterDetailsScrapper struct{}

func NewChapterDetailsPageScrapper() ChapterDetailsPageScrapper {
	return chapterDetailsScrapper{}
}

func (c chapterDetailsScrapper) GetContent(path string) (domain.Chapter, error) {
	const op = "ChapterDetailsPageScrapper.GetContent"
	waitSelector := `document.querySelector("#node-337513")`
	wantedSelector := `document.querySelector("body > section > div > div > div > div > section > div")`
	document, err := getCrawledHtmlDocument(buildPageURL(path), waitSelector, wantedSelector)
	if err != nil {
		return domain.Chapter{}, errors.New(errors.WithOp(op), errors.WithError(err))
	}

	return c.buildReaderContent(document), nil
}

func (c chapterDetailsScrapper) buildReaderContent(document *goquery.Document) domain.Chapter {
	if document == nil {
		return domain.Chapter{}
	}

	titleChannel := make(chan string)
	imagesChannel := make(chan domain.Images)
	defer func() {
		close(titleChannel)
		close(imagesChannel)
	}()

	go func() {
		titleChannel <- c.getTitle(*document)
	}()
	go func() {
		imagesChannel <- c.getImages(*document)
	}()

	return domain.Chapter{
		Title:  <-titleChannel,
		Images: <-imagesChannel,
	}
}

func (c chapterDetailsScrapper) getTitle(document goquery.Document) string {
	return document.Find("h1.page-header a.active").Text()
}

func (c chapterDetailsScrapper) getImages(document goquery.Document) domain.Images {
	var images domain.Images
	document.Find("#node-337513 #images").Each(func(pos int, selection *goquery.Selection) {
		selection.Find("img").Each(func(iPos int, iSelection *goquery.Selection) {
			image := domain.Image{}
			image.ID, _ = iSelection.Attr("name")
			image.ImageURL, _ = iSelection.Attr("src")
			images = append(images, image)
		})
	})
	return images
}
