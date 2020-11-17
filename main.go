package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	start := time.Now()
	// navigate to a page, wait for an element, click
	err := chromedp.Run(ctx,
		emulation.SetUserAgentOverride("WebScraper 1.0"),
		chromedp.Navigate(`https://github.com`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.ScrollIntoView(`footer`),
		chromedp.WaitVisible(`footer < div`),
		chromedp.Text(`h1`, &res, chromedp.NodeVisible, chromedp.ByQuery),
		// Prerender the DOM
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			res, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			fmt.Print(res)
			return err
		}),
		// Take a screenshot
		chromedp.Screenshot(`#some-element`, res, chromedp.NodeVisible, chromedp.ByID))

	// write the screenshot to the file after .Run()
	if err := ioutil.WriteFile("el-screenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("h1 contains: '%s'\n", res)
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}