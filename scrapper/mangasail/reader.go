package mangasail

import (
	"context"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/bigscreen/manga-scrapper/domain"
)

type ReaderPageScrapper interface {
	GetContent(path string) (domain.MangasailChapter, error)
}

type readerScrapper struct {
	chromeCtx context.Context
}

func NewReaderPageScrapper(chromeCtx context.Context) ReaderPageScrapper {
	return readerScrapper{chromeCtx: chromeCtx}
}

func (r readerScrapper) GetContent(path string) (domain.MangasailChapter, error) {
	waitSelector := `document.querySelector("#node-337513")`
	wantedSelector := `document.querySelector("body > section > div > div > div > div > section > div")`
	document, err := getCrawledHtmlDocument(r.chromeCtx, r.buildDetailsURL(path), waitSelector, wantedSelector)
	if err != nil {
		fmt.Println("GetReaderContent, failed to get html document, err:", err)
		return domain.MangasailChapter{}, err
	}

	return r.buildReaderContent(document), nil
}

func (r readerScrapper) buildDetailsURL(path string) string {
	return HomeURL + path
}

func (r readerScrapper) buildReaderContent(document *goquery.Document) domain.MangasailChapter {
	if document == nil {
		return domain.MangasailChapter{}
	}

	titleChannel := make(chan string)
	imagesChannel := make(chan domain.MangasailImages)
	defer func() {
		close(titleChannel)
		close(imagesChannel)
	}()

	go func() {
		titleChannel <- r.getTitle(*document)
	}()
	go func() {
		imagesChannel <- r.getImages(*document)
	}()

	return domain.MangasailChapter{
		Title:  <-titleChannel,
		Images: <-imagesChannel,
	}
}

func (r readerScrapper) getTitle(document goquery.Document) string {
	return document.Find("h1.page-header a.active").Text()
}

func (r readerScrapper) getImages(document goquery.Document) domain.MangasailImages {
	var images domain.MangasailImages
	document.Find("#node-337513 #images").Each(func(pos int, selection *goquery.Selection) {
		selection.Find("img").Each(func(iPos int, iSelection *goquery.Selection) {
			image := domain.MangasailImage{}
			image.ID, _ = iSelection.Attr("name")
			image.ImageURL, _ = iSelection.Attr("src")
			images = append(images, image)
		})
	})
	return images
}
