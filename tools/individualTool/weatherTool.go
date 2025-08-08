package individualtool

import (
	"log"
	"server/utils"

	owm "github.com/briandowns/openweathermap"
)

func WeatherTool(place string) (string, float64, string) {
	w, err := owm.NewCurrent("C", "en", utils.WEATHER_API)
	if err != nil {
		log.Println("Failed to initialize OpenWeather client:", err)
		return "", -256, ""
	}

	err = w.CurrentByName(place)
	if err != nil {
		log.Println("Failed to get weather:", err)
		return "", -256, ""
	}
	return w.Name, w.Main.Temp, w.Weather[0].Description
}
