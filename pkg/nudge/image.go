package nudge

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"path"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/sajari/fuzzy"
	"golang.org/x/image/draw"
)

var (
	imageCache = make(map[string]image.Image)
)

type fuzzySort struct {
	arr  []string
	dist []int
}

const maxTokens = 5
const maxChoices = 7

type lengthSorter []string

func loadImage(kind, description string, maxSize fyne.Size) (image.Image, error) {
	description = strings.ToLower(description)
	path := path.Join("media", kind)
	imgPath, err := getImageFileName(path, description)
	if err != nil {
		return nil, err
	}

	if img, ok := imageCache[imgPath]; ok {
		return img, nil
	}

	f, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("Opening %s failed: %v", imgPath, err)
	}
	defer f.Close()
	src, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	whRatio := float32(src.Bounds().Dx()) / float32(src.Bounds().Dy())
	var img *image.RGBA
	if maxSize.Width > maxSize.Height {
		img = image.NewRGBA(image.Rect(0, 0, int(maxSize.Height*whRatio), int(maxSize.Height)))
	} else {
		img = image.NewRGBA(image.Rect(0, 0, int(maxSize.Width), int(maxSize.Width/whRatio)))
	}
	draw.NearestNeighbor.Scale(img, img.Bounds(), src, src.Bounds(), draw.Over, nil)

	imageCache[imgPath] = img
	return img, nil
}

func getImageFileName(pth, description string) (string, error) {
	dir, err := os.ReadDir(pth)
	if err != nil {
		return "", err
	}

	names := fuzzySort{
		arr: make([]string, 0, len(dir)),
	}
	for _, d := range dir {
		name := d.Name()
		if !strings.HasSuffix(name, ".png") {
			continue
		}
		if len(name) < 8 {
			continue
		}
		names.Add(path.Join(pth, name))
	}
	if names.Len() == 0 {
		return "", fmt.Errorf("no images")
	}

	return names.choose(description), nil
}

func (s *fuzzySort) choose(description string) string {
	s.precalc(description)

	sort.Sort(s)
	if s.Len() > maxChoices {
		s.arr = s.arr[:maxChoices]
	}

	choice := rand.Intn(s.Len())

	return s.arr[choice]
}

func (s *fuzzySort) Add(name string) {
	s.arr = append(s.arr, name)
}

func cleanText(text string) []string {
	for _, r := range ",.<>;'`!@#$%^&*():\\|/?[]{}\"-=~" {
		text = strings.ReplaceAll(text, string(r), " ")
	}
	textMap := make(map[string]struct{})
	for _, token := range strings.Split(strings.ToLower(text), " ") {
		if _, ok := textMap[token]; ok {
			continue
		}
		textMap[token] = struct{}{}
	}
	textArr := make([]string, 0, len(textMap))
	for token := range textMap {
		textArr = append(textArr, token)
	}
	sort.Sort(lengthSorter(textArr))

	if len(textArr) > maxTokens {
		return textArr[:maxTokens]
	}
	return textArr
}

func (s *fuzzySort) precalc(text string) {
	s.dist = make([]int, len(s.arr))
	textArr := cleanText(text)
	for i, a := range s.arr {
		a = path.Base(a)
		a = strings.ReplaceAll(a, "_", " ")
		a = strings.ReplaceAll(a, "-", " ")
		a = strings.ReplaceAll(a, ".png", "")
		s.dist[i] = calcDistance(strings.ToLower(a), textArr)
	}
}

func calcDistance(a string, b []string) int {
	min := 65535
	aArr := strings.Split(a, " ")
	for _, i := range aArr {
		for _, j := range b {
			dist := fuzzy.Levenshtein(&i, &j)
			if min > dist {
				min = dist
			}
		}
	}
	return min
}

func (s lengthSorter) Len() int { return len(s) }
func (s lengthSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s lengthSorter) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func (s *fuzzySort) Len() int { return len(s.arr) }
func (s *fuzzySort) Swap(i, j int) {
	s.arr[i], s.arr[j] = s.arr[j], s.arr[i]
	s.dist[i], s.dist[j] = s.dist[j], s.dist[i]
}
func (s *fuzzySort) Less(i, j int) bool {
	return s.dist[i] < s.dist[j]
}
