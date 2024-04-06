package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {

	address := "Москва, Пресненская набережная, 2"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		//chromedp.ProxyServer("http://username:password@proxyserver.com:31280"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	var res []string

	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://sbermarket.ru/`),
		chromedp.WaitVisible(`img[alt="METRO"]`, chromedp.BySearch),
		chromedp.Click(`img[alt="METRO"]`, chromedp.NodeVisible),
		chromedp.WaitVisible(`//button[contains(text(), 'Выберите адрес доставки')]`, chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`//button[contains(text(), 'Выберите адрес доставки')]`, chromedp.NodeVisible),
		chromedp.WaitVisible(`.address-modal__content input`, chromedp.BySearch),
		chromedp.SendKeys(`.address-modal__content input`, address, chromedp.ByQuery),

		chromedp.WaitVisible(`div[class^="styles_option"]`, chromedp.BySearch),
		chromedp.Click(`div[class^="styles_option"]`, chromedp.NodeVisible),

		chromedp.WaitEnabled(`.address-modal__content button`, chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`.address-modal__content button`, chromedp.NodeVisible),

		chromedp.WaitNotPresent(`.address-modal__content button`, chromedp.BySearch),
		chromedp.Click(`//a[.//div[contains(text(), 'Овощи, фрукты, орехи')]]`, chromedp.NodeVisible),

		chromedp.WaitReady(`a[data-qa="category_department_taxons_list_taxon_item_0_link"]`, chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`a[data-qa="category_department_taxons_list_taxon_item_0_link"]`, chromedp.NodeVisible),
		chromedp.Sleep(5*time.Second),

		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('div[class^=ProductCard_root]')).map(el => {
				let imgElement = el.querySelector('[class*=ProductCard_image]');
				let imgUrl = imgElement.src;
				let largeImgUrl = imgUrl.replace('220-220', '1646-1646');
				let name = el.querySelector('[class*=ProductCard_titleContainer]').innerText;
				let pricePerKg = el.querySelector('[class*=ProductCard_volume]').innerText.trim();
				let discountedPrice = "";
				let proElement = el.querySelector('[class*=ProductCardPrice_original]');
				
				if ( proElement ) {
					let spanText = proElement.querySelector('span').innerText;
					discountedPrice = proElement.innerText.replace(spanText, '').trim() + "\n";
				}
				
				let prElement = el.querySelector('[class*=ProductCardPrice_price]');
				let spanText = prElement.querySelector('span').innerText;
				let originalPrice = prElement.innerText.replace(spanText, '').trim();
				
				return imgUrl + "\n"  + largeImgUrl + "\n"  + name + "\n"  + pricePerKg + "\n"  + discountedPrice + originalPrice + "\n\n";
			})
		`, &res),
	)

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile("products.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	for _, product := range res {
		_, err := fmt.Fprint(file, product)
		if err != nil {
			log.Println("Failed to write to file:", err)
			return
		}
	}
}
