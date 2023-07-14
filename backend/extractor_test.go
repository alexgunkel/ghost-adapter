package backend_test

import (
	"os"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"

	"github.com/alexgunkel/ghost-adapter/backend"
)

func TestExtract(t *testing.T) {
	body, err := os.ReadFile("./example.json")
	if err != nil {
		panic(err)
	}

	res, err := backend.Extract(body)
	assert.NoError(t, err)
	assert.Len(t, res, 5)
	assert.Equal(t, "CodeCrafters – Fortbildung braucht einen neuen Ansatz", res[0].Title)
	assert.Equal(t, "https://blog.alexandergunkel.eu/content/images/2023/04/Penguins_in_classroom.png", res[0].FeatureImage)
	assert.Equal(t, "Pinguine während der Fortbildung", res[0].FeatureImageAlt)
	assert.Equal(t, "", res[0].FeatureImageCaption)
	assert.Equal(t, "https://blog.alexandergunkel.eu/codecrafters/", res[0].Url)
}

func TestExtractBuildsTeaser(t *testing.T) {
	body, err := os.ReadFile("./example.json")
	if err != nil {
		panic(err)
	}

	res, err := backend.Extract(body)
	assert.NoError(t, err)
	assert.Len(t, res, 5)
	assert.Equal(t,
		`Wer als Junior-Entwickler:in in die Berufswelt einsteigt, hat in aller Regel bereits fundierte Kenntnisse in mindestens einer Programmiersprache und bereits kleinere Projekte (oft zu Übungszwecken) absolviert. Um gute Senior- oder gar Lead-Developer zu werden, müssen Entwickler:innen nun aber weitere Fähigkeiten erwerben, die oft ganz anders geartet sind, als die ursprünglich erworbenen Programmierkenntnisse. Zunächst kommen Fähigkeiten hinzu, sauberen, wartbaren und aufgeräumten Code zu schreiben, Projekte zu analysieren und Schwachstellen zu verbessern. Von Senior-Entwickler:innen werden Kompetenzen im Bereich der Softwarearchitektur verlangt und Lead-Developer müssen insbesondere über kommunikative und Führungskompetenzen verfügen.`,
		res[0].Teaser,
	)
}

func FuzzConvert(f *testing.F) {
	f.Fuzz(func(t *testing.T, input string) {
		if utf8.ValidString(input) {
			assert.True(t, utf8.ValidString(backend.Convert(input)))
		}
	})
}
