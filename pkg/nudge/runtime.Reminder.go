package nudge

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/HealthyNudges/pkg/config"
)

func (r *nudgeRuntime) runReminder(cfg *config.Config, nudge *config.Nudge, how config.ExceptionHow) {
	adur, _ := cfg.ReminderAdvance.Get()
	advance := time.Duration(float64(time.Minute) * adur)
	reminder := MatchSchedule(nudge.Schedule, advance)
	event := MatchSchedule(nudge.Schedule, 0)

	if r.splash != nil {
		return
	}

	if reminder && !event {
		r.NotifyReminder(cfg, nudge)
	}

	if event {
		if nudge.Notification {
			r.NotifyEvent(cfg, nudge)
		}

		if !nudge.Window {
			return
		}

		title := canvas.NewText(nudge.Name, theme.ForegroundColor())
		title.Alignment = fyne.TextAlignCenter
		title.TextSize = theme.TextHeadingSize()

		description := canvas.NewText(nudge.Description, theme.ForegroundColor())
		description.Alignment = fyne.TextAlignCenter
		description.TextSize = theme.TextSubHeadingSize()

		legend := canvas.NewText("", theme.PlaceHolderColor())
		legend.Alignment = fyne.TextAlignCenter
		legend.TextSize = theme.TextSize()

		im, err := loadImage("remind", fmt.Sprint(nudge.Name, " ", nudge.Description), cfg.GetMaxImageSize())
		if err != nil {
			log.Println("Failed to load an image: ", err)
			return
		}
		img := canvas.NewImageFromImage(im)
		img.FillMode = canvas.ImageFillOriginal

		button := widget.NewButton("Got it!", func() {
			r.Lock()
			defer r.Unlock()
			if r.splash == nil {
				return
			}
			r.splash.Close()
		})
		content := container.NewVBox(title, description, container.NewCenter(img), legend, button)
		r.ShowEvent(cfg, nudge, content, nil)
	}
}
