package cmd

import (
	"bufio"
	"image"
	"log"
	"os"
	"os/signal"
	"time"

	sdeck "github.com/KarpelesLab/streamdeck"
	"github.com/KarpelesLab/streamdeck-examples/assets"
	"github.com/spf13/cobra"
)

var slideshowCmd = &cobra.Command{
	Use:   "slideshow",
	Short: "picture slideshow across all buttons",
	Long: `This example creates a slideshow on the Stream Deck, across all buttons.
Images of different formats (png, jpeg, gif) are loaded, resized to match
the panel and if necessary cropped to the center.`,
	Run: slideshow,
}

func init() {
	rootCmd.AddCommand(slideshowCmd)
}

func slideshow(cmd *cobra.Command, args []string) {

	sd, err := sdeck.NewStreamDeck()
	if err != nil {
		log.Panic(err)
	}
	defer sd.ClearAllBtns()

	_dices, err := assets.FS.Open("images/dices.png")
	if err != nil {
		log.Fatal(err)
	}
	defer _dices.Close()

	dices, _, err := image.Decode(bufio.NewReader(_dices))
	if err != nil {
		log.Panic(err)
	}

	_dna, err := assets.FS.Open("images/dna.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer _dna.Close()

	dna, _, err := image.Decode(bufio.NewReader(_dna))
	if err != nil {
		log.Panic(err)
	}

	_octocat, err := assets.FS.Open("images/octocat.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer _octocat.Close()

	octocat, _, err := image.Decode(bufio.NewReader(_octocat))
	if err != nil {
		log.Panic(err)
	}

	// start drawing octocat
	if err := sd.FillPanel(octocat); err != nil {
		log.Panic(err)
	}

	images := []image.Image{dices, dna, octocat}

	// launch a ticker for the slideshow
	ticker := time.NewTicker(time.Second * 3)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	position := 0
	for {
		select {
		case <-ticker.C:
			if err := sd.FillPanel(images[position]); err != nil {
				log.Panic(err)
			}
			if position == len(images)-1 {
				position = 0
				break
			}
			if position < len(images)-1 {
				position++
			}
		case <-c:
			return
		}
	}
}
