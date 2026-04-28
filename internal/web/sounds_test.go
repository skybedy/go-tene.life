package web

import "testing"

func TestSoundTrackTitleUsesCzechDisplayNames(t *testing.T) {
	cases := map[string]string{
		"01_nakupovani_a_jidlo.mp3":              "Nakupování a jídlo",
		"05_lekar_a_zdravi.mp3":                  "Lékař a zdraví",
		"10_bydleni_a_domacnost_pokracovani.mp3": "Bydlení a domácnost, pokračování",
		"22_bezpecnost_a_nouzove_situace.mp3":    "Bezpečnost a nouzové situace",
		"spanelsko_ceska_slovicka_251_500.mp3":   "Španělsko-česká slovíčka 251-500",
	}

	for fileName, expected := range cases {
		if got := soundTrackTitle(fileName); got != expected {
			t.Fatalf("soundTrackTitle(%q) = %q, want %q", fileName, got, expected)
		}
	}
}
